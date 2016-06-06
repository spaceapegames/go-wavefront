package cmd

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	ui "github.com/gizak/termui"
	"github.com/spaceapegames/go-wavefront/api"
	"github.com/spf13/cobra"
)

var startTime, endTime int64
var granularity, period, maxPoints string
var liveGraph, includeObsolete bool
var maxGraphs int

// cmdQuery represents the query command
var cmdQuery = &cobra.Command{
	Use:   "query \"query-string\" [1m,3h,1d,1w] [ --end-time TIME ] ",
	Short: "Query Wavefront",
	Long: `Query Wavefront for time-series data.
If no time period given, returns data for the last hour.
Specify a period (e.g. 2h) to return data for period up to Now (i.e. last 2 hours).
Override Now with --end-time (i.e. 2 hours before end-time)
Alternatively provide a --start-time instead of a duration.
All times should be Unix epoch formatted.`,
	Run: queryCmd,
}

func init() {
	RootCmd.AddCommand(cmdQuery)

	cmdQuery.Flags().Int64VarP(&startTime, "start-time", "s", 0, "Start time for query in Unix epoch format.")
	cmdQuery.Flags().Int64VarP(&endTime, "end-time", "e", 0, "End time for query in Unix epoch format. Defaults to now")
	cmdQuery.Flags().StringVarP(&period, "period", "p", "1h", "Period of time before end time for which to query. Default 1h")
	cmdQuery.Flags().StringVarP(&granularity, "granularity", "g", "s", "Granularity of query. Defaults to seconds")
	cmdQuery.Flags().StringVarP(&maxPoints, "max-points", "x", "", "Maximum number of datapoints to return")
	cmdQuery.Flags().BoolVarP(&includeObsolete, "include-obsolete", "o", false, "Include obsolete metrics. Defaults to false.")
	cmdQuery.Flags().BoolVarP(&liveGraph, "live-graph", "l", false, "Display live graph of the output.")
	cmdQuery.Flags().IntVarP(&maxGraphs, "max-graphs", "m", 4, "Maximum number of graphs to show in live-graph mode. Max is 4.")
}

func queryCmd(cmd *cobra.Command, args []string) {
	createClient()

	if len(args) != 1 {
		log.Fatal("Please provide a query string (or --help for usage).")
	}

	queryString := args[0]

	params := wavefront.QueryParams{}

	params["g"] = granularity
	params["includeObsoleteMetrics"] = strconv.FormatBool(includeObsolete)

	if maxPoints != "" {
		params["p"] = maxPoints
	}

	query := client.Query.NewQuery(queryString)
	if endTime != 0 {
		query.SetEndTime(time.Unix(endTime, 0))
	} else {
		// the client defaults endTime to Now() but store
		// this here for periodSecs calculation below
		endTime = time.Now().Unix()
	}

	var periodSecs int64
	var err error
	if startTime != 0 {
		params["s"] = strconv.FormatInt(startTime, 10)
		periodSecs = endTime - startTime
	} else {
		if periodSecs, err = expandDuration(period); err != nil {
			log.Fatal(err)
		} else {
			query.SetStartTime(periodSecs)
		}
	}

	query.SetParams(params)

	res, err := query.Execute()
	if err != nil {
		log.Fatal(err)
	}

	if rawResponse == true {
		prettyPrint(&res.RawResponse)
		return
	}
	if liveGraph == true {
		doLiveGraph(res, query, periodSecs)
		return
	}

	// fall back to just printing the output
	fmt.Printf("%s", res)
}

// expandDuration converts a period like 1h, 2d etc.
// to seconds
func expandDuration(duration string) (int64, error) {
	r := regexp.MustCompile(`(\d+)([m,s,h,d,w])`)
	match := r.FindStringSubmatch(duration)
	if len(match) == 0 {
		return 0, errors.New("Please provide a valid duration (e.g. 10m,1h,2d,3w) or alternatively specify a start-time")
	}

	n, _ := strconv.Atoi(match[1]) //number
	d := match[2]                  //duration

	var seconds int

	switch d {
	case "s":
		seconds = 1
	case "m":
		seconds = 60
	case "h":
		seconds = 60 * 60
	case "d":
		seconds = 60 * 60 * 24
	case "w":
		seconds = 60 * 60 * 24 * 7
	}

	return int64(seconds * n), nil

}

//doLiveGraph builds a graph in the terminal window that
//updates every graphUpdate seconds
//It will build up to maxGraphs graphs with one time-series
//per graph
func doLiveGraph(res *wavefront.QueryResponse, query *wavefront.Querying, period int64) {

	err := ui.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	if maxGraphs > len(res.TimeSeries) {
		maxGraphs = len(res.TimeSeries)
	}

	var wDivisor, hDivisor int
	switch maxGraphs {
	case 1:
		wDivisor = 1
		hDivisor = 1
	case 2:
		wDivisor = 2
		hDivisor = 1
	case 3, 4:
		wDivisor = 2
		hDivisor = 2
	}

	height := ui.TermHeight() / hDivisor
	width := ui.TermWidth() / wDivisor
	xVals, yVals := calculateCoords(maxGraphs, ui.TermWidth()/wDivisor, ui.TermHeight()/hDivisor)
	graphs := buildGraphs(res, height, width, xVals, yVals)

	ui.Render(graphs...)
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		// press q to quit
		ui.StopLoop()
	})

	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		// handle Ctrl + c combination
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {
		query.SetEndTime(time.Now())
		query.SetStartTime(period)
		res, err := query.Execute()
		if err != nil {
			log.Fatal(err)
		}
		graphs := buildGraphs(res, height, width, xVals, yVals)
		ui.Render(graphs...)
	})

	ui.Loop()

}

func buildGraphs(res *wavefront.QueryResponse, height, width int, xVals, yVals []int) []ui.Bufferer {

	graphs := make([]ui.Bufferer, maxGraphs)

	for i, timeSeries := range res.TimeSeries[:maxGraphs] {
		datapoints := make([]float64, len(timeSeries.DataPoints))
		datalabels := make([]string, len(timeSeries.DataPoints))
		for i, d := range timeSeries.DataPoints {
			datapoints[i] = d[1]
			datalabels[i] = time.Unix(int64(d[0]), 0).Format(time.Kitchen)
		}

		var label string
		if timeSeries.Host != "" {
			label = timeSeries.Host
		} else {
			label = res.Query
		}

		lc := ui.NewLineChart()
		lc.BorderLabel = label
		lc.Data = datapoints
		lc.DataLabels = datalabels
		lc.Width = width
		lc.Height = height
		lc.X = xVals[i]
		lc.Y = yVals[i]
		lc.AxesColor = ui.ColorWhite
		lc.LineColor = ui.ColorYellow | ui.AttrBold

		graphs[i] = lc

	}
	return graphs
}

func calculateCoords(numGraphs, eachWidth, eachHeight int) ([]int, []int) {
	switch numGraphs {
	case 1:
		return []int{0}, []int{0}
	case 2:
		return []int{0, eachWidth}, []int{0, 0}
	case 3, 4:
		return []int{0, eachWidth, eachWidth, 0},
			[]int{0, eachHeight, 0, eachHeight}
	}
	return []int{}, []int{}
}

package wavefront

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

type DataPoint []float64

type TimeSeries struct {
	DataPoints []DataPoint
	Label      string
	Host       string
}

// QueryReponse stores the Wavefront query response
type QueryResponse struct {
	RawResponse   []byte
	RawTimeSeries []interface{} `json:"timeseries"`
	TimeSeries    []TimeSeries
	Query         string         `json:"query"`
	Stats         map[string]int `json:"stats"`
	Name          string         `json:"name"`
	Granularity   int            `json:"granularity"`
	Hosts         []string       `json:"hostsUsed"`
}

// Querying is used to query the charts API
type Querying struct {
	client *Client
	params QueryParams
}

const (
	baseQueryPath = "/chart/api"
	// some constants provided for time convenience
	LAST_HOUR     = 60 * 60
	LAST_3_HOURS  = LAST_HOUR * 3
	LAST_6_HOURS  = LAST_HOUR * 6
	LAST_24_HOURS = LAST_HOUR * 24
)

// Execute is used to execute a query against the Wavefront Chart API
func (q Querying) Execute() (*QueryResponse, error) {
	var queryResp *QueryResponse

	// quick sanity check
	start, _ := strconv.Atoi(q.params["s"])
	end, _ := strconv.Atoi(q.params["e"])
	if start > end {
		return nil, errors.New("Query start time is after end time")
	}

	req, err := q.client.NewRequest("GET", baseQueryPath, &q.params)
	if err != nil {
		return nil, err
	}
	resp, err := q.client.Do(req, &queryResp)
	if err != nil {
		return nil, err
	}

	queryResp.RawResponse, err = ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	// Parse the timeseries into a format we can work with
	// This involves an abominable type assertion
	for _, t := range queryResp.RawTimeSeries {
		tt := t.(map[string]interface{})
		ts := TimeSeries{Label: tt["label"].(string),
			DataPoints: []DataPoint{},
		}
		// Not all queries return Host
		if h, ok := tt["host"]; ok {
			ts.Host = h.(string)
		}

		for _, d := range (tt["data"]).([]interface{}) {
			timestamp, value := d.([]interface{})[0].(float64), d.([]interface{})[1].(float64)
			ts.DataPoints = append(ts.DataPoints, []float64{timestamp, value})
		}

		// Not interested in empty datasets
		if len(ts.DataPoints) != 0 {
			queryResp.TimeSeries = append(queryResp.TimeSeries, ts)
		}
	}

	return queryResp, nil
}

// NewQuery builds a query based on a query string and an optional
// set of key-value params
// The default end time will be Now(), and the start time one hour ago
// Strict mode is set to True by default
func (q *Querying) NewQuery(query string, opts ...map[string]string) *Querying {
	q.params = QueryParams{"q": query}
	if len(opts) > 0 {
		q.SetParams(opts[0])
	}

	// Default the start time to one hour ago
	if _, found := q.params["s"]; !found {
		hourPast := strconv.FormatInt(time.Now().Unix()-3600, 10)
		q.SetParams(map[string]string{"s": hourPast})
	}

	// Default end time to now
	if _, found := q.params["e"]; !found {
		now := strconv.FormatInt(time.Now().Unix(), 10)
		q.SetParams(map[string]string{"e": now})
	}

	// Default strict mode to True
	if _, found := q.params["strict"]; !found {
		q.SetParams(map[string]string{"strict": "true"})
	}
	return q
}

// SetParams adds or sets query parameters on a Querying object
func (q *Querying) SetParams(p QueryParams) {
	for opt, val := range p {
		q.params[opt] = val
	}
}

// GetParams returns a map of parameters
func (q Querying) GetParams() QueryParams {
	return q.params
}

// SetStartTime sets the time from which to query for points
// seconds is the number of seconds before the end-time that the query should encompass
// Some constants are provided as helpers: LAST_HOUR, LAST_3_HOURS etc
func (q *Querying) SetStartTime(seconds int64) {
	end, _ := strconv.Atoi(q.params["e"])
	start := strconv.FormatInt(int64(end)-seconds, 10)
	q.SetParams(QueryParams{"s": start})
}

// SetEndTime sets the time at which the query should end
// endTime is a time.Time type
func (q *Querying) SetEndTime(endTime time.Time) {
	q.SetParams(QueryParams{"e": strconv.FormatInt(endTime.Unix(), 10)})
}

// String outputs the time-series of a QueryResponse object
// in a human-readable format
func (q QueryResponse) String() string {
	var out string
	for _, t := range q.TimeSeries {
		if t.Host != "" {
			out += fmt.Sprintf("%s : %s\n", t.Label, t.Host)
		} else {
			out += fmt.Sprintf("%s\n", t.Label)
		}
		for _, d := range t.DataPoints {
			out += fmt.Sprintf("%d %f\n", int64(d[0]), d[1])
		}
	}
	return out
}

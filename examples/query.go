package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/WavefrontHQ/go-wavefront-management-api"
)

func main() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	client, err := wavefront.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// enable debug - all requests get dumped to stdout before being executed
	client.Debug(true)

	// NewQueryParams generates a query using the given ts expression.
	// By default the query period will be one hour since the current time.
	query := client.NewQuery(wavefront.NewQueryParams(
		`ts("cpu.load.1m.avg", dc=dc1)`,
	))

	// Set the query period to be one day instead of one hour
	query.SetStartTime(int64(time.Duration.Hours(24) * 60 * 60))

	// Execute carries out the query
	result, err := query.Execute()
	if err != nil {
		log.Fatal(err)
	}

	// The raw JSON response is available as RawResponse.
	// This can be useful for debugging
	b, _ := ioutil.ReadAll(result.RawResponse)
	fmt.Println(string(b))

	// The timeseries response can now be used to explore the results
	fmt.Println(result.TimeSeries[0].Label)
	fmt.Println(result.TimeSeries[0].DataPoints[0])

}

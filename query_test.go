package wavefront

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

type MockWavefrontClient struct {
	Client
	Response []byte
}

func (m MockWavefrontClient) Do(req *http.Request) (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewReader(m.Response)), nil
}

func TestQuery(t *testing.T) {
	query := "dudndun"
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	q := &Query{
		Params: NewQueryParams(query),
		client: &MockWavefrontClient{
			Response: []byte(`{"valid":"json"}`),
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
		},
	}

	// check correct default timewindow applied
	end, _ := strconv.Atoi(q.Params.EndTime)
	start, _ := strconv.Atoi(q.Params.StartTime)
	if end-start != 3600 {
		t.Errorf("query window, expected 3600, got %d", end-start)
	}

	q.SetEndTime(time.Now())
	q.SetStartTime(LastDay)
	end, _ = strconv.Atoi(q.Params.EndTime)
	start, _ = strconv.Atoi(q.Params.StartTime)
	if end-start != LastDay {
		t.Errorf("query window, expected %d, got %d", LastDay, end-start)
	}

	resp, err := q.Execute()
	if err != nil {
		t.Fatal("error executing query:", err)
	}

	raw, err := ioutil.ReadAll(resp.RawResponse)
	if err != nil {
		t.Error(err)
	}

	if err := json.Unmarshal(raw, new(map[string]interface{})); err != nil {
		fmt.Println(string(raw))
		t.Error("raw response is invalid JSON", err)
	}
}

func getQueryOutputFromFixture(fixture string) (*QueryResponse, error) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	response, err := ioutil.ReadFile(fixture)
	if err != nil {
		return nil, err
	}
	q := &Query{
		Params: NewQueryParams("ts(some.query)"),
		client: &MockWavefrontClient{
			Response: response,
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
		},
	}
	return q.Execute()
}

func TestQuery_SingleSeries(t *testing.T) {
	resp, err := getQueryOutputFromFixture("./fixtures/single-series.json")
	if err != nil {
		t.Fatal("error executing query:", err)
	}

	ts := resp.TimeSeries[0]
	if ts.Host != "server1.example.net" {
		t.Errorf("timeseries host, exepected server1.example.net, got %s", ts.Host)
	}

	if ts.Label != "servers.load.load.longterm" {
		t.Errorf("timeseries label, exepected servers.load.load.longterm, got %s", ts.Label)
	}

	if len(ts.DataPoints) != 60 {
		t.Errorf("datapoints, expected 60, got %d", len(ts.DataPoints))
	}

}

func TestQuery_MultiSeries(t *testing.T) {
	resp, err := getQueryOutputFromFixture("./fixtures/multi-series.json")
	if err != nil {
		t.Fatal("error executing query:", err)
	}

	if len(resp.TimeSeries) != 2 {
		t.Fatalf("expected 2 timeseries, got %d", len(resp.TimeSeries))
	}

	ts := resp.TimeSeries[0]
	if ts.Host != "server1.example.net" {
		t.Errorf("timeseries host, exepected server1.example.net, got %s", ts.Host)
	}

	if ts.Label != "servers.load.load.longterm" {
		t.Errorf("timeseries label, exepected servers.load.load.longterm, got %s", ts.Label)
	}

	if len(ts.DataPoints) != 59 {
		t.Errorf("datapoints, expected 59, got %d", len(ts.DataPoints))
	}

	ts = resp.TimeSeries[1]
	if ts.Host != "server2.example.net" {
		t.Errorf("timeseries host, exepected server2.example.net, got %s", ts.Host)
	}

	if ts.Label != "servers.load.load.longterm" {
		t.Errorf("timeseries label, exepected servers.load.load.longterm, got %s", ts.Label)
	}

	if len(ts.DataPoints) != 59 {
		t.Errorf("datapoints, expected 59, got %d", len(ts.DataPoints))
	}

}

func TestQuery_Error(t *testing.T) {
	resp, err := getQueryOutputFromFixture("./fixtures/failed-query.json")
	if err != nil {
		t.Fatal("error executing query:", err)
	}

	// must have the ErrType field populated if there is one
	if resp.ErrType == "" {
		t.Errorf("Expected to have ErrType field populated.")
	}

	if resp.ErrType != "QueryExecutionFailed" {
		t.Errorf("Unexpected Error Type found.")
	}

	// must have the ErrType field populated if there is one
	if resp.ErrMessage == "" {
		t.Errorf("Expected to have ErrMessage field populated.")
	}

	if resp.ErrMessage != "Terminating Query to save resources" {
		t.Errorf("Unexpected Error Message found.")
	}
}

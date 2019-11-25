package wavefront

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"time"
)

// Query represents a query to be made against the Charts API
type Query struct {
	// client is the Wavefront client used to execute queries
	client Wavefronter

	// Params is the set of parameters that will be used when executing the Query
	Params *QueryParams

	// Response will be the response of the last executed Query
	Response *QueryResponse
}

// QueryParams represents parameters that will be passed when making a Query
type QueryParams struct {
	// Name is an optional name to identify the query
	Name string `query:"n"`

	// QueryString is the actual timeseries query to be executed
	QueryString string `query:"q"`

	// StartTime is the start time for the query in epoch milliseconds
	StartTime string `query:"s"`

	// EndTime is the end time for the query in epoch milliseconds
	EndTime string `query:"e"`

	// Granularity is the granularity of the points returned, and can be one of
	// d,h,m or s
	Granularity string `query:"g"`

	// MaxPoints is the maximum number of points to return
	MaxPoints string `query:"p"`

	// SeriesOutsideTimeWindow is a boolean to indicate whether series with only
	// points that are outside  of the query window will be returned
	SeriesOutsideTimeWindow bool `query:"i"`

	// AutoEvents is a boolean to indicate whether to return Events for sources
	// included in the query
	AutoEvents bool `query:"autoEvents"`

	// SummarizationStrategy is the strategy to be used when grouping points together.
	// Valid values are MEAN, MEDIAN, MIN, MAX, SUN, COUNT, LAST, FIRST
	SummarizationStrategy string `query:"summarization"`

	// ListMode is a boolean to indicate whether to retrieve events more optimally
	// displayed for a list.
	ListMode bool `query:"listmode"`

	// StrictMode is a boolean which, if true, will not return points outside of the
	// query window. Defaults to false if ommitted.
	StrictMode bool `query:"strict"`

	// IncludeObsoleteMetrics is a boolean to indicate whether to return points from
	// sources which have stopped reporting. Defaults to false if ommitted.
	IncludeObsoleteMetrics bool `query:"includeObsoleteMetrics"`
}

// QueryResponse is used to represent a Wavefront query response
type QueryResponse struct {
	RawResponse *bytes.Reader
	TimeSeries  []TimeSeries   `json:"timeseries"`
	Query       string         `json:"query"`
	Stats       map[string]int `json:"stats"`
	Name        string         `json:"name"`
	Granularity int            `json:"granularity"`
	Hosts       []string       `json:"hostsUsed"`
	Warnings    string         `json:"warnings"`

	// ErrType : ref https://code.vmware.com/apis/714/wavefront-rest#/Query/queryApi
	ErrType string `json:"errorType"`

	// ErrMessage : ref https://code.vmware.com/apis/714/wavefront-rest#/Query/queryApi
	ErrMessage string `json:"errorMessage"`
}

// DataPoint represents a single timestamp/value data point as returned
// by Wavefront
type DataPoint []float64

// TimeSeries represents a single TimeSeries as returned by Wavefront
type TimeSeries struct {
	DataPoints []DataPoint       `json:"data"`
	Label      string            `json:"label"`
	Host       string            `json:"host"`
	Tags       map[string]string `json:"tags"`
}

const (
	baseQueryPath = "/api/v2/chart/api"
	// some constants provided for time convenience
	LastHour    = 60 * 60
	Last3Hours  = LastHour * 3
	Last6Hours  = LastHour * 6
	Last24Hours = LastHour * 24
	LastDay     = Last24Hours
	LastWeek    = LastDay * 7
)

// NewQueryParams takes a query string and returns a set of QueryParams with
// a query window of one hour since now and a set of sensible default vakues
func NewQueryParams(query string) *QueryParams {
	endTime := time.Now().Unix()
	startTime := endTime - LastHour
	return &QueryParams{
		QueryString: query,
		EndTime:     strconv.FormatInt(endTime, 10),
		StartTime:   strconv.FormatInt(startTime, 10),
		Granularity: "s",
		StrictMode:  true,
	}
}

func NewQueryParamsNoStrict(query string) *QueryParams {
	endTime := time.Now().Unix()
	startTime := endTime - LastHour
	return &QueryParams{
		QueryString: query,
		EndTime:     strconv.FormatInt(endTime, 10),
		StartTime:   strconv.FormatInt(startTime, 10),
		Granularity: "s",
		StrictMode:  false,
	}
}

// NewQuery returns a Query based on QueryParams
func (c *Client) NewQuery(params *QueryParams) *Query {
	return &Query{
		client: c,
		Params: params,
	}
}

// Execute is used to execute a query against the Wavefront Chart API
func (q *Query) Execute() (*QueryResponse, error) {
	queryResp := &QueryResponse{}

	params := map[string]string{}

	qpType := reflect.TypeOf(q.Params).Elem()
	qp := reflect.ValueOf(q.Params).Elem()

	for i := 0; i < qpType.NumField(); i++ {
		if qp.Field(i).String() != "" {

			if qp.Field(i).Type().String() == "bool" {
				params[qpType.Field(i).Tag.Get("query")] = strconv.FormatBool(qp.Field(i).Bool())
			} else {
				params[qpType.Field(i).Tag.Get("query")] = qp.Field(i).String()
			}
		}
	}

	req, err := q.client.NewRequest("GET", baseQueryPath, &params, nil)
	if err != nil {
		return nil, err
	}
	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	// bytes.Reader implements Seek, which we need to use to 'rewind' the Body below
	queryResp.RawResponse = bytes.NewReader(body)
	err = json.Unmarshal(body, queryResp)
	if err != nil {
		return nil, err
	}

	// 'rewind' the raw response
	queryResp.RawResponse.Seek(0, 0)

	return queryResp, nil
}

// SetStartTime sets the time from which to query for points.
// 'seconds' is the number of seconds before the end-time that the query will
// be inclusive of. EndTime must be set before calling this function.
// Some constants are provided for convenience: LastHour, Last3Hours, LastDay etc.
func (q *Query) SetStartTime(seconds int64) error {
	if q.Params.EndTime == "" {
		return fmt.Errorf("ensure end-time is configured")
	}
	end, err := strconv.Atoi(q.Params.EndTime)
	if err != nil {
		return err
	}
	q.Params.StartTime = strconv.FormatInt(int64(end)-seconds, 10)
	return nil
}

// SetEndTime sets the time at which the query should end
func (q *Query) SetEndTime(endTime time.Time) {
	q.Params.EndTime = strconv.FormatInt(endTime.Unix(), 10)
}

func (qr *QueryResponse) UnmarshalJSON(data []byte) error {
	// Aliasing the type avoids recursive calls to UnmarshalJSON
	type Alias QueryResponse
	tmp := struct {
		*Alias
	}{
		Alias: (*Alias)(qr),
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	return nil
}

// String outputs the time-series of a QueryResponse object
// in a human-readable format
func (qr QueryResponse) String() string {
	var out string
	if qr.Warnings != "" {
		out += fmt.Sprintf("Warnings : %s\n", qr.Warnings)
	}
	for _, t := range qr.TimeSeries {
		if t.Host != "" {
			out += fmt.Sprintf("%s : %s\n", t.Label, t.Host)
		} else {
			out += fmt.Sprintf("%s\n", t.Label)
		}
		if t.Tags != nil {
			for k, v := range t.Tags {
				out += fmt.Sprintf("%s : %s\n", k, v)
			}
		}
		for _, d := range t.DataPoints {
			out += fmt.Sprintf("%d %f\n", int64(d[0]), d[1])
		}
	}
	return out
}

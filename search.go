package wavefront

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// Search represents a search to be made against the Search API
type Search struct {
	// client is the Wavefront client used to effect the Search
	client Wavefronter

	// Type is the type of entity to be searched for (i.e. alert, event, dashboard,
	// extlink, cloudintegration etc.)
	Type string

	// Params are the Search parameters to be applied
	Params *SearchParams

	// Deleted is whether to search against the /{entity}/deleted endpoint (for
	// deleted items) instead of the normal one. Defaults to false.
	Deleted bool
}

// SearchParams represents paramaters used to effect a Search.
// If multiple search terms are given they will act like a logical AND.
// If Conditions is nil, all items of the given Type will be returned
type SearchParams struct {
	// Conditions are the search conditions to be matched.
	// If multiple are given they will act like a logical AND
	Conditions []*SearchCondition `json:"query"`

	// Limit is the max number of results to be returned. Defaults to 100.
	Limit int `json:"limit"`

	// Offset is the offset from the first result to be returned.
	// For instance, an Offset of 100 will yield results 101 - 200
	// (assuming a Limit of 100). Defaults to 0.
	Offset int `json:"offset"`

	// TimeRange is the range between which results will be searched.
	// This is only valid for certain search types (e.g. Events)
	TimeRange *TimeRange `json:"timeRange,omitempty"`
}

// TimeRange represents a range of times to search between. It is only valid
// for certain search types (e.g. Events)
type TimeRange struct {
	// StartTime is the time, in epoch milliseconds from which search results
	// will be returned.
	StartTime int64 `json:"earliestStartTimeEpochMillis"`

	// EndTime is the time, in epoch milliseconds up to which search results
	// will be returned.
	EndTime int64 `json:"latestStartTimeEpochMillis"`
}

// SearchCondition represents a single search condition.
// Multiple conditions can be applied to one search, they will act as a logical AND.
type SearchCondition struct {
	// Key is the type of parameter to be matched (e.g. tags, status, id)
	Key string `json:"key"`

	// Value is the value of Key to be searched for (e.g. the tag name, or snoozed)
	Value string `json:"value"`

	// MatchingMethod must be one of CONTAINS, STARTSWITH, EXACT, TAGPATH
	MatchingMethod string `json:"matchingMethod"`
}

// SearchResponse represents the result of a successful search operation
type SearchResponse struct {
	// RawResponse is the raw JSON response returned by Wavefront from a Search
	// operation
	RawResponse *bytes.Reader

	// Response is the response body of a Search operation
	Response struct {
		// Items will be the Wavefront entities returned by a successful search
		// operation (i.e. the Alerts, or Dashboards etc.)
		Items json.RawMessage

		// MoreResults indicates whether there are further items to be returned in a
		// paginated response.
		MoreItems bool `json:"moreItems"`
	} `json:"response"`

	// NextOffset is the offset that should be used to retrieve the next page of
	// results in a paginated response. If there are no more results, it will be zero.
	NextOffset int
}

const baseSearchPath = "/api/v2/search"

// NewSearch returns a Search based on SearchParams.
// searchType is the type of entity to be searched for (i.e. alert, event, dashboard,
// extlink, cloudintegration etc.)
func (c *Client) NewSearch(searchType string, params *SearchParams) *Search {
	paramsCopy := *params
	return &Search{
		client:  c,
		Type:    searchType,
		Params:  &paramsCopy,
		Deleted: false,
	}
}

// NewTimeRange returns a *TimeRange encompassing the period seconds before the given
// endTime. If endTime is 0, the current time will be used.
func NewTimeRange(endTime, period int64) (*TimeRange, error) {
	if endTime == 0 {
		endTime = time.Now().Unix()
	}
	if period < 0 {
		return nil, fmt.Errorf("time period must be a positive number")
	}
	startTime := endTime - period
	return &TimeRange{
		StartTime: startTime * 1000,
		EndTime:   endTime * 1000,
	}, nil
}

// Execute is used to carry out a search
func (s *Search) Execute() (*SearchResponse, error) {
	paramsCopy := *s.Params
	// set defaults
	if paramsCopy.Limit == 0 {
		paramsCopy.Limit = 100
	}

	path := baseSearchPath + "/" + s.Type
	if s.Deleted {
		path += "/deleted"
	}

	payload, err := json.Marshal(&paramsCopy)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("POST", path, nil, payload)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	searchResp := &SearchResponse{}
	// bytes.Reader implements Seek, which we need to use to 'rewind' the Body below
	searchResp.RawResponse = bytes.NewReader(body)
	err = json.Unmarshal(body, searchResp)
	if err != nil {
		return nil, err
	}

	if searchResp.Response.MoreItems {
		searchResp.NextOffset = paramsCopy.Offset + paramsCopy.Limit
	} else {
		searchResp.NextOffset = 0
	}

	// 'rewind' the raw response
	_, err = searchResp.RawResponse.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return searchResp, nil
}

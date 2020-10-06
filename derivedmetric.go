package wavefront

import (
	"fmt"
	"strconv"
)

type DerivedMetric struct {
	ID                       *string  `json:"id,omitempty"`
	Name                     string   `json:"name,omitempty"`
	Query                    string   `json:"query,omitempty"`
	Minutes                  int      `json:"minutes,omitempty"`
	Tags                     WFTags   `json:"tags,omitempty"`
	Status                   []string `json:"status,omitempty"`
	InTrash                  bool     `json:"inTrash,omitempty"`
	QueryFailing             bool     `json:"queryFailing,omitempty"`
	LastFailedTime           int      `json:"lastFailedTime,omitempty"`
	LastErrorMessage         string   `json:"lastErrorMessage,omitempty"`
	AdditionalInformation    string   `json:"additionalInformation,omitempty"`
	HostsUsed                []string `json:"hostsUsed,omitempty"`
	UpdateUserId             string   `json:"updateUserId,omitempty"`
	CreateUserId             string   `json:"createUserId,omitempty"`
	LastProcessedMillis      int      `json:"lastProcessedMillis,omitempty"`
	ProcessRateMinutes       int      `json:"processRateMinutes,omitempty"`
	PointsScannedAtLastQuery int      `json:"pointsScannedAtLastQuery,omitempty"`
	IncludeObsoleteMetrics   bool     `json:"includeObsoleteMetrics,omitempty"`
	LastQueryTime            int      `json:"lastQueryTime,omitempty"`
	MetricsUsed              []string `json:"metricsUsed,omitempty"`
	QueryQBEnabled           bool     `json:"queryQBEnabled,omitempty"`
	UpdatedEpochMillis       int      `json:"updatedEpochMillis,omitempty"`
	CreatedEpochMillis       int      `json:"createdEpochMillis,omitempty"`
	Deleted                  bool     `json:"deleted,omitempty"`
}

type DerivedMetrics struct {
	client Wavefronter
}

type WFTags struct {
	CustomerTags []string `json:"customerTags"`
}

const baseDerivedMetricsPath = "/api/v2/derivedmetric"

func (c *Client) DerivedMetrics() *DerivedMetrics {
	return &DerivedMetrics{client: c}
}

// Get is used to retrieve an existing DerivedMetric by ID.
// The ID field must be specified
func (dm DerivedMetrics) Get(metric *DerivedMetric) error {
	if metric.ID == nil || *metric.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", baseDerivedMetricsPath, *metric.ID),
		dm.client,
		doResponse(metric))
}

// Find returns all DerivedMetrics filtered by the given search conditions.
// If filter is nil, all DerivedMetrics are returned.
func (dm DerivedMetrics) Find(filter []*SearchCondition) (
	results []*DerivedMetric, err error) {
	err = doSearch(filter, "derivedmetric", dm.client, &results)
	return
}

// Create a DerivedMetric, name, query, and minutes are required
func (dm DerivedMetrics) Create(metric *DerivedMetric) error {
	if metric.Name == "" || metric.Query == "" || metric.Minutes == 0 {
		return fmt.Errorf("name, query, and minutes must be specified to create a derived metric")
	}

	return doRest(
		"POST",
		baseDerivedMetricsPath,
		dm.client,
		doPayload(metric),
		doResponse(metric))
}

// Update a DerivedMetric all fields are optional except for ID
func (dm DerivedMetrics) Update(metric *DerivedMetric) error {
	if metric.ID == nil || *metric.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", baseDerivedMetricsPath, *metric.ID),
		dm.client,
		doPayload(metric),
		doResponse(metric))
}

// Delete a DerivedMetric all fields are optional except for ID
func (dm DerivedMetrics) Delete(metric *DerivedMetric, skipTrash bool) error {
	if metric.ID == nil || *metric.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	params := map[string]string{
		"skipTrash": strconv.FormatBool(skipTrash),
	}

	err := doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", baseDerivedMetricsPath, *metric.ID),
		dm.client,
		doParams(params))
	if err != nil {
		return err
	}
	empty := ""
	metric.ID = &empty
	return nil
}

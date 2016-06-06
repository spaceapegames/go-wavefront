package wavefront

import (
	"io/ioutil"
)

// Alert type is used to store individual alerts
type Alert struct {
	Name         string                 `json:"name"`
	MetricsUsed  []string               `json:"metricsUsed"`
	UserTags     map[string]int         `json:"userTagsWithCounts"`     // user tags along with a count thereof
	CustomerTags map[string]int         `json:"customerTagsWithCounts"` // customer tags along with a count thereof
	Severity     string                 `json:"severity"`
	Hosts        []string               `json:"hostsUsed"`
	Condition    string                 `json:"condition"`
	Event        map[string]interface{} `json:"event"`
}

// Alerting is used to query the Alerts API
type Alerting struct {
	client      *Client
	RawResponse []byte //Raw JSON response of the last request
}

// baseAlertPath is the base API path for retrieving alerts
const baseAlertPath = "/api/alert/"

// getAlerts is used to retrieve alerts given a specific type.
// There should be little need to call it directly, one of the convenience wrappers should be used instead.
func (a *Alerting) getAlerts(alertType string, params *QueryParams) ([]*Alert, error) {
	var alerts []*Alert
	req, err := a.client.NewRequest("GET", baseAlertPath+alertType, params)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req, &alerts)
	if err != nil {
		return nil, err
	}
	a.RawResponse, err = ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

// All returns all alerts
func (a *Alerting) All(params *QueryParams) ([]*Alert, error) {
	return a.getAlerts("all", params)
}

// Active returns active(firing) alerts
func (a *Alerting) Active(params *QueryParams) ([]*Alert, error) {
	return a.getAlerts("active", params)
}

// Snoozed returns snoozed alerts
func (a *Alerting) Snoozed(params *QueryParams) ([]*Alert, error) {
	return a.getAlerts("snoozed", params)
}

// Invalid returns invalid alerts
func (a *Alerting) Invalid(params *QueryParams) ([]*Alert, error) {
	return a.getAlerts("invalid", params)
}

// AffectedByMaintenance returns alerts which are currently affected by maintenance windows
func (a *Alerting) AffectedByMaintenance(params *QueryParams) ([]*Alert, error) {
	return a.getAlerts("affected_by_maintenance", params)
}

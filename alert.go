package wavefront

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	AlertTypeThreshold = "THRESHOLD"
	AlertTypeClassic   = "CLASSIC"
)

// Alert represents a single Wavefront Alert
type Alert struct {
	// Name is the name given to an Alert
	Name string `json:"name"`

	// ID is the Wavefront-assigned ID of an existing Alert
	ID *string `json:"id,omitempty"`

	// AlertType should be either CLASSIC or THRESHOLD
	AlertType string `json:"alertType,omitempty"`

	// AdditionalInfo is any extra information about the Alert
	AdditionalInfo string `json:"additionalInformation"`

	// Target is a comma-separated list of targets for the Alert
	Target string `json:"target,omitempty"`

	// For THRESHOLD alerts. Targets is a map[string]string. This maps severity to lists of targets.
	// Valid keys are: severe, smoke, warn or info
	Targets map[string]string `json:"targets"`

	// Condition is the condition under which the Alert will fire
	Condition string `json:"condition"`

	// For THRESHOLD alerts. Conditions is a map[string]string. This maps severity to respective conditions.
	// Valid keys are: severe, smoke, warn or info
	Conditions map[string]string `json:"conditions"`

	// DisplayExpression is the ts query to generate a graph of this Alert, in the UI
	DisplayExpression string `json:"displayExpression,omitempty"`

	// Minutes is the number of minutes the Condition must be met, before the
	// Alert will fire
	Minutes int `json:"minutes"`

	// ResolveAfterMinutes is the number of minutes the Condition must be un-met
	// before the Alert is considered resolved
	ResolveAfterMinutes int `json:"resolveAfterMinutes,omitempty"`

	// Minutes to wait before re-sending notification of firing alert.
	NotificationResendFrequencyMinutes int `json:"notificationResendFrequencyMinutes"`

	// Severity is the severity of the Alert, and can be one of SEVERE,
	// SMOKE, WARN or INFO
	Severity string `json:"severity,omitempty"`

	// For THRESHOLD alerts. SeverityList is a list of strings. Different severities applicable to this alert.
	// Valid elements are: SEVERE, SMOKE, WARN or INFO
	SeverityList []string `json:"severityList"`

	// Status is the current status of the Alert
	Status []string `json:"status"`

	// Tags are the tags applied to the Alert
	Tags []string

	FailingHostLabelPairs       []SourceLabelPair `json:"failingHostLabelPairs,omitempty"`
	InMaintenanceHostLabelPairs []SourceLabelPair `json:"inMaintenanceHostLabelPairs,omitempty"`
}

type SourceLabelPair struct {
	Host   string `json:"host"`
	Firing int    `json:"firing"`
}

// Alerts is used to perform alert-related operations against the Wavefront API
type Alerts struct {
	// client is the Wavefront client used to perform alert-related operations
	client Wavefronter
}

const baseAlertPath = "/api/v2/alert"

// UnmarshalJSON is a custom JSON unmarshaller for an Alert, used in order to
// populate the Tags field in a more intuitive fashion
func (a *Alert) UnmarshalJSON(b []byte) error {
	type alert Alert
	temp := struct {
		Tags map[string][]string `json:"tags"`
		*alert
	}{
		alert: (*alert)(a),
	}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	a.Tags = temp.Tags["customerTags"]
	return nil
}

func (a *Alert) MarshalJSON() ([]byte, error) {
	type alert Alert
	return json.Marshal(&struct {
		Tags map[string][]string `json:"tags"`
		*alert
	}{
		Tags: map[string][]string{
			"customerTags": a.Tags,
		},
		alert: (*alert)(a),
	})
}

// Alerts is used to return a client for alert-related operations
func (c *Client) Alerts() *Alerts {
	return &Alerts{client: c}
}

// Get is used to retrieve an existing Alert by ID.
// The ID field must be provided
func (a Alerts) Get(alert *Alert) error {
	if *alert.ID == "" {
		return fmt.Errorf("Alert id field is not set")
	}

	return a.crudAlert("GET", fmt.Sprintf("%s/%s", baseAlertPath, *alert.ID), alert)
}

// Find returns all alerts filtered by the given search conditions.
// If filter is nil, all alerts are returned.
func (a Alerts) Find(filter []*SearchCondition) ([]*Alert, error) {
	search := &Search{
		client: a.client,
		Type:   "alert",
		Params: &SearchParams{
			Conditions: filter,
		},
	}

	var results []*Alert
	moreItems := true
	for moreItems == true {
		resp, err := search.Execute()
		if err != nil {
			return nil, err
		}
		var tmpres []*Alert
		err = json.Unmarshal(resp.Response.Items, &tmpres)
		if err != nil {
			return nil, err
		}
		results = append(results, tmpres...)
		moreItems = resp.Response.MoreItems
		search.Params.Offset = resp.NextOffset
	}

	return results, nil
}

// Create is used to create an Alert in Wavefront.
// If successful, the ID field of the alert will be populated.
func (a Alerts) Create(alert *Alert) error {
	return a.crudAlert("POST", baseAlertPath, alert)
}

// Update is used to update an existing Alert.
// The ID field of the alert must be populated
func (a Alerts) Update(alert *Alert) error {
	if alert.ID == nil {
		return fmt.Errorf("alert id field not set")
	}

	return a.crudAlert("PUT", fmt.Sprintf("%s/%s", baseAlertPath, *alert.ID), alert)

}

// Delete is used to delete an existing Alert.
// The ID field of the alert must be populated
func (a Alerts) Delete(alert *Alert) error {
	if alert.ID == nil {
		return fmt.Errorf("alert id field not set")
	}

	err := a.crudAlert("DELETE", fmt.Sprintf("%s/%s", baseAlertPath, *alert.ID), alert)
	if err != nil {
		return err
	}

	//reset the ID field so deletion is not attempted again
	alert.ID = nil
	return nil

}

func (a Alerts) crudAlert(method, path string, alert *Alert) error {
	payload, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	req, err := a.client.NewRequest(method, path, nil, payload)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &struct {
		Response *Alert `json:"response"`
	}{
		Response: alert,
	})

}

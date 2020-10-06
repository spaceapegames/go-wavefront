package wavefront

import (
	"encoding/json"
	"fmt"
	"strconv"
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

	// Access Control Lists for who can view or modify (user or group IDs)// ACL to apply to the alert
	ACL AccessControlList `json:"acl"`

	FailingHostLabelPairs       []SourceLabelPair `json:"failingHostLabelPairs,omitempty"`
	InMaintenanceHostLabelPairs []SourceLabelPair `json:"inMaintenanceHostLabelPairs,omitempty"`

	// The interval between checks for this alert, in minutes
	CheckingFrequencyInMinutes int `json:"processRateMinutes,omitempty"`

	// Real-Time Alerting, evaluate the alert strictly on ingested data without accounting for delays
	EvaluateRealtimeData bool `json:"evaluateRealtimeData,omitempty"`

	// Include obsolete metrics in alert query
	IncludeObsoleteMetrics bool `json:"includeObsoleteMetrics,omitempty"`
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
		return fmt.Errorf("alert id field is not set")
	}

	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", baseAlertPath, *alert.ID),
		a.client,
		doResponse(alert))
}

// Find returns all alerts filtered by the given search conditions.
// If filter is nil, all alerts are returned.
func (a Alerts) Find(filter []*SearchCondition) (alerts []*Alert, err error) {
	err = doSearch(filter, "alert", a.client, &alerts)
	return
}

// Create is used to create an Alert in Wavefront.
// If successful, the ID field of the alert will be populated.
func (a Alerts) Create(alert *Alert) error {
	return doRest(
		"POST",
		baseAlertPath,
		a.client,
		doPayload(alert),
		doResponse(alert))
}

// Update is used to update an existing Alert.
// The ID field of the alert must be populated
func (a Alerts) Update(alert *Alert) error {
	if alert.ID == nil {
		return fmt.Errorf("alert id field not set")
	}

	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", baseAlertPath, *alert.ID),
		a.client,
		doPayload(alert),
		doResponse(alert))
}

// Delete is used to delete an existing Alert.
// The ID field of the alert must be populated
func (a Alerts) Delete(alert *Alert, skipTrash bool) error {
	if alert.ID == nil {
		return fmt.Errorf("alert id field not set")
	}

	params := map[string]string{
		"skipTrash": strconv.FormatBool(skipTrash),
	}

	err := doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", baseAlertPath, *alert.ID),
		a.client,
		doParams(params))
	if err != nil {
		return err
	}

	//reset the ID field so deletion is not attempted again
	alert.ID = nil
	return nil

}

// Sets the ACL on the alert with the supplied list of IDs for canView and canModify
// an empty []string on canView will remove all values set
// an empty []string on canModify will set the value to the owner of the token issuing the API call
func (a Alerts) SetACL(id string, canView, canModify []string) error {
	return putEntityACL(id, canView, canModify, baseAlertPath, a.client)
}

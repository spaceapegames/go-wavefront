package wavefront

import (
	"fmt"
)

// Target represents a Wavefront Alert Target, for routing notifications
// associated with Alerts.
// Targets can be either email or webhook targets, and the Method must be set
// appropriately.
type Target struct {
	// Description is a description of the target Target
	Description string `json:"description"`

	// ID is the Wavefront-assigned ID of an existing Target
	ID *string `json:"id"`

	// Template is the Mustache template for the notification body
	Template string `json:"template"`

	// Title is the title(name) of the Target
	Title string `json:"title"`

	// Method must be EMAIL, WEBHOOK or PAGERDUTY
	Method string `json:"method"`

	// List of routing targets that this notificant will notify
	Routes []AlertRoute `json:"routes"`

	// Recipient is a comma-separated list of email addresses, webhook URL,
	// or 32-digit PagerDuty key  to which notifications will be sent for this Target
	Recipient string `json:"recipient"`

	// EmailSubject is the subject of the email which will be sent for this Target
	// (EMAIL targets only)
	EmailSubject string `json:"emailSubject"`

	// IsHTMLContent is a boolean value for wavefront to add HTML Boilerplate
	// while using HTML Templates as email.
	// (EMAIL targets only)
	IsHtmlContent bool `json:"isHtmlContent"`

	// ContentType is the content type for webhook posts (e.g. application/json)
	// (WEBHOOK targets only)
	ContentType string `json:"contentType"`

	// CustomHeaders are any custom HTTP headers that should be sent with webhook,
	// in key:value syntax (WEBHOOK targets only)
	CustomHeaders map[string]string `json:"customHttpHeaders"`

	// Triggers is a list of Alert states that will trigger this notification
	// and can include ALERT_OPENED, ALERT_RESOLVED, ALERT_STATUS_RESOLVED,
	// ALERT_AFFECTED_BY_MAINTENANCE_WINDOW, ALERT_SNOOZED, ALERT_NO_DATA,
	// ALERT_NO_DATA_RESOLVED
	Triggers []string `json:"triggers"`
}

type AlertRoute struct {
	// The notification target method, supports values in EMAIL, PAGERDUTY, WEBHOOK
	Method string `json:"method,omitempty"`

	// The endpoint for the notification target
	Target string `json:"target,omitempty"`

	// A space delimited string to filter on tagk=value e.g. env prod
	Filter string `json:"filter,omitempty"`
}

// Targets is used to perform target-related operations against the Wavefront API
type Targets struct {
	// client is the Wavefront client used to perform target-related operations
	client Wavefronter
}

const baseTargetPath = "/api/v2/notificant"

// Targets is used to return a client for target-related operations
func (c *Client) Targets() *Targets {
	return &Targets{client: c}
}

// Get is used to retrieve an existing Target by ID.
// The ID field must be provided
func (t Targets) Get(target *Target) error {
	if target.ID == nil || *target.ID == "" {
		return fmt.Errorf("target id field is not set")
	}
	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", baseTargetPath, *target.ID),
		t.client,
		doResponse(target))
}

// Find returns all targets filtered by the given search conditions.
// If filter is nil, all targets are returned.
func (t Targets) Find(filter []*SearchCondition) (
	results []*Target, err error) {
	err = doSearch(filter, "notificant", t.client, &results)
	return
}

// Create is used to create a Target in Wavefront.
// If successful, the ID field of the target will be populated.
func (t Targets) Create(target *Target) error {
	return doRest(
		"POST",
		baseTargetPath,
		t.client,
		doPayload(target),
		doResponse(target))
}

// Update is used to update an existing Target.
// The ID field of the target must be populated
func (t Targets) Update(target *Target) error {
	if target.ID == nil {
		return fmt.Errorf("target id field not set")
	}

	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", baseTargetPath, *target.ID),
		t.client,
		doPayload(target),
		doResponse(target))
}

// Delete is used to delete an existing Target.
// The ID field of the target must be populated
func (t Targets) Delete(target *Target) error {
	if target.ID == nil {
		return fmt.Errorf("target id field not set")
	}

	err := doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", baseTargetPath, *target.ID),
		t.client)
	if err != nil {
		return err
	}

	//reset the ID field so deletion is not attempted again
	target.ID = nil
	return nil

}

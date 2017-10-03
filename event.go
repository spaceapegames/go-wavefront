package wavefront

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// Event represents a single Wavefront Event
type Event struct {
	// Name is the name given to the Event
	Name string `json:"name"`

	// ID is the Wavefront-assigned ID of an existing Event
	ID *string `json:"id,omitempty"`

	// StartTime is the start time, in epoch milliseconds, of the Event.
	// If zero, it will be set to current time
	StartTime int64 `json:"startTime"`

	// EndTime is the end time, in epoch milliseconds, of the Event
	EndTime int64 `json:"endTime"`

	// Tags are the tags associated with the Event
	Tags []string `json:"tags"`

	// Severity is the severity category of the Event, can be INFO, WARN,
	// SEVERE or UNCLASSIFIED
	Severity string

	// Type is the type of the Event, e.g. "Alert", "Deploy" etc.
	Type string

	// Details is a description of the Event
	Details string

	// Instantaneous, if true, creates a point-in-time Event (i.e. with no duration)
	Instantaneous bool `json:"isEphemeral"`
}

// Events is used to perform event-related operations against the Wavefront API
type Events struct {
	// client is the Wavefront client used to perform event-related operations
	client Wavefronter
}

const baseEventPath = "/api/v2/event"

// UnmarshalJSON is a custom JSON unmarshaller for an Event, used to explode
// the annotations.
func (e *Event) UnmarshalJSON(b []byte) error {
	type event Event
	temp := struct {
		Annotations map[string]string `json:"annotations"`
		*event
	}{
		event: (*event)(e),
	}

	e.Severity = temp.Annotations["severity"]
	e.Type = temp.Annotations["type"]
	e.Details = temp.Annotations["details"]
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	return nil
}

func (e *Event) MarshalJSON() ([]byte, error) {
	type event Event
	return json.Marshal(&struct {
		Annotations map[string]string `json:"annotations"`
		*event
	}{
		Annotations: map[string]string{
			"severity": e.Severity,
			"details":  e.Details,
			"type":     e.Type,
		},
		event: (*event)(e),
	})
}

// Events is used to return a client for event-related operations
func (c *Client) Events() *Events {
	return &Events{client: c}
}

// Find returns all events filtered by the given search conditions.
// If filter is nil then all Events are returned. The result set is limited to
// the first 100 entries. If more results are required the Search type can
// be used directly.
func (e Events) Find(filter []*SearchCondition, timeRange *TimeRange) ([]*Event, error) {
	search := &Search{
		client: e.client,
		Type:   "event",
		Params: &SearchParams{
			Conditions: filter,
			TimeRange:  timeRange,
		},
	}
	var results []*Event
	resp, err := search.Execute()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp.Response.Items, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// FindByID returns the Error with the Wavefront-assigned ID.
// If not found an error is returned
func (e Events) FindByID(id string) (*Event, error) {
	res, err := e.Find([]*SearchCondition{
		&SearchCondition{
			Key:            "id",
			Value:          id,
			MatchingMethod: "EXACT",
		},
	}, nil)

	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("no event found with ID %s", id)
	}

	return res[0], nil
}

// Create is used to create an Event in Wavefront.
// If successful, the ID field of the event will be populated.
func (a Events) Create(event *Event) error {
	if event.StartTime == 0 {
		event.StartTime = time.Now().Unix() * 1000
	}
	if event.Instantaneous == true {
		event.EndTime = event.StartTime + 1
	}
	return a.crudEvent("POST", baseEventPath, event)
}

// Update is used to update an existing Event.
// The ID field of the Event must be populated
func (e Events) Update(event *Event) error {
	if event.ID == nil {
		return fmt.Errorf("Event id field not set")
	}

	return e.crudEvent("PUT", fmt.Sprintf("%s/%s", baseEventPath, *event.ID), event)

}

// Close is used to close an existing Event
func (e Events) Close(event *Event) error {
	if event.ID == nil {
		return fmt.Errorf("Event id field not set")
	}

	return e.crudEvent("POST", fmt.Sprintf("%s/%s/close", baseEventPath, *event.ID), event)
}

// Delete is used to delete an existing Event.
// The ID field of the Event must be populated
func (e Events) Delete(event *Event) error {
	if event.ID == nil {
		return fmt.Errorf("Event id field not set")
	}

	err := e.crudEvent("DELETE", fmt.Sprintf("%s/%s", baseEventPath, *event.ID), event)
	if err != nil {
		return err
	}

	//reset the ID field so deletion is not attempted again
	event.ID = nil
	return nil

}

func (e Events) crudEvent(method, path string, event *Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	req, err := e.client.NewRequest(method, path, nil, payload)
	if err != nil {
		return err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &struct {
		Response *Event `json:"response"`
	}{
		Response: event,
	})
}

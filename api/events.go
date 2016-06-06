package wavefront

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// WavefrontEvent type is used to store events that will be posted to Wavefront
// The 'param' struct tags are used to convert the data into a set of
// request parameters suitable for WF events
type WavefrontEvent struct {
	Name          string   `param:"n"` // Name for the event
	Start         int64    `param:"s"` // Start time of the event
	End           int64    `param:"e"` // End time of the event
	Instantaneous bool     `param:"c"` // Flag for instantaneous event
	Details       string   `param:"d"` // any additional details for the event
	Hosts         []string `param:"h"` // Hostname(s) affected by the event
	Severity      string   `param:"l"` // User-defined severity
	Type          string   `param:"t"` // Event type
}

// Event is used to create Wavefront events
type Event struct {
	client *Client
	event  WavefrontEvent
}

// baseEventPath is the base API path for retrieving alerts
const baseEventPath = "/api/events"

// NewEvent returns an Event struct with default
// fields sensibly pre-filled
func (e *Event) NewEvent(name string) *Event {
	e.event = WavefrontEvent{Name: name, Start: time.Now().UnixNano() / 1000000}
	return e
}

// params converts an Event struct into a QueryParams type
func (e Event) params() *QueryParams {
	params := QueryParams{}
	t := reflect.TypeOf(e.event)
	v := reflect.ValueOf(e.event)
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("param")
		name := t.Field(i).Name

		// deal with affected hosts array separately
		if name == "Hosts" && len(e.event.Hosts) > 0 {
			join := fmt.Sprintf("&%s=", tag)
			params[tag] = strings.Join(e.event.Hosts, join)
		}

		switch v.Field(i).Kind() {
		case reflect.String:
			if v.Field(i).String() != "" {
				params[tag] = v.FieldByName(name).String()
			}
		case reflect.Int64:
			if v.Field(i).Int() != 0 {
				params[tag] = strconv.FormatInt(v.FieldByName(name).Int(), 10)
			}
		case reflect.Bool:
			params[tag] = strconv.FormatBool(v.FieldByName(name).Bool())
		}
	}

	return &params
}

// post uses the Wavefront client to create/delete Events
func (e Event) post(path string, method string) ([]byte, error) {
	req, err := e.client.NewRequest(method, path, e.params())
	if err != nil {
		return nil, err
	}
	resp, err := e.client.Do(req, new(interface{}))
	if err != nil {
		return nil, err
	}

	raw, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// Create posts the Event to Wavefront
func (e Event) Create() ([]byte, error) {
	return e.post(baseEventPath, "POST")
}

// Instant creates an Instantaneous event
func (e Event) Instant() ([]byte, error) {
	e.event.Instantaneous = true
	return e.Create()
}

// End is used to end a (non-instantaneous) event
func (e Event) End() ([]byte, error) {
	return e.post(baseEventPath+"/close", "POST")
}

// Delete is used to delete an event
func (e Event) Delete() ([]byte, error) {
	p := *e.params()
	path, err := url.Parse(fmt.Sprintf("%s/%s/%s", baseEventPath, p["s"], p["n"]))
	if err != nil {
		return nil, err
	}
	return e.post(path.String(), "DELETE")
}

// AddAffectedHost adds a hostname to the list of those affected by this Event
func (e *Event) AddAffectedHost(hostname string) {
	e.event.Hosts = append(e.event.Hosts, hostname)
}

// SetStartTime should be used when acting on existing Events
// (i.e. Ending or Deleting them) for which the start time is known
func (e *Event) SetStartTime(start int64) {
	e.event.Start = start
}

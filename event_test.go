package wavefront

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type MockEventClient struct {
	Client
	T *testing.T
}

type MockCrudEventClient struct {
	Client
	T      *testing.T
	method string
}

func (m *MockEventClient) Do(req *http.Request) (io.ReadCloser, error) {
	body, _ := ioutil.ReadAll(req.Body)
	search := SearchParams{}
	err := json.Unmarshal(body, &search)
	if err != nil {
		m.T.Fatal(err)
	}
	if search.TimeRange.EndTime != 1498723080000 && search.TimeRange.StartTime != 1498719480000 {
		m.T.Errorf("expected time range 1498719480000 - 1498723080000, got %d - %d",
			search.TimeRange.StartTime, search.TimeRange.EndTime,
		)
	}
	response, err := ioutil.ReadFile("./fixtures/search-event-response.json")
	if err != nil {
		m.T.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

func TestEvents_Find(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	e := &Events{
		client: &MockEventClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
	tr, _ := NewTimeRange(1498723080, LastHour)
	events, err := e.Find(nil, tr)

	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 3 {
		t.Errorf("expected 3 events, got %d", len(events))
	}

	if *events[0].ID != "1498664617084:Alert Fired: Service Errors" {
		t.Errorf("expected first event ID '1498664617084:Alert Fired: Service Errors', got %s", *events[0].ID)
	}

	if events[0].Severity != "warn" || events[0].Type != "alert-detail" || events[0].Details != "some details" {
		t.Errorf("unexpected annotations on event")
	}

}

func (m *MockCrudEventClient) Do(req *http.Request) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile("./fixtures/create-event-response.json")
	if err != nil {
		m.T.Fatal(err)
	}
	if req.Method != m.method {
		m.T.Errorf("request method expected '%s' got '%s'", m.method, req.Method)
	}
	body, _ := ioutil.ReadAll(req.Body)
	event := Event{}
	err = json.Unmarshal(body, &event)
	if err != nil {
		m.T.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

func TestEvents_CreateUpdateDeleteEvent(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	e := &Events{
		client: &MockCrudEventClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			method: "PUT",
			T:      t,
		},
	}

	event := Event{
		Name:      "test event",
		StartTime: time.Now().Unix() * 1000,
		Tags:      []string{"mytag1"},
		Severity:  "warn",
	}

	if err := e.Update(&event); err == nil {
		t.Errorf("expected event update to error with no ID")
	}

	e.client.(*MockCrudEventClient).method = "POST"

	e.Create(&event)
	if *event.ID != "1234" {
		t.Errorf("event ID expected 1234, got %s", *event.ID)
	}

	e.client.(*MockCrudEventClient).method = "PUT"
	if err := e.Update(&event); err != nil {
		t.Error(err)
	}

	e.client.(*MockCrudEventClient).method = "POST"
	if err := e.Close(&event); err != nil {
		t.Error(err)
	}

	e.client.(*MockCrudEventClient).method = "DELETE"
	if err := e.Delete(&event); err != nil {
		t.Error(err)
	}

	if event.ID != nil {
		t.Errorf("expected event ID to be reset after deletion")
	}

}

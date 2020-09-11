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

	asserts "github.com/stretchr/testify/assert"
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
	search := &SearchParams{}
	resp, err := testDo(m.T, req, "./fixtures/search-event-response.json", "POST", search)

	assertEqual(m.T, int64(1498719480000), search.TimeRange.StartTime)
	assertEqual(m.T, int64(1498723080000), search.TimeRange.EndTime)
	return resp, err
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

	assertEqual(t, 3, len(events))
	assertEqual(t, "1498664617084:Alert Fired: Service Errors", *events[0].ID)
	assertEqual(t, "warn", events[0].Severity)
	assertEqual(t, "alert-detail", events[0].Type)
	assertEqual(t, "some details", events[0].Details)
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
	assert := asserts.New(t)
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

	// Update should fail because no ID is set
	assert.Error(e.Update(&event))

	e.client.(*MockCrudEventClient).method = "POST"

	assert.NoError(e.Create(&event))
	assert.Equal("1234", *event.ID)

	e.client.(*MockCrudEventClient).method = "PUT"
	assert.NoError(e.Update(&event))

	e.client.(*MockCrudEventClient).method = "POST"
	assert.NoError(e.Close(&event))

	e.client.(*MockCrudEventClient).method = "DELETE"
	assert.NoError(e.Delete(&event))

	assert.Nil(event.ID)
}

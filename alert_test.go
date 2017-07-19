package wavefront

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type MockAlertClient struct {
	Client
	InvokedCount int
	T            *testing.T
}

type MockCrudAlertClient struct {
	Client
	method string
	T      *testing.T
}

func (m *MockAlertClient) Do(req *http.Request) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile(fmt.Sprintf("./fixtures/paginated-alert-%d.json", m.InvokedCount))
	if err != nil {
		m.T.Fatal(err)
	}
	body, _ := ioutil.ReadAll(req.Body)
	search := SearchParams{}
	err = json.Unmarshal(body, &search)
	if err != nil {
		m.T.Fatal(err)
	}
	if search.Offset != search.Limit*m.InvokedCount {
		m.T.Errorf("offset, expected %d, got %d", search.Limit*m.InvokedCount, search.Offset)
	}
	m.InvokedCount++
	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

func TestAlerts_PaginatedFind(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	a := &Alerts{
		client: &MockAlertClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
	alerts, err := a.Find(nil)
	if err != nil {
		t.Fatal(err)
	}
	invoked := ((a.client).(*MockAlertClient)).InvokedCount
	if invoked != 2 {
		t.Errorf("paginated search, expected 2, got %d", invoked)
	}

	if alerts[0].Name != "Excessive consumption of inodes" {
		t.Errorf("alert name incorrect: %s", alerts[0].Name)
	}

	if len(alerts[0].Tags) != 2 {
		t.Errorf("alert tags, expected 2, got %d", len(alerts[0].Tags))
	}

}

func (m *MockCrudAlertClient) Do(req *http.Request) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile("./fixtures/create-alert-response.json")
	if err != nil {
		m.T.Fatal(err)
	}
	if req.Method != m.method {
		m.T.Errorf("request method expected '%s' got '%s'", m.method, req.Method)
	}
	body, _ := ioutil.ReadAll(req.Body)
	alert := Alert{}
	err = json.Unmarshal(body, &alert)
	if err != nil {
		m.T.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

func TestAlerts_CreateUpdateDeleteAlert(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	a := &Alerts{
		client: &MockCrudAlertClient{
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

	alert := Alert{
		Name:                "test alert",
		Target:              "test@example.com",
		Condition:           "ts(servers.cpu.usage) > 10 * 10",
		DisplayExpression:   "ts(servers.cpu.usage)",
		Minutes:             2,
		ResolveAfterMinutes: 2,
		Severity:            "WARN",
		Tags:                []string{"mytag1", "mytag2"},
	}

	if err := a.Update(&alert); err == nil {
		t.Errorf("expected alert update to error with no ID")
	}

	a.client.(*MockCrudAlertClient).method = "POST"

	a.Create(&alert)
	if *alert.ID != "1234" {
		t.Errorf("alert ID expected 1234, got %s", *alert.ID)
	}

	a.client.(*MockCrudAlertClient).method = "PUT"
	if err := a.Update(&alert); err != nil {
		t.Error(err)
	}

	a.client.(*MockCrudAlertClient).method = "DELETE"
	if err := a.Delete(&alert); err != nil {
		t.Error(err)
	}

	if alert.ID != nil {
		t.Error("expected alert ID to be reset after deletion")
	}

}

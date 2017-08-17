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

type MockDashboardClient struct {
	Client
	InvokedCount int
	T            *testing.T
}

type MockCrudDashboardClient struct {
	Client
	method string
	T      *testing.T
}

func (m *MockDashboardClient) Do(req *http.Request) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile(fmt.Sprintf("./fixtures/paginated-dashboard-%d.json", m.InvokedCount))
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

func TestDashboards_PaginatedFind(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	a := &Dashboards{
		client: &MockDashboardClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
	dashboards, err := a.Find(nil)
	if err != nil {
		t.Fatal(err)
	}
	invoked := ((a.client).(*MockDashboardClient)).InvokedCount
	if invoked != 2 {
		t.Errorf("paginated search, expected 2, got %d", invoked)
	}

	if dashboards[0].Name != "test 1" {
		t.Errorf("dashboard name incorrect: %s", dashboards[0].Name)
	}

	if len(dashboards[0].Tags) != 2 {
		t.Errorf("dashboard tags, expected 2, got %d", len(dashboards[0].Tags))
	}

}

func (m *MockCrudDashboardClient) Do(req *http.Request) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile("./fixtures/create-dashboard-response.json")
	if err != nil {
		m.T.Fatal(err)
	}
	if req.Method != m.method {
		m.T.Errorf("request method expected '%s' got '%s'", m.method, req.Method)
	}
	body, _ := ioutil.ReadAll(req.Body)
	dashboard := Dashboard{}
	err = json.Unmarshal(body, &dashboard)
	if err != nil {
		m.T.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

func TestDashboards_CreateUpdateDeleteDashboard(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	a := &Dashboards{
		client: &MockCrudDashboardClient{
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

	dashboard := Dashboard{
		Name:        "Dashboard API example",
		Tags:        []string{"test"},
		Description: "Dashboard Description",
		Url:         "api-example",
	}

	if err := a.Update(&dashboard); err == nil {
		t.Errorf("expected dashboard update to error with no ID")
	}

	a.client.(*MockCrudDashboardClient).method = "POST"

	a.Create(&dashboard)
	if dashboard.ID != "api-example" {
		t.Errorf("dashboard ID expected api-example, got %s", dashboard.ID)
	}

	a.client.(*MockCrudDashboardClient).method = "PUT"
	if err := a.Update(&dashboard); err != nil {
		t.Error(err)
	}

	a.client.(*MockCrudDashboardClient).method = "DELETE"
	if err := a.Delete(&dashboard); err != nil {
		t.Error(err)
	}

	if dashboard.ID != "" {
		t.Error("expected dashboard ID to be reset after deletion")
	}

}

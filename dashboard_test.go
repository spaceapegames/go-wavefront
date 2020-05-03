package wavefront

import (
	"io"
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
	return testPaginatedDo(m.T, req, "./fixtures/paginated-dashboard-%d.json", &m.InvokedCount)
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
	return testDo(m.T, req, "./fixtures/create-dashboard-response.json", m.method, Dashboard{})
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
	assertEqual(t, "api-example", dashboard.ID)

	a.client.(*MockCrudDashboardClient).method = "PUT"
	if err := a.Update(&dashboard); err != nil {
		t.Error(err)
	}

	a.client.(*MockCrudDashboardClient).method = "DELETE"
	if err := a.Delete(&dashboard, true); err != nil {
		t.Error(err)
	}

	assertEqual(t, "", dashboard.ID)
}

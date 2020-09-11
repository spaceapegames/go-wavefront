package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	asserts "github.com/stretchr/testify/assert"
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
	assert := asserts.New(t)
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
	assert.Equal(2, invoked)
	assert.Equal("test 1", dashboards[0].Name)
	assert.Len(dashboards[0].Tags, 2)
}

func (m *MockCrudDashboardClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/create-dashboard-response.json", m.method, &Dashboard{})
}

func TestDashboards_CreateUpdateDeleteDashboard(t *testing.T) {
	assert := asserts.New(t)
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

	// Update should fail because no ID is set
	assert.Error(a.Update(&dashboard))

	a.client.(*MockCrudDashboardClient).method = "POST"

	assert.NoError(a.Create(&dashboard))
	assert.Equal("api-example", dashboard.ID)

	a.client.(*MockCrudDashboardClient).method = "PUT"
	assert.NoError(a.Update(&dashboard))

	a.client.(*MockCrudDashboardClient).method = "DELETE"
	assert.NoError(a.Delete(&dashboard, true))

	assert.Equal("", dashboard.ID)
}

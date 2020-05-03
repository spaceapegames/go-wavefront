package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockDerivedMetricsClient struct {
	Client
	T *testing.T
}

type MockCrudDerivedMetricsClient struct {
	Client
	T      *testing.T
	method string
}

func (m *MockDerivedMetricsClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/search-derivedmetrics-response.json", "POST", &SearchParams{})
}

func (m *MockCrudDerivedMetricsClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/crud-derivedmetric-response.json", m.method, &DerivedMetric{})
}

func TestDerivedMetrics_Find(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	d := &DerivedMetrics{
		client: &MockDerivedMetricsClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	derivedMetrics, err := d.Find(nil)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 1, len(derivedMetrics))
	assertEqual(t, "1234567891011", *derivedMetrics[0].ID)
}

func TestDerivedMetrics_CRUD(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	d := &DerivedMetrics{
		client: &MockCrudDerivedMetricsClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	derivedMetric := &DerivedMetric{
		Name: "example",
	}
	d.client.(*MockCrudDerivedMetricsClient).method = "POST"
	if err := d.Create(derivedMetric); err == nil {
		t.Errorf("expected an error on missing name, query, and/or minutes")
	}

	derivedMetric.Query = "ts(cpu.usage.*)"
	derivedMetric.Minutes = 10

	var _ = d.Create(derivedMetric)
	assertEqual(t, "1234567891011", *derivedMetric.ID)

	d.client.(*MockCrudDerivedMetricsClient).method = "GET"
	_ = d.Get(derivedMetric)
	assertEqual(t, "example", derivedMetric.Name)

	d.client.(*MockCrudDerivedMetricsClient).method = "PUT"
	_ = d.Update(derivedMetric)

	d.client.(*MockCrudDerivedMetricsClient).method = "DELETE"
	_ = d.Delete(derivedMetric, true)
	assertEqual(t, "", *derivedMetric.ID)
}

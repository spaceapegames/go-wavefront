package wavefront

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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
	body, _ := ioutil.ReadAll(req.Body)
	search := SearchParams{}
	err := json.Unmarshal(body, &search)
	if err != nil {
		m.T.Fatal(err)
	}

	response, err := ioutil.ReadFile("./fixtures/search-derivedmetrics-response.json")
	if err != nil {
		m.T.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

func (m *MockCrudDerivedMetricsClient) Do(req *http.Request) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile("./fixtures/crud-derivedmetric-response.json")
	if err != nil {
		m.T.Fatal(err)
	}
	if req.Method != m.method {
		m.T.Errorf("request method expected '%s' got '%s'", m.method, req.Method)
	}

	body, _ := ioutil.ReadAll(req.Body)
	derivedMetrics := DerivedMetric{}
	err = json.Unmarshal(body, &derivedMetrics)
	if err != nil {
		m.T.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
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

	if len(derivedMetrics) != 1 {
		t.Errorf("expected to find one derived metric, got %d", len(derivedMetrics))
	}

	if *derivedMetrics[0].ID != "1234567891011" {
		t.Errorf("expected first id to be 1234567891011, got %s", "1234567891011")
	}
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
	if *derivedMetric.ID != "1234567891011" {
		t.Errorf("expected id returned from create to be 1234567891011, got %s", "1234567891011")
	}

	d.client.(*MockCrudDerivedMetricsClient).method = "GET"
	_ = d.Get(derivedMetric)
	if derivedMetric.Name != "example" {
		t.Errorf("expected to get derived metric with name example, got %s", derivedMetric.Name)
	}

	d.client.(*MockCrudDerivedMetricsClient).method = "PUT"
	_ = d.Update(derivedMetric)

	d.client.(*MockCrudDerivedMetricsClient).method = "DELETE"
	_ = d.Delete(derivedMetric)
	if *derivedMetric.ID != "" {
		t.Errorf("expected id to be empty, got %s", *derivedMetric.ID)
	}
}

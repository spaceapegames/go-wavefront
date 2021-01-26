package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockIngestionPoliciesClient struct {
	Client
	T *testing.T
}

type MockCrudIngestionPoliciesClient struct {
	Client
	T      *testing.T
	method string
}

func (pol MockIngestionPoliciesClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(pol.T, req, "./fixtures/search-ingestionpolicy-response.json", "POST", &SearchParams{})
}

func (pol MockCrudIngestionPoliciesClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(pol.T, req, "./fixtures/crud-ingestionpolicy-response.json", pol.method, &IngestionPolicy{})
}

func TestIngestionPolicies_Find(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	pol := &IngestionPolicies{
		client: &MockIngestionPoliciesClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	ingestionPolicies, err := pol.Find(nil)

	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 1, len(ingestionPolicies))
	assertEqual(t, "someid", ingestionPolicies[0].ID)
}

func TestIngestionPolicies_CreateUpdateDelete(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	pol := &IngestionPolicies{
		client: &MockCrudIngestionPoliciesClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	pol.client.(*MockCrudIngestionPoliciesClient).method = "POST"

	ingestionPolicy := &IngestionPolicy{}

	if err := pol.Create(ingestionPolicy); err == nil {
		t.Errorf("expected to receive error for missing fields")
	}

	ingestionPolicy.Name = "Example"
	ingestionPolicy.Description = "someDescription"

	if err := pol.Create(ingestionPolicy); err != nil {
		t.Fatal(err)
	}

	pol.client.(*MockCrudIngestionPoliciesClient).method = "GET"
	var _ = pol.Get(ingestionPolicy)

	pol.client.(*MockCrudIngestionPoliciesClient).method = "PUT"
	var _ = pol.Update(ingestionPolicy)

	pol.client.(*MockCrudIngestionPoliciesClient).method = "DELETE"
	var _ = pol.Delete(ingestionPolicy)

	assertEqual(t, "test-policy-1607616352537", ingestionPolicy.ID)
}

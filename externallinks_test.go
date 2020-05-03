package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockExternalLinksClient struct {
	Client
	T *testing.T
}

type MockCrudExternalLinksClient struct {
	Client
	T      *testing.T
	method string
}

func (e MockExternalLinksClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(e.T, req, "./fixtures/search-extlinks-response.json", "POST", &SearchParams{})
}

func (e MockCrudExternalLinksClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(e.T, req, "./fixtures/crud-extlink-response.json", e.method, &ExternalLink{})
}

func TestExternalLinks_Find(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	e := &ExternalLinks{
		client: &MockExternalLinksClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	externalLinks, err := e.Find(nil)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 1, len(externalLinks))
	assertEqual(t, "someid", *externalLinks[0].ID)
}

func TestExternalLinks_CreateUpdateDelete(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	e := &ExternalLinks{
		client: &MockCrudExternalLinksClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	e.client.(*MockCrudExternalLinksClient).method = "POST"

	externalLink := &ExternalLink{}
	if err := e.Create(externalLink); err == nil {
		t.Errorf("expected to receive error for missing fields")
	}

	externalLink.Name = "Example"
	externalLink.Description = "someDescription"
	externalLink.Template = "https://www.someSystem.com/events?logSource={{{source}}}&startTime={{startEpochMillis}}&endTime={{endEpochMillis}}"
	if err := e.Create(externalLink); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, ".*", externalLink.SourceFilterRegex)

	e.client.(*MockCrudExternalLinksClient).method = "GET"
	var _ = e.Get(externalLink)

	e.client.(*MockCrudExternalLinksClient).method = "PUT"
	var _ = e.Update(externalLink)

	e.client.(*MockCrudExternalLinksClient).method = "DELETE"
	var _ = e.Delete(externalLink)

	assertEqual(t, "", *externalLink.ID)
}

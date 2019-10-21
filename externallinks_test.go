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
	body, _ := ioutil.ReadAll(req.Body)
	search := SearchParams{}
	err := json.Unmarshal(body, &search)
	if err != nil {
		e.T.Fatal(err)
	}

	response, err := ioutil.ReadFile("./fixtures/search-extlinks-response.json")
	if err != nil {
		e.T.Fatal(err)
	}

	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

func (e MockCrudExternalLinksClient) Do(req *http.Request) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile("./fixtures/crud-extlink-response.json")
	if err != nil {
		e.T.Fatal(err)
	}

	if req.Method != e.method {
		e.T.Errorf("request method expected '%s' got '%s'", e.method, req.Method)
	}

	body, _ := ioutil.ReadAll(req.Body)
	link := ExternalLink{}
	err = json.Unmarshal(body, &link)
	if err != nil {
		e.T.Fatal(err)
	}

	return ioutil.NopCloser(bytes.NewReader(response)), nil
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

	if len(externalLinks) != 1 {
		t.Errorf("expected one ExternalLink returned, got %d", len(externalLinks))
	}

	if *externalLinks[0].ID != "someid" {
		t.Errorf("expected first ExternalLink id to be someid, got %s", *externalLinks[0].ID)
	}
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

	if externalLink.SourceFilterRegex != ".*" {
		t.Errorf("expected to find a Source Filter Regex of '*', got %s", externalLink.SourceFilterRegex)
	}

	e.client.(*MockCrudExternalLinksClient).method = "GET"
	var _ = e.Get(externalLink)

	e.client.(*MockCrudExternalLinksClient).method = "PUT"
	var _ = e.Update(externalLink)

	e.client.(*MockCrudExternalLinksClient).method = "DELETE"
	var _ = e.Delete(externalLink)

	if *externalLink.ID != "" {
		t.Errorf("exected ExternalLink ID to be blank, got %s", *externalLink.ID)
	}

}

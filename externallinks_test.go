package wavefront

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
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
		e.T.Errorf("request method expected '%s' got '%s'", m.method, req.Method)
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

}

func TestExternalLinks_CreateUpdateDelete(t *testing.T) {

}

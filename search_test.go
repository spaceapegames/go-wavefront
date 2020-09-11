package wavefront

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	asserts "github.com/stretchr/testify/assert"
)

type MockSearchClient struct {
	Client
	Response  []byte
	T         *testing.T
	isDeleted bool
}

func (m MockSearchClient) Do(req *http.Request) (io.ReadCloser, error) {
	p := SearchParams{}
	b, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(b, &p)
	if err != nil {
		m.T.Fatal(err)
	}
	// check defaults
	if p.Offset != 0 || p.Limit != 100 {
		m.T.Errorf("default offset and limit, expected 0, 100; got %d, %d", p.Offset, p.Limit)
	}

	if m.isDeleted == true && req.URL.Path != "/api/v2/search/alert/deleted" {
		m.T.Errorf("deleted search path expected /api/v2/search/alert/deleted, got %s", req.URL.Path)
	}

	return ioutil.NopCloser(bytes.NewReader(m.Response)), nil
}

func TestSearch(t *testing.T) {
	assert := asserts.New(t)
	sc := &SearchCondition{
		Key:            "tags",
		Value:          "myTag",
		MatchingMethod: "EXACT",
	}

	sp := &SearchParams{
		Conditions: []*SearchCondition{sc},
	}
	response, err := ioutil.ReadFile("./fixtures/search-alert-response.json")
	if err != nil {
		t.Fatal(err)
	}
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	s := &Search{
		Params: sp,
		Type:   "alert",
		client: &MockSearchClient{
			Response: response,
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	resp, err := s.Execute()
	if err != nil {
		t.Fatal("error executing query:", err)
	}

	raw, err := ioutil.ReadAll(resp.RawResponse)
	if err != nil {
		t.Error(err)
	}

	if err := json.Unmarshal(raw, new(map[string]interface{})); err != nil {
		t.Error("raw response is invalid JSON", err)
	}

	// check offset of next page in paginated response
	if resp.NextOffset != 100 {
		t.Errorf("next offset, expected 100, got %d", resp.NextOffset)
	}

	// check deleted path appended
	s.Deleted = true
	((s.client).(*MockSearchClient)).isDeleted = true
	_, err = s.Execute()
	assert.NoError(err)
}

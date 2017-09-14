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

type MockTargetClient struct {
	Client
	T *testing.T
}

type MockCrudTargetClient struct {
	Client
	method string
	T      *testing.T
}

func (m *MockTargetClient) Do(req *http.Request) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile("./fixtures/list-targets.json")
	if err != nil {
		m.T.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

func TestTargets_Find(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	tgts := &Targets{
		client: &MockTargetClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
	targets, err := tgts.Find(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(targets) != 2 {
		t.Errorf("list target response, expected 2, got %d", len(targets))
	}

	if targets[0].Method != "WEBHOOK" {
		t.Errorf("expected target method to be WEBHOOK, got %s", targets[0].Method)
	}
}

func (m *MockCrudTargetClient) Do(req *http.Request) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile("./fixtures/create-target-response.json")
	if err != nil {
		m.T.Fatal(err)
	}
	if req.Method != m.method {
		m.T.Errorf("request method expected '%s' got '%s'", m.method, req.Method)
	}
	body, _ := ioutil.ReadAll(req.Body)
	target := Target{}
	err = json.Unmarshal(body, &target)
	if err != nil {
		m.T.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

func TestTargets_CreateUpdateDeleteTarget(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	a := &Targets{
		client: &MockCrudTargetClient{
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

	tmpl, _ := ioutil.ReadFile("./target-template.tmpl")

	target := Target{
		Title:       "test target",
		Description: "testing something",
		Method:      "WEBHOOK",
		Recipient:   "https://hooks.slack.com/services/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",
		ContentType: "application/json",
		CustomHeaders: map[string]string{
			"Testing": "true",
		},
		Triggers: []string{"ALERT_OPENED", "ALERT_RESOLVED"},
		Template: string(tmpl),
	}

	if err := a.Update(&target); err == nil {
		t.Errorf("expected target update to error with no ID")
	}

	a.client.(*MockCrudTargetClient).method = "POST"

	a.Create(&target)
	if *target.ID != "7" {
		t.Errorf("target ID expected 7, got %s", *target.ID)
	}

	a.client.(*MockCrudTargetClient).method = "PUT"
	if err := a.Update(&target); err != nil {
		t.Error(err)
	}

	a.client.(*MockCrudTargetClient).method = "DELETE"
	if err := a.Delete(&target); err != nil {
		t.Error(err)
	}

	if target.ID != nil {
		t.Error("expected target ID to be reset after deletion")
	}

}

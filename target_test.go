package wavefront

import (
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
	return testDo(m.T, req, "./fixtures/list-targets.json", "POST", &SearchParams{})
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

	assertEqual(t, 2, len(targets))
	assertEqual(t, "WEBHOOK", targets[0].Method)
}

func (m *MockCrudTargetClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/create-target-response.json", m.method, &Target{})
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
		Recipient:   "https://hooks.slack.com/services/test/me",
		Routes: []AlertRoute{
			{
				Method: "WEBHOOK",
				Target: "https://hooks.slack.com/services/test/me",
				Filter: "env prod*",
			},
		},
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
	assertEqual(t, "7", *target.ID)

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

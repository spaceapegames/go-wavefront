package wavefront

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	asserts "github.com/stretchr/testify/assert"
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
	assert := asserts.New(t)
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

	// Update should fail as no ID is set
	assert.Error(a.Update(&target))

	a.client.(*MockCrudTargetClient).method = "POST"
	assert.NoError(a.Create(&target))
	assert.Equal("7", *target.ID)

	a.client.(*MockCrudTargetClient).method = "PUT"
	assert.NoError(a.Update(&target))

	a.client.(*MockCrudTargetClient).method = "DELETE"
	assert.NoError(a.Delete(&target))
	assert.Nil(target.ID)
}

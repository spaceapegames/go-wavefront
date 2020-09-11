package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	asserts "github.com/stretchr/testify/assert"
)

type MockAlertClient struct {
	Client
	InvokedCount int
	T            *testing.T
}

type MockCrudAlertClient struct {
	Client
	method string
	T      *testing.T
}

func (m *MockAlertClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testPaginatedDo(m.T, req, "./fixtures/paginated-alert-%d.json", &m.InvokedCount)
}

func TestAlerts_PaginatedFind(t *testing.T) {
	assert := asserts.New(t)
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	a := &Alerts{
		client: &MockAlertClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
	alerts, err := a.Find(nil)
	if err != nil {
		t.Fatal(err)
	}
	invoked := ((a.client).(*MockAlertClient)).InvokedCount
	assert.Equal(2, invoked)

	assert.Equal("Excessive consumption of inodes", alerts[0].Name)
	assert.Len(alerts[0].Tags, 2)
}

func (m *MockCrudAlertClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/create-alert-response.json", m.method, &Alert{})
}

func TestAlerts_CreateUpdateDeleteAlert(t *testing.T) {
	assert := asserts.New(t)
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	a := &Alerts{
		client: &MockCrudAlertClient{
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

	alert := Alert{
		Name:                "test alert",
		Target:              "test@example.com",
		Condition:           "ts(servers.cpu.usage) > 10 * 10",
		DisplayExpression:   "ts(servers.cpu.usage)",
		Minutes:             2,
		ResolveAfterMinutes: 2,
		Severity:            "WARN",
		AdditionalInfo:      "please resolve this alert",
		Tags:                []string{"mytag1", "mytag2"},
	}

	// Update should fail because no ID is set
	assert.Error(a.Update(&alert))

	a.client.(*MockCrudAlertClient).method = "POST"
	assert.NoError(a.Create(&alert))
	assert.Equal("1234", *alert.ID)

	a.client.(*MockCrudAlertClient).method = "PUT"
	assert.NoError(a.Update(&alert))

	a.client.(*MockCrudAlertClient).method = "DELETE"
	assert.NoError(a.Delete(&alert, true))
	assert.Nil(alert.ID)
}

func TestMultiThresholdAlerts_CreateUpdateDeleteAlert(t *testing.T) {
	assert := asserts.New(t)
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	a := &Alerts{
		client: &MockCrudAlertClient{
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

	alert := Alert{
		Name: "test alert",
		Targets: map[string]string{
			"smoke": "test@example.com",
			"warn":  "test2@example.com",
		},
		Conditions: map[string]string{
			"smoke": "ts(servers.cpu.usage) > 5 * 10",
			"warn":  "ts(servers.cpu.usage) > 10 * 10",
		},
		DisplayExpression:   "ts(servers.cpu.usage)",
		Minutes:             2,
		ResolveAfterMinutes: 2,
		SeverityList:        []string{"SMOKE", "WARN"},
		AdditionalInfo:      "please resolve this alert",
		Tags:                []string{"mytag1", "mytag2"},
	}

	// Update should fail because no ID is set
	assert.Error(a.Update(&alert))

	a.client.(*MockCrudAlertClient).method = "POST"
	assert.NoError(a.Create(&alert))
	assert.Equal("1234", *alert.ID)

	a.client.(*MockCrudAlertClient).method = "PUT"
	assert.NoError(a.Update(&alert))

	a.client.(*MockCrudAlertClient).method = "DELETE"
	assert.NoError(a.Delete(&alert, true))
	assert.Nil(alert.ID)
}

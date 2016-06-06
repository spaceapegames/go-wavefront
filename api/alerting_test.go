package wavefront

import (
	_ "fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var dummyAlert = []byte(`[ {
        "name" : "TEST alert",
        "metricsUsed" : [ "some.metric" ],
        "userTagsWithCounts": {}, 
        "customerTagsWithCounts": {
            "Test": 1
          },
        "severity" : "WARN",
        "hostsUsed" : [ "my.host" ],
        "condition" : "some.condition",
        "event" : { "something" : 1 }
        } ]`)

// test retrieval of alerts
func TestRetrieveAlerts(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(dummyAlert)

		r.ParseForm()
		tag := r.Form.Get("customerTag")
		if tag != "Testing" {
			t.Errorf("alerts customerTag, expected Testing, got %s", tag)
		}

		path := r.URL.Path
		if path != "/api/alert/all" {
			t.Errorf("alerts path, expected /api/alert/all, got %s", path)
		}
	}))
	defer ts.Close()

	conf := &Config{Address: strings.TrimLeft(ts.URL, "https://"),
		Token:         "123456789",
		SkipTLSVerify: true}

	if tc, err := NewClient(conf); err == nil {
		alerting := Alerting{client: tc}
		params := &QueryParams{"customerTag": "Testing"}
		if _, err := alerting.All(params); err != nil {
			t.Fatalf("creating Alerting failed with %s", err)
		} else {
			if string(alerting.RawResponse) != string(dummyAlert) {
				t.Errorf("get all alerts, incorrect response")
			}
		}

	}

}

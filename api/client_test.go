package wavefront

import (
	_ "fmt"
	_ "net/http/httputil"
	"testing"
)

var (
	testConfig = &Config{Address: "local.wavefront.com",
		Token: "123456789"}
	err error
	c   *Client
)

func TestNewClient(t *testing.T) {
	if c, err = NewClient(testConfig); err != nil {
		t.Fatal(err)
	}

	if c.BaseURL.String() != "https://local.wavefront.com" {
		t.Errorf("BaseURL expected https://local.wavefront.com, got %s", c.BaseURL)
	}
}

func TestNewRequest(t *testing.T) {
	p := &QueryParams{"test": "this"}

	if c, err = NewClient(testConfig); err != nil {
		t.Fatal(err)
	}

	if req, err := c.NewRequest("GET", "alerts", p); err != nil {
		t.Fatal(err)
	} else {
		// Test for Auth Token being passed correctly
		if h, ok := req.Header["X-Auth-Token"]; ok {
			if h[0] != testConfig.Token {
				t.Errorf("Auth token, expected %s got %s", testConfig.Token, h)
			}
		} else {
			t.Errorf("Missing X-AUTH-TOKEN header in request: %s", req.Header)
		}

		// Test that parameters are passed through
		if req.URL.RawQuery != "test=this" {
			t.Errorf("Query string expected test=this got %s", req.URL.RawQuery)
		}
	}
}

func TestNewRequestPost(t *testing.T) {
	if c, err = NewClient(testConfig); err != nil {
		t.Fatal(err)
	}

	if req, err := c.NewRequest("POST", "events", &QueryParams{}); err != nil {
		t.Fatal(err)
	} else {
		// Test that Content-Type header is set
		if h, ok := req.Header["Content-Type"]; ok {
			if h[0] != "application/x-www-form-urlencoded" {
				t.Errorf("Content-Type, expected application/x-www-form-urlencoded, got %s", h[0])
			}
		} else {
			t.Errorf("Missing x-www-form-urlencoded header in POST request.")
		}
	}

}

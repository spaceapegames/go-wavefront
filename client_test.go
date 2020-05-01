package wavefront

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientGet(t *testing.T) {
	params := &map[string]string{
		"s":                      "144242525262",
		"e":                      "142252272822",
		"includeObsoleteMetrics": "true",
	}

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		if r.URL.Path != "/api/v2/test/thing" {
			t.Errorf("request path, expected /api/v2/test/thing, got %s", r.URL.Path)
		}

		if header, ok := r.Header["Authorization"]; ok {
			if header[0] != "Bearer 123456789" {
				t.Errorf("Authorization header, expected 'Bearer 123456789', got %s", header[0])
			}
		} else {
			t.Errorf("no Authorization header set")
		}

		for k, v := range *params {
			if r.Form.Get(k) != v {
				t.Errorf("request param, expected %s, got %s", v, r.Form.Get(k))
			}
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer srv.Close()

	client, err := NewClient(&Config{
		Address:       strings.TrimLeft(srv.URL, "https://"),
		Token:         "123456789",
		SkipTLSVerify: true,
	})

	if err != nil {
		t.Fatal("error initiating client:", err)
	}

	req, err := client.NewRequest("GET", "test/thing", params, nil)
	if err != nil {
		t.Fatal("error creating request:", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error executing request:", err)
	}
	output, _ := ioutil.ReadAll(resp)
	fmt.Println(string(output))
}

func TestClientPost(t *testing.T) {
	params := &map[string]string{
		"s":                      "144242525262",
		"e":                      "142252272822",
		"includeObsoleteMetrics": "true",
	}
	body := []byte(`{ "some" : "json" }`)

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		if r.URL.Path != "/api/v2/test/thing" {
			t.Errorf("request path, expected /api/v2/test/thing, got %s", r.URL.Path)
		}

		if header, ok := r.Header["Authorization"]; ok {
			if header[0] != "Bearer 123456789" {
				t.Errorf("Authorization header, expected 'Bearer 123456789', got %s", header[0])
			}
		} else {
			t.Errorf("no Authorization header set")
		}

		if header, ok := r.Header["Content-Type"]; ok {
			if header[0] != "application/json" {
				t.Errorf("Authorization header, expected 'application/json', got %s", header[0])
			}
		} else {
			t.Errorf("no Content-Type header set")
		}

		actualBody, _ := ioutil.ReadAll(r.Body)
		if string(actualBody) != string(body) {
			t.Errorf("request body, expected %s got %s", string(body), string(actualBody))
		}

		for k, v := range *params {
			if r.Form.Get(k) != v {
				t.Errorf("request param, expected %s, got %s", v, r.Form.Get(k))
			}
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer srv.Close()

	client, err := NewClient(&Config{
		Address:       strings.TrimLeft(srv.URL, "https://"),
		Token:         "123456789",
		SkipTLSVerify: true,
	})

	if err != nil {
		t.Fatal("error initiating client:", err)
	}

	req, err := client.NewRequest("POST", "test/thing", params, body)
	if err != nil {
		t.Fatal("error creating request:", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error executing request:", err)
	}
	output, _ := ioutil.ReadAll(resp)
	fmt.Println(string(output))
}

func TestClientTLSConfig(t *testing.T) {
	params := &map[string]string{
		"s":                      "144242525262",
		"e":                      "142252272822",
		"includeObsoleteMetrics": "true",
	}

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		if r.URL.Path != "/api/v2/test/thing" {
			t.Errorf("request path, expected /api/v2/test/thing, got %s", r.URL.Path)
		}

		if header, ok := r.Header["Authorization"]; ok {
			if header[0] != "Bearer 123456789" {
				t.Errorf("Authorization header, expected 'Bearer 123456789', got %s", header[0])
			}
		} else {
			t.Errorf("no Authorization header set")
		}

		for k, v := range *params {
			if r.Form.Get(k) != v {
				t.Errorf("request param, expected %s, got %s", v, r.Form.Get(k))
			}
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer srv.Close()

	pool := x509.NewCertPool()
	pool.AddCert(srv.Certificate())
	tlsConfig := tls.Config{
		RootCAs: pool,
	}

	client, err := NewClient(&Config{
		Address:         strings.TrimLeft(srv.URL, "https://"),
		Token:           "123456789",
		TLSClientConfig: &tlsConfig,
	})

	if err != nil {
		t.Fatal("error initiating client:", err)
	}

	req, err := client.NewRequest("GET", "test/thing", params, nil)
	if err != nil {
		t.Fatal("error creating request:", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal("error executing request:", err)
	}
	output, _ := ioutil.ReadAll(resp)
	fmt.Println(string(output))
}

func TestClientTLSWithProxy(t *testing.T) {
	pool := x509.NewCertPool()
	tlsConfig := tls.Config{
		RootCAs: pool,
	}
	httpProxy := "http://example.com:8080"

	client, err := NewClient(&Config{
		Address:         "example.org",
		Token:           "123456789",
		TLSClientConfig: &tlsConfig,
		HttpProxy:       httpProxy,
	})

	if err != nil {
		t.Fatal("error initiating client:", err)
	}

	transport := client.httpClient.Transport.(*http.Transport)
	if transport.TLSClientConfig != &tlsConfig {
		t.Fatal("TLS Client Config not preserved")
	}

	transportProxy, _ := transport.Proxy(nil)
	if httpProxy != transportProxy.String() {
		t.Fatal("HttpProxy not preserved")
	}
}

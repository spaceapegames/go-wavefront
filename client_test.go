package wavefront

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	asserts "github.com/stretchr/testify/assert"
)

func TestClientGet(t *testing.T) {
	params := &map[string]string{
		"s":                      "144242525262",
		"e":                      "142252272822",
		"includeObsoleteMetrics": "true",
	}

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Error parsing form: %v", err)
		}

		if r.URL.Path != "/api/v2/test/thing" {
			t.Errorf("request path, expected /api/v2/test/thing, got %s", r.URL.Path)
		}

		if header, ok := r.Header["Authorization"]; ok {
			if header[0] != "Bearer 123456789" {
				t.Errorf("authorization header, expected 'Bearer 123456789', got %s", header[0])
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
		Address:       strings.TrimPrefix(srv.URL, "https://"),
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
		if err := r.ParseForm(); err != nil {
			t.Fatalf("Error parsing form: %v", err)
		}

		if r.URL.Path != "/api/v2/test/thing" {
			t.Errorf("request path, expected /api/v2/test/thing, got %s", r.URL.Path)
		}

		if header, ok := r.Header["Authorization"]; ok {
			if header[0] != "Bearer 123456789" {
				t.Errorf("authorization header, expected 'Bearer 123456789', got %s", header[0])
			}
		} else {
			t.Errorf("no Authorization header set")
		}

		if header, ok := r.Header["Content-Type"]; ok {
			if header[0] != "application/json" {
				t.Errorf("authorization header, expected 'application/json', got %s", header[0])
			}
		} else {
			t.Errorf("no Content-Type header set")
		}

		actualBody, _ := ioutil.ReadAll(r.Body)
		// The request body is buffered since we need to replay it on failure
		// this means the first read will fire this function above with an empty body (because we read it)
		if string(actualBody) != "" {
			if string(actualBody) != string(body) {
				t.Errorf("request body, expected %s got %s", string(body), string(actualBody))
			}
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
		Address:       strings.TrimPrefix(srv.URL, "https://"),
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

type testPointType struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type fakeWavefronter struct {
	newRequestError error
	method          string
	path            string
	params          map[string]string
	body            string
	doError         error
	response        string
}

func (f *fakeWavefronter) NewRequest(method, path string, params *map[string]string, body []byte) (*http.Request, error) {
	f.method = method
	f.path = path
	if params != nil {
		f.params = make(map[string]string, len(*params))
		for k, v := range *params {
			f.params[k] = v
		}
	} else {
		f.params = nil
	}
	if body != nil {
		f.body = string(body)
	} else {
		f.body = ""
	}
	return nil, f.newRequestError
}

func (f *fakeWavefronter) Do(_ *http.Request) (io.ReadCloser, error) {
	if f.doError != nil {
		return nil, f.doError
	}
	result := ioutil.NopCloser(strings.NewReader(f.response))
	return result, nil
}

var (
	doRestError = errors.New("a REST error")
)

func TestDoRest_NewRequestError(t *testing.T) {
	assert := asserts.New(t)
	fake := &fakeWavefronter{newRequestError: doRestError}
	err := doRest(
		"DELETE",
		"/a/rest/path",
		fake)
	assert.Equal(doRestError, err)
}

func TestDoRest_DoError(t *testing.T) {
	assert := asserts.New(t)
	fake := &fakeWavefronter{doError: doRestError}
	err := doRest(
		"DELETE",
		"/a/rest/path",
		fake)
	assert.Equal(doRestError, err)
}

func TestDoRest_NoOptions(t *testing.T) {
	assert := asserts.New(t)
	fake := &fakeWavefronter{response: "Some response."}
	err := doRest(
		"DELETE",
		"/a/rest/path",
		fake)
	assert.NoError(err)
	assert.Equal("DELETE", fake.method)
	assert.Equal("/a/rest/path", fake.path)
	assert.Nil(fake.params)
	assert.Empty(fake.body)
}

func TestDoRest_UnexpectedResponse(t *testing.T) {
	assert := asserts.New(t)
	fake := &fakeWavefronter{response: "A bad response."}
	var result testPointType
	err := doRest(
		"POST",
		"/a/rest/path",
		fake,
		doOutput(&result))
	assert.Error(err)
}

func TestDoRest(t *testing.T) {
	responseStr := `
{
  "status": {
      "code": 200,
      "message": "OK"
  },
  "response": {
    "x": 42,
    "y": 63
  }
}`
	bodyStr := `{"x":3,"y":5}`
	assert := asserts.New(t)
	fake := &fakeWavefronter{response: responseStr}
	var result testPointType
	params := map[string]string{"email": "true"}
	paramOption := doParams(params)

	// Test that mutating the map doesn't affect our option
	params["extra"] = "big"
	err := doRest(
		"POST",
		"/a/rest/path",
		fake,
		doInput(&testPointType{X: 3, Y: 5}),
		doOutput(&result),
		paramOption)
	assert.NoError(err)
	assert.Equal("POST", fake.method)
	assert.Equal("/a/rest/path", fake.path)
	assert.Equal(map[string]string{"email": "true"}, fake.params)
	assert.Equal(bodyStr, fake.body)
	assert.Equal(testPointType{X: 42, Y: 63}, result)
}

func TestConfigDefensiveCopy(t *testing.T) {
	assert := asserts.New(t)
	config := &Config{
		Address:       "somehost.wavefront.com",
		Token:         "123456789",
		SkipTLSVerify: true,
	}
	client, _ := NewClient(config)
	assert.NotSame(config, client.Config)
	assert.Equal("somehost.wavefront.com", client.Config.Address)
	assert.Equal("123456789", client.Config.Token)
}

func TestDoRest_DirectResponse(t *testing.T) {
	responseStr := `
{
  "x": 42,
  "y": 63
}`
	assert := asserts.New(t)
	fake := &fakeWavefronter{response: responseStr}
	var result testPointType
	err := doRest(
		"GET",
		"/a/rest/path",
		fake,
		doOutput(&result),
		doDirectResponse())
	assert.NoError(err)
	assert.Equal("GET", fake.method)
	assert.Equal("/a/rest/path", fake.path)
	assert.Equal(testPointType{X: 42, Y: 63}, result)
}

type testPrimeStructType struct {
	Primes []int `json:"primes"`
}

func TestDoRest_SafeModify(t *testing.T) {
	responseStr := `
{
    "primes": [2, 3, 5]
}`
	assert := asserts.New(t)
	fake := &fakeWavefronter{response: responseStr}
	result := testPrimeStructType{Primes: []int{0, 1, 2, 3, 4}}
	original := result
	assert.NoError(doRest(
		"GET",
		"/a/rest/path",
		fake,
		doOutput(&result),
		doDirectResponse()))
	assert.Equal([]int{2, 3, 5}, result.Primes)

	// doRest should modify result in such a way that the slice from the
	// shallow copy remains intact
	assert.Equal([]int{0, 1, 2, 3, 4}, original.Primes)
}

func TestDoRest_SafeModifyResponse(t *testing.T) {
	responseStr := `
{
  "status": {
      "code": 200,
      "message": "OK"
  },
  "response": {
    "primes": [2, 3, 5]
  }
}`
	assert := asserts.New(t)
	fake := &fakeWavefronter{response: responseStr}
	result := testPrimeStructType{Primes: []int{0, 1, 2, 3, 4}}
	original := result
	assert.NoError(doRest(
		"GET",
		"/a/rest/path",
		fake,
		doOutput(&result)))
	assert.Equal([]int{2, 3, 5}, result.Primes)

	// doRest should modify result in such a way that the slice from the
	// shallow copy remains intact
	assert.Equal([]int{0, 1, 2, 3, 4}, original.Primes)
}

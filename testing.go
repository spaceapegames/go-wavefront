package wavefront

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

// Helps expedite the boiler plate code for testing client requests
// iface is the destination type we need marshal data into from the request
// iface can be utilized to ensure the request body is properly being marshalled and values are set as expected
func testDo(t *testing.T, req *http.Request, fixture, method string, iface interface{}) (io.ReadCloser, error) {
	response, err := ioutil.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t, method, req.Method)
	if req.Body != nil {
		body, _ := ioutil.ReadAll(req.Body)
		err = json.Unmarshal(body, iface)
		if err != nil {
			t.Fatal(err)
		}
	} else {
		t.Log("Request body was nil, if this is expected please ignore...")
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
}

// Helps expedite the boiler plate code for testing client requests against paginated results
// tMarshal is the destination type we need marshal data into from the response body of the fixture
func testPaginatedDo(t *testing.T, req *http.Request, fixture string, invokedCount *int) (io.ReadCloser, error) {
	search := SearchParams{}
	resp, err := testDo(t, req, fmt.Sprintf(fixture, *invokedCount), "POST", &search)
	assertEqual(t, search.Limit*(*invokedCount), search.Offset)
	*invokedCount++
	return resp, err
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected %v (type %v), actual %v (type %v)", expected, reflect.TypeOf(expected),
			actual, reflect.TypeOf(actual))
	}
}

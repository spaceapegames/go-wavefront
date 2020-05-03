package wavefront

import (
	"encoding/json"
	"io/ioutil"
)

// Performs a basic basicCrud action against a particular endpoint and
// automatically marshals the data back to the appropriate concrete type
// Expects the response to always be a json object containing at least a `response` struct
func basicCrud(client Wavefronter, method, path string, t interface{}, params *map[string]string) error {
	payload, err := json.Marshal(t)
	if err != nil {
		return err
	}

	req, err := client.NewRequest(method, path, params, payload)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Close()
	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &struct {
		Response interface{} `json:"response"`
	}{
		Response: t,
	})
}

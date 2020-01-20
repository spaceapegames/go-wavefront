package wavefront

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type AccessControlList struct {
	CanView   []string `json:"canView,omitempty"`
	CanModify []string `json:"canModify,omitempty"`
}

func putEntityACL(id string, canView []string, canModify []string, basePath string, client Wavefronter) error {
	if id == "" {
		return fmt.Errorf("id must not be empty")
	}
	payload, err := json.Marshal(&[]struct {
		EntityID  string   `json:"entityId"`
		ViewACL   []string `json:"viewAcl,omitempty"`
		ModifyACL []string `json:"modifyAcl,omitempty"`
	}{
		{
			EntityID:  id,
			ViewACL:   canView,
			ModifyACL: canModify,
		},
	})

	if err != nil {
		return err
	}

	req, err := client.NewRequest("PUT", fmt.Sprintf("%s/acl/set", basePath), nil, payload)
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

	if len(body) > 0 {
		return fmt.Errorf("expected no response, got %s", string(body))
	}

	return nil
}

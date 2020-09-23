package wavefront

import (
	"fmt"
)

type AccessControlList struct {
	CanView   []string `json:"canView,omitempty"`
	CanModify []string `json:"canModify,omitempty"`
}

func putEntityACL(id string, canView []string, canModify []string, basePath string, client Wavefronter) error {
	if id == "" {
		return fmt.Errorf("id must not be empty")
	}
	acls := []struct {
		EntityID  string   `json:"entityId"`
		ViewACL   []string `json:"viewAcl,omitempty"`
		ModifyACL []string `json:"modifyAcl,omitempty"`
	}{
		{
			EntityID:  id,
			ViewACL:   canView,
			ModifyACL: canModify,
		},
	}
	return doRest(
		"PUT",
		fmt.Sprintf("%s/acl/set", basePath),
		client,
		doPayload(acls))
}

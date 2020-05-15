package wavefront

import (
	"encoding/json"
	"fmt"
)

type Role struct {
	ID                   string       `json:"id"`
	CreatedEpochMillis   int          `json:"createdEpochMillis,omitempty"`
	LastUpdatedMs        int          `json:"lastUpdatedMs,omitempty"`
	SampleLinkedGroups   *[]UserGroup `json:"sampleLinkedGroups,omitempty"`
	SampleLinkedAccounts *[]string    `json:"sampleLinkedAccounts,omitempty"`
	LinkedGroupsCount    int          `json:"linkedGroupsCount,omitempty"`
	LinkedAccountsCount  int          `json:"linkedAccountsCount,omitempty"`
	Customer             string       `json:"customer,omitempty"`
	LastUpdatedAccountId string       `json:"lastUpdatedAccountId,omitempty"`
	Name                 string       `json:"name"`
	Permissions          []string     `json:"permissions,omitempty"`
	Description          string       `json:"description,omitempty"`
}

const baseRoleUrlPath = "/api/v2/role"

type Roles struct {
	client Wavefronter
}

func (c *Client) Roles() *Roles {
	return &Roles{client: c}
}

func (r Roles) Find(filter []*SearchCondition) ([]*Role, error) {
	search := &Search{
		client: r.client,
		Type:   "role",
		Params: &SearchParams{
			Conditions: filter,
		},
	}

	var results []*Role
	moreItems := true
	for moreItems == true {
		resp, err := search.Execute()
		if err != nil {
			return nil, err
		}
		var tmpres []*Role
		err = json.Unmarshal(resp.Response.Items, &tmpres)
		if err != nil {
			return nil, err
		}
		results = append(results, tmpres...)
		moreItems = resp.Response.MoreItems
		search.Params.Offset = resp.NextOffset
	}

	return results, nil
}

func (r Roles) Create(role *Role) error {
	if role.Name == "" {
		return fmt.Errorf("name must be specified while creating a role")
	}

	return basicCrud(r.client, "POST", baseRoleUrlPath, role, nil)
}

// Get a specific Role for the given ID
// ID field must be specified
func (r Roles) Get(role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}

	return basicCrud(r.client, "GET", fmt.Sprintf("%s/%s", baseRoleUrlPath, role.ID), role, nil)
}

// Update a specific Role for the given ID
// ID field must be specified
func (r Roles) Update(role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}

	return basicCrud(r.client, "PUT", fmt.Sprintf("%s/%s", baseRoleUrlPath, role.ID), role, nil)
}

func (r Roles) Delete(role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}

	return basicCrud(r.client, "DELETE", fmt.Sprintf("%s/%s", baseRoleUrlPath, role.ID), role, nil)
}

func (r Roles) AddAssignees(assignees []string, role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}
	return crudWithPayload(r.client, "POST",
		fmt.Sprintf("%s/%s/%s", baseRoleUrlPath, role.ID, "addAssignees"), assignees, role, nil)
}

func (r Roles) RemoveAssignees(assignees []string, role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}
	return crudWithPayload(r.client, "POST",
		fmt.Sprintf("%s/%s/%s", baseRoleUrlPath, role.ID, "removeAssignees"), assignees, role, nil)
}

func (r Roles) GrantPermission(permission string, roles []*Role) error {
	var roleIds []string
	if len(roles) == 0 {
		return fmt.Errorf("must specify at least one role to modify")
	}

	for _, role := range roles {
		if role.ID == "" {
			return fmt.Errorf("the ID field must be specified")
		}
		roleIds = append(roleIds, role.ID)
	}

	return crudWithPayload(r.client, "POST",
		fmt.Sprintf("%s/%s/%s", baseRoleUrlPath, "grant", permission), roleIds, roles, nil)

}

func (r Roles) RevokePermission(permission string, roles []*Role) error {
	var roleIds []string
	if len(roles) == 0 {
		return fmt.Errorf("must specify at least one role to modify")
	}

	for _, role := range roles {
		if role.ID == "" {
			return fmt.Errorf("the ID field must be specified")
		}
		roleIds = append(roleIds, role.ID)
	}

	return crudWithPayload(r.client, "POST",
		fmt.Sprintf("%s/%s/%s", baseRoleUrlPath, "revoke", permission), roleIds, roles, nil)
}

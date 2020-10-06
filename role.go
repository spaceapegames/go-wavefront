package wavefront

import (
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

func (r Roles) Find(filter []*SearchCondition) (roles []*Role, err error) {
	err = doSearch(filter, "role", r.client, &roles)
	return
}

func (r Roles) Create(role *Role) error {
	if role.Name == "" {
		return fmt.Errorf("name must be specified while creating a role")
	}

	return doRest(
		"POST",
		baseRoleUrlPath,
		r.client,
		doPayload(role),
		doResponse(role))
}

// Get a specific Role for the given ID
// ID field must be specified
func (r Roles) Get(role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}

	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", baseRoleUrlPath, role.ID),
		r.client,
		doResponse(role))
}

// Update a specific Role for the given ID
// ID field must be specified
func (r Roles) Update(role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}

	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", baseRoleUrlPath, role.ID),
		r.client,
		doPayload(role),
		doResponse(role))
}

func (r Roles) Delete(role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}

	return doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", baseRoleUrlPath, role.ID),
		r.client)
}

func (r Roles) AddAssignees(assignees []string, role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}
	return doRest(
		"POST",
		fmt.Sprintf("%s/%s/addAssignees", baseRoleUrlPath, role.ID),
		r.client,
		doPayload(assignees),
		doResponse(role))
}

func (r Roles) RemoveAssignees(assignees []string, role *Role) error {
	if role.ID == "" {
		return fmt.Errorf("the ID field must be specified")
	}
	return doRest(
		"POST",
		fmt.Sprintf("%s/%s/removeAssignees", baseRoleUrlPath, role.ID),
		r.client,
		doPayload(assignees),
		doResponse(role))
}

func (r Roles) GrantPermission(permission string, roles []*Role) error {
	if len(roles) == 0 {
		return fmt.Errorf("must specify at least one role to modify")
	}

	roleIds := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.ID == "" {
			return fmt.Errorf("the ID field must be specified")
		}
		roleIds = append(roleIds, role.ID)
	}
	return doRest(
		"POST",
		fmt.Sprintf("%s/grant/%s", baseRoleUrlPath, permission),
		r.client,
		doPayload(roleIds))
}

func (r Roles) RevokePermission(permission string, roles []*Role) error {
	if len(roles) == 0 {
		return fmt.Errorf("must specify at least one role to modify")
	}

	roleIds := make([]string, 0, len(roles))
	for _, role := range roles {
		if role.ID == "" {
			return fmt.Errorf("the ID field must be specified")
		}
		roleIds = append(roleIds, role.ID)
	}
	return doRest(
		"POST",
		fmt.Sprintf("%s/%s/%s", baseRoleUrlPath, "revoke", permission),
		r.client,
		doPayload(roleIds))
}

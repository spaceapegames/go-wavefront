package wavefront

import (
	"fmt"
)

type UserGroup struct {
	// Unique ID for the user group
	ID *string `json:"id,omitempty"`

	// Name of the user group
	Name string `json:"name,omitempty"`

	// Customer
	Customer string `json:"customer,omitempty"`

	// Roles assigned to the group
	Roles []Role

	// Users that are members of the group
	Users []string `json:"users,omitempty"`

	// Total number of users that are members of the user group
	UserCount int `json:"userCount,omitempty"`

	// Which properties of the user group are editable
	Properties UserGroupPropertiesDTO `json:"properties,omitempty"`

	// Description of the Group purpose
	Description string `json:"description,omitempty"`

	// When the group was created
	CreatedEpochMillis int `json:"createdEpochMillis,omitempty"`
}

type UserGroupPropertiesDTO struct {
	NameEditable bool `json:"nameEditable"`

	PermissionsEditable bool `json:"permissionsEditable"`

	UsersEditable bool `json:"usersEditable"`

	RolesEditable bool `json:"rolesEditable"`
}

const baseUserGroupPath = "/api/v2/usergroup"

type UserGroups struct {
	client Wavefronter
}

// UserGroups is used to return a client for user-group related operations
func (c *Client) UserGroups() *UserGroups {
	return &UserGroups{client: c}
}

func (g UserGroups) Create(userGroup *UserGroup) error {
	if userGroup.Name == "" {
		return fmt.Errorf("name must be specified when creating a usergroup")
	}

	return doRest(
		"POST",
		baseUserGroupPath,
		g.client,
		doPayload(userGroup),
		doResponse(userGroup))
}

// Gets a specific UserGroup by ID
// The ID field must be specified
func (g UserGroups) Get(userGroup *UserGroup) error {
	if userGroup.ID == nil || *userGroup.ID == "" {
		return fmt.Errorf("usergroup ID field is not set")
	}

	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", baseUserGroupPath, *userGroup.ID),
		g.client,
		doResponse(userGroup))
}

// Find returns all UsersGroups filtered by the given search conditions.
// If filter is nil, all UserGroups are returned.
func (g UserGroups) Find(filter []*SearchCondition) (
	results []*UserGroup, err error) {
	err = doSearch(filter, "usergroup", g.client, &results)
	return
}

// Update does not support updating the users on the group
// To update the users in a group use AddUsers and RemoveUsers
// The ID field must be specified
func (g UserGroups) Update(userGroup *UserGroup) error {
	if userGroup.ID == nil || *userGroup.ID == "" {
		return fmt.Errorf("usergroup ID must be specified")
	}

	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", baseUserGroupPath, *userGroup.ID),
		g.client,
		doPayload(userGroup),
		doResponse(userGroup))
}

// Adds the specified users to the group
func (g UserGroups) AddUsers(id *string, users *[]string) error {
	if *id == "" {
		return fmt.Errorf("usergroup ID must be specified")
	}

	return g.updateUserGroupUsers(users, id, "addUsers")
}

// Removes the specified users from the group
func (g UserGroups) RemoveUsers(id *string, users *[]string) error {
	if *id == "" {
		return fmt.Errorf("usergroup ID must be specified")
	}

	return g.updateUserGroupUsers(users, id, "removeUsers")
}

func (g UserGroups) Delete(userGroup *UserGroup) error {
	if *userGroup.ID == "" {
		return fmt.Errorf("usergroup ID must be specified")
	}

	err := doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", baseUserGroupPath, *userGroup.ID),
		g.client)
	if err != nil {
		return err
	}

	// Clear the ID so that we do not accidentally re-submit
	empty := ""
	userGroup.ID = &empty
	return nil
}

func (g UserGroups) updateUserGroupUsers(users *[]string, id *string, endpoint string) error {
	return doRest(
		"POST",
		fmt.Sprintf("%s/%s/%s", baseUserGroupPath, *id, endpoint),
		g.client,
		doPayload(users))
}

package wavefront

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type UserGroup struct {
	// Unique ID for the user group
	ID *string `json:"id,omitempty"`

	// Name of the user group
	Name string `json:"name,omitempty"`

	// Permission(s) assigned to the user group
	Permissions []string `json:"permissions,omitempty"`

	// Customer
	Customer string `json:"customer,omitempty"`

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
	if len(userGroup.Permissions) == 0 {
		return fmt.Errorf("permissions must be specified when creating a usergroup")
	}

	return g.crudUserGroup("POST", baseUserGroupPath, userGroup)
}

// Gets a specific UserGroup by ID
// The ID field must be specified
func (g UserGroups) Get(userGroup *UserGroup) error {
	if *userGroup.ID == "" {
		return fmt.Errorf("usergroup ID field is not set")
	}

	return g.crudUserGroup("GET", fmt.Sprintf("%s/%s", baseUserGroupPath, *userGroup.ID), userGroup)
}

// Find returns all UsersGroups filtered by the given search conditions.
// If filter is nil, all UserGroups are returned.
func (g UserGroups) Find(filter []*SearchCondition) ([]*UserGroup, error) {
	search := &Search{
		client: g.client,
		Type:   "usergroup",
		Params: &SearchParams{
			Conditions: filter,
		},
	}

	var results []*UserGroup
	moreItems := true
	for moreItems == true {
		resp, err := search.Execute()
		if err != nil {
			return nil, err
		}
		var tmpres []*UserGroup
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

// Update does not support updating the users on the group
// To update the users in a group use AddUsers and RemoveUsers
// The ID field must be specified
func (g UserGroups) Update(userGroup *UserGroup) error {
	if *userGroup.ID == "" {
		return fmt.Errorf("usergroup ID must be specified")
	}

	return g.crudUserGroup("PUT", fmt.Sprintf("%s/%s", baseUserGroupPath, *userGroup.ID), userGroup)
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

	err := g.crudUserGroup("DELETE", fmt.Sprintf("%s/%s", baseUserGroupPath, *userGroup.ID), userGroup)
	if err != nil {
		return err
	}

	// Clear the ID so that we do not accidentally re-submit
	*userGroup.ID = ""
	return nil
}

func (g UserGroups) updateUserGroupUsers(users *[]string, id *string, endpoint string) error {
	payload, err := json.Marshal(users)
	if err != nil {
		return err
	}

	req, err := g.client.NewRequest("POST", fmt.Sprintf("%s/%s/%s", baseUserGroupPath, *id, endpoint), nil, payload)
	if err != nil {
		return err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	_, err = ioutil.ReadAll(resp)
	return err
}

func (g UserGroups) crudUserGroup(method, path string, userGroup *UserGroup) error {
	payload, err := json.Marshal(userGroup)
	if err != nil {
		return err
	}
	req, err := g.client.NewRequest(method, path, nil, payload)
	if err != nil {
		return err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &struct {
		Response *UserGroup `json:"response"`
	}{
		Response: userGroup,
	})
}

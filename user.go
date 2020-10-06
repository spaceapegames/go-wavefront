package wavefront

import (
	"encoding/json"
	"fmt"
)

type NewUserRequest struct {
	// The only time it is referred to as emailAddress is when it's a new user
	EmailAddress string `json:"emailAddress"`

	// The permissions granted to this user
	Permissions []string `json:"groups,omitempty"`

	// Groups this user belongs to
	// This is wrapped with a Wrapper to manage the serialization between what we send to the API
	// And what the API sends back (which is to say, we send just IDs but we always receive a complete object)
	Groups UserGroupsWrapper `json:"userGroups,omitempty"`
}

// User represents a Wavefront User
type User struct {
	// The email identifier for a user
	ID *string `json:"identifier"`

	// The customer the user is a member of
	Customer string `json:"customer,omitempty"`

	// Last successful login in epoch millis
	LastSuccessfulLogin uint `json:"lastSuccessfulLogin,omitempty"`

	// The permissions granted to this user
	Permissions []string `json:"groups,omitempty"`

	// Groups this user belongs to
	// This is wrapped with a Wrapper to manage the serialization between what we send to the API
	// And what the API sends back (which is to say, we send just IDs but we always receive a complete object)
	Groups UserGroupsWrapper `json:"userGroups,omitempty"`

	// Used during an PUT call to modify a users password
	Credential string `json:"credential,omitempty"`
}

type UserGroupsWrapper struct {
	UserGroups []UserGroup
}

// Users is used to perform user-related operations against the Wavefront API
type Users struct {
	// client is the Wavefront client used to perform target-related operations
	client Wavefronter
}

const (
	baseUserPath               = "/api/v2/user"
	AGENT_MANAGEMENT           = "agent_management"
	ALERTS_MANAGEMENT          = "alerts_management"
	DASHBOARD_MANAGEMENT       = "dashboard_management"
	EMBEDDED_CHARTS_MANAGEMENT = "embedded_charts"
	EVENTS_MANAGEMENT          = "events_management"
	EXTERNAL_LINKS_MANAGEMENT  = "external_links_management"
	HOST_TAG_MANAGEMENT        = "host_tag_management"
	METRICS_MANAGEMENT         = "metrics_management"
	USER_MANAGEMENT            = "user_management"
	INTEGRATIONS_MANAGEMENT    = "application_management"
	DIRECT_INGESTION           = "ingestion"
	BATCH_QUERY_PRIORITY       = "batch_query_priority"
	DERIVED_METRICS_MANAGEMENT = "derived_metrics_management"
)

// Users is used to return a client for user-related operations
func (c *Client) Users() *Users {
	return &Users{client: c}
}

// Get is used to retrieve an existing User by ID.
// The identifier field must be specified
func (u Users) Get(user *User) error {
	if user.ID == nil || *user.ID == "" {
		return fmt.Errorf("user ID field is not set")
	}

	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", baseUserPath, *user.ID),
		u.client,
		doDirectResponse(user))
}

// Find returns all Users filtered by the given search conditions.
// If filter is nil, all Users are returned.
// UserGroups returned on the User from this call will be ID only
func (u Users) Find(filter []*SearchCondition) (users []*User, err error) {
	err = doSearch(filter, "user", u.client, &users)
	return
}

// Does not support specifying a credential
// The EmailAddress field must be specified
func (u Users) Create(newUser *NewUserRequest, user *User, sendEmail bool) error {
	if newUser.EmailAddress == "" {
		return fmt.Errorf("a valid email address must be specified")
	}

	params := map[string]string{
		"sendEmail": fmt.Sprintf("%t", sendEmail),
	}
	return doRest(
		"POST",
		baseUserPath,
		u.client,
		doParams(params),
		doPayload(newUser),
		doResponse(user))
}

// Supports specifying the credential
// The identifier field must be specified
func (u Users) Update(user *User) error {
	if user.ID == nil || *user.ID == "" {
		return fmt.Errorf("user ID field is not set")
	}
	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", baseUserPath, *user.ID),
		u.client,
		doPayload(user),
		doResponse(user))
}

// Deletes the specified user
// The ID field must be specified
func (u Users) Delete(user *User) error {
	if user.ID == nil || *user.ID == "" {
		return fmt.Errorf("user ID field is not set")
	}

	err := doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", baseUserPath, *user.ID),
		u.client)
	if err != nil {
		return err
	}
	empty := ""
	user.ID = &empty
	return nil
}

// During a GET operation or returned AFTER a POST/PUT/DELETE we receive a complete UserGroup struct
func (w *UserGroupsWrapper) UnmarshalJSON(data []byte) error {
	// First try to unmarshal it as an array of string IDs
	// The Search API returns only group ids when looking up Users
	var groupIds []*string
	if err := json.Unmarshal(data, &groupIds); err == nil {
		// We need to go ahead and bind these IDs to empty groups on the UserGroupsWrapper
		for _, v := range groupIds {
			w.UserGroups = append(w.UserGroups, UserGroup{
				ID: v,
			})
		}
		return nil
	}

	// Failing that lets try to just unmarshal the groups directly
	return json.Unmarshal(data, &w.UserGroups)
}

// During a POST/PUT/DELETE on a User only the UserGroup.ID is transmitted
func (w *UserGroupsWrapper) MarshalJSON() ([]byte, error) {
	var ids []*string
	if w.UserGroups != nil {
		for _, v := range w.UserGroups {
			if *v.ID != "" {
				ids = append(ids, v.ID)
			}
		}
	}
	return json.Marshal(ids)
}

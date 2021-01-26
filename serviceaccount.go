package wavefront

import (
	"fmt"
)

const (
	saEndpoint = "/api/v2/account/serviceaccount"
)

// ServiceAccount represents a ServiceAccount which exists in Wavefront. Note that here, roles,
// groups, tokens and the ingestion policy are embedded structs.
type ServiceAccount struct {
	ID              string          `json:"identifier"`
	Description     string          `json:"description"`
	Permissions     []string        `json:"groups"`
	Active          bool            `json:"active"`
	Roles           []Role          `json:"roles"`
	UserGroups      []UserGroup     `json:"userGroups"`
	Tokens          []Token         `json:"tokens"`
	IngestionPolicy IngestionPolicy `json:"ingestionPolicy"`
}

// TokenIds returns the Ids of the tokens in this instance.
func (s *ServiceAccount) TokenIds() []string {
	result := make([]string, 0, len(s.Tokens))
	for i := range s.Tokens {
		result = append(result, s.Tokens[i].ID)
	}
	return result
}

// RoleIds returns the Ids of the Roles in this instance
func (s *ServiceAccount) RoleIds() []string {
	result := make([]string, 0, len(s.Roles))
	for i := range s.Roles {
		result = append(result, s.Roles[i].ID)
	}
	return result
}

// UserGroupIds returns the Ids of the UserGroups in this instance
func (s *ServiceAccount) UserGroupIds() []string {
	result := make([]string, 0, len(s.UserGroups))
	for i := range s.UserGroups {
		if s.UserGroups[i].ID != nil {
			result = append(result, *s.UserGroups[i].ID)
		}
	}
	return result
}

// IngestionPolicyId returns the Id of the ingestion policy in this instance
func (s *ServiceAccount) IngestionPolicyId() string {
	return s.IngestionPolicy.ID
}

// Options returns a ServiceAccountOptions prepopulated with the settings
// from this instance. Use this to update a ServiceAccount.
func (s *ServiceAccount) Options() *ServiceAccountOptions {
	return &ServiceAccountOptions{
		ID:                s.ID,
		Active:            s.Active,
		Description:       s.Description,
		Permissions:       s.Permissions,
		Roles:             s.RoleIds(),
		UserGroups:        s.UserGroupIds(),
		IngestionPolicyID: s.IngestionPolicyId(),
	}
}

// ServiceAccountOptions represents the options for creating or updating
// a ServiceAccount in Wavefront. Note that here, Roles, UserGroups and the ingestion policy are
// the IDs of existing objects.
type ServiceAccountOptions struct {

	// Required
	ID string `json:"identifier"`

	// Required
	Active bool `json:"active"`

	// Always leave empty for now.
	Tokens []string `json:"tokens"`

	Description       string   `json:"description,omitempty"`
	Permissions       []string `json:"groups,omitempty"`
	Roles             []string `json:"roles,omitempty"`
	UserGroups        []string `json:"userGroups,omitempty"`
	IngestionPolicyID string   `json:"ingestionPolicyId,omitempty"`
}

// ServiceAccounts is used to perform service account related operations
// against the Wavefront API
type ServiceAccounts struct {
	// client is the Wavefront client used to perform target-related operations
	client Wavefronter
}

// ServiceAccounts is used to return a client for service account related
// operations
func (c *Client) ServiceAccounts() *ServiceAccounts {
	return &ServiceAccounts{client: c}
}

// Find returns all service accounts filtered by the given search conditions.
// If filter is nil, all service accounts are returned.
func (s *ServiceAccounts) Find(filter []*SearchCondition) (
	results []*ServiceAccount, err error) {
	err = doSearch(filter, "serviceaccount", s.client, &results)
	return
}

// GetByID returns the ServiceAccount with given ID. If no such
// ServiceAccount exists, GetByID returns an error. The caller can call
// NotFound on err to determine whether or not the error is because the
// ServiceAccount doesn't exist.
func (s *ServiceAccounts) GetByID(id string) (
	serviceAccount *ServiceAccount, err error) {
	var result ServiceAccount
	err = doRest(
		"GET",
		fmt.Sprintf("%s/%s", saEndpoint, id),
		s.client,
		doResponse(&result))
	if err != nil {
		return
	}
	return &result, nil
}

// Create creates a ServiceAccount according to options and returns the
// newly created ServiceAccount.
func (s *ServiceAccounts) Create(options *ServiceAccountOptions) (
	serviceAccount *ServiceAccount, err error) {
	var result ServiceAccount
	err = doRest(
		"POST",
		saEndpoint,
		s.client,
		doPayload(fixServiceAccountOptions(*options)),
		doResponse(&result))
	if err != nil {
		return
	}
	return &result, nil
}

// Update updates a ServiceAccount according to options and returns the
// updated ServiceAccount.
func (s *ServiceAccounts) Update(options *ServiceAccountOptions) (
	serviceAccount *ServiceAccount, err error) {
	var result ServiceAccount
	err = doRest(
		"PUT",
		fmt.Sprintf("%s/%s", saEndpoint, options.ID),
		s.client,
		doPayload(fixServiceAccountOptions(*options)),
		doResponse(&result))
	if err != nil {
		return
	}
	return &result, nil
}

// DeleteById deletes the ServiceAccount with given id.
func (s *ServiceAccounts) DeleteByID(id string) error {
	return doRest(
		"DELETE",
		fmt.Sprintf("/api/v2/account/%s", id),
		s.client)
}

// This is a work around for a bug in the wavefront Rest API. It
// returns a value like options except with the token field set to an
// empty slice.
func fixServiceAccountOptions(
	options ServiceAccountOptions) *ServiceAccountOptions {
	options.Tokens = []string{}
	return &options
}

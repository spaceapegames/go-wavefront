package wavefront

import (
	"fmt"
)

const (
	tokenEndpoint = "/api/v2/apitoken/serviceaccount"
)

// Token represents an API token in Wavefront
type Token struct {
	ID   string `json:"tokenID"`
	Name string `json:"tokenName"`
}

// Options returns the options for this token.
func (s *Token) Options() *TokenOptions {
	return &TokenOptions{
		ID:   s.ID,
		Name: s.Name,
	}
}

// TokenOptions represents the options for creating or modifying
// tokens
type TokenOptions struct {

	// ID must be blank when creating a token.
	ID string `json:"tokenID,omitempty"`

	Name string `json:"tokenName,omitempty"`
}

// Tokens is used to perform service account token related operations
// against the Wavefront API
type Tokens struct {
	// client is the Wavefront client used to perform target-related operations
	client Wavefronter
}

// Tokens is used to return a client for service account token related
// operations
func (c *Client) Tokens() *Tokens {
	return &Tokens{client: c}
}

// Create creates a Token according to options and returns all the
// tokens including the newly created one for the given serviceAccountID.
func (s *Tokens) Create(
	serviceAccountID string, options *TokenOptions) (
	tokens []Token, err error) {
	var result []Token
	err = doRest(
		"POST",
		fmt.Sprintf("%s/%s", tokenEndpoint, serviceAccountID),
		s.client,
		doPayload(options),
		doResponse(&result))
	if err != nil {
		return
	}
	return result, nil
}

// Delete deletes a token
func (s *Tokens) Delete(serviceAccountID, tokenID string) error {
	return doRest(
		"DELETE",
		fmt.Sprintf("%s/%s/%s", tokenEndpoint, serviceAccountID, tokenID),
		s.client)
}

// Update updates a token and returns the updated token.
func (s *Tokens) Update(
	serviceAccountID string, options *TokenOptions) (
	token *Token, err error) {
	var result Token
	err = doRest(
		"PUT",
		fmt.Sprintf("%s/%s/%s", tokenEndpoint, serviceAccountID, options.ID),
		s.client,
		doPayload(options),
		doResponse(&result))
	if err != nil {
		return
	}
	return &result, nil
}

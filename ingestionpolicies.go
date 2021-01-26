package wavefront

import (
	"fmt"
)

const basePolicyPath = "/api/v2/usage/ingestionpolicy"

type IngestionPolicy struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	UserAccountCount    int    `json:"userAccountCount"`
	ServiceAccountCount int    `json:"serviceAccountCount"`
}

type IngestionPolicies struct {
	client Wavefronter
}

func (c *Client) IngestionPolicies() *IngestionPolicies {
	return &IngestionPolicies{client: c}
}

func (p IngestionPolicies) Find(conditions []*SearchCondition) (
	results []*IngestionPolicy, err error) {
	err = doSearch(conditions, "ingestionpolicy", p.client, &results)
	return
}

func (p IngestionPolicies) Get(policy *IngestionPolicy) error {
	if policy.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", basePolicyPath, policy.ID),
		p.client,
		doResponse(policy))
}

func (p IngestionPolicies) Create(policy *IngestionPolicy) error {
	if policy.Name == "" {
		return fmt.Errorf("ingestion policy name must be specified")
	}
	return doRest(
		"POST",
		basePolicyPath,
		p.client,
		doPayload(policy),
		doResponse(policy))
}

func (p IngestionPolicies) Update(policy *IngestionPolicy) error {
	if policy.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", basePolicyPath, policy.ID),
		p.client,
		doPayload(policy),
		doResponse(policy))
}

func (p IngestionPolicies) Delete(policy *IngestionPolicy) error {
	if policy.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", basePolicyPath, policy.ID),
		p.client)
}

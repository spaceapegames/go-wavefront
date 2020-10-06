package wavefront

import (
	"fmt"
)

type ExternalLink struct {
	ID                    *string           `json:"id,omitempty"`
	Name                  string            `json:"name"`
	Description           string            `json:"description"`
	CreatorId             string            `json:"creatorId,omitempty"`
	UpdaterId             string            `json:"updaterId,omitempty"`
	UpdatedEpochMillis    int               `json:"updatedEpochMillis,omitempty"`
	CreatedEpochMillis    int               `json:"createdEpochMillis,omitempty"`
	Template              string            `json:"template"`
	MetricFilterRegex     string            `json:"metricFilterRegex,omitempty"`
	SourceFilterRegex     string            `json:"sourceFilterRegex,omitempty"`
	PointTagFilterRegexes map[string]string `json:"pointTagFilterRegexes,omitempty"`
	IsLogIntegration      bool              `json:"isLogIntegration,omitempty"`
}

const baseExtLinkPath = "/api/v2/extlink"

type ExternalLinks struct {
	client Wavefronter
}

func (c *Client) ExternalLinks() *ExternalLinks {
	return &ExternalLinks{client: c}
}

func (e ExternalLinks) Find(conditions []*SearchCondition) (
	results []*ExternalLink, err error) {
	err = doSearch(conditions, "extlink", e.client, &results)
	return
}

func (e ExternalLinks) Get(link *ExternalLink) error {
	if link.ID == nil || *link.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", baseExtLinkPath, *link.ID),
		e.client,
		doResponse(link))
}

func (e ExternalLinks) Create(link *ExternalLink) error {
	if link.Name == "" || link.Description == "" || link.Template == "" {
		return fmt.Errorf("externa link name, description, and template must be specified")
	}
	return doRest(
		"POST",
		baseExtLinkPath,
		e.client,
		doPayload(link),
		doResponse(link))
}

func (e ExternalLinks) Update(link *ExternalLink) error {
	if link.ID == nil || *link.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", baseExtLinkPath, *link.ID),
		e.client,
		doPayload(link),
		doResponse(link))
}

func (e ExternalLinks) Delete(link *ExternalLink) error {
	if *link.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	err := doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", baseExtLinkPath, *link.ID),
		e.client)
	if err != nil {
		return err
	}

	// Clear out the id to prevent re-submission
	empty := ""
	link.ID = &empty
	return nil
}

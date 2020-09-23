package wavefront

import (
	"encoding/json"
	"fmt"
)

type ExternalLink struct {
	ID                    *string           `json:"id"`
	Name                  string            `json:"name"`
	Description           string            `json:"description"`
	CreatorId             string            `json:"creatorId"`
	UpdaterId             string            `json:"updaterId"`
	UpdatedEpochMillis    int               `json:"updatedEpochMillis"`
	CreatedEpochMillis    int               `json:"createdEpochMillis"`
	Template              string            `json:"template"`
	MetricFilterRegex     string            `json:"metricFilterRegex,omitempty"`
	SourceFilterRegex     string            `json:"SourceFilterRegex,omitempty"`
	PointTagFilterRegexes map[string]string `json:"PointTagFilterRegexes,omitempty"`
}

const baseExtLinkPath = "/api/v2/extlink"

type ExternalLinks struct {
	client Wavefronter
}

func (c *Client) ExternalLinks() *ExternalLinks {
	return &ExternalLinks{client: c}
}

func (e ExternalLinks) Find(conditions []*SearchCondition) ([]*ExternalLink, error) {
	search := Search{
		client: e.client,
		Type:   "extlink",
		Params: &SearchParams{
			Conditions: conditions,
		},
	}

	var results []*ExternalLink
	moreItems := true
	for moreItems {
		resp, err := search.Execute()
		if err != nil {
			return nil, err
		}
		var tmpres []*ExternalLink
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

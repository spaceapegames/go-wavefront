package wavefront

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	if *link.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return e.crudExtLinks("GET", fmt.Sprintf("%s/%s", baseExtLinkPath, *link.ID), link)
}

func (e ExternalLinks) Create(link *ExternalLink) error {
	if link.Name == "" || link.Description == "" || link.Template == "" {
		return fmt.Errorf("externa link name, description, and template must be specified")
	}

	return e.crudExtLinks("POST", baseExtLinkPath, link)
}

func (e ExternalLinks) Update(link *ExternalLink) error {
	if *link.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	return e.crudExtLinks("POST", fmt.Sprintf("%s/%s", baseExtLinkPath, *link.ID), link)
}

func (e ExternalLinks) Delete(link *ExternalLink) error {
	if *link.ID == "" {
		return fmt.Errorf("id must be specified")
	}

	err := e.crudExtLinks("DELETE", fmt.Sprintf("%s/%s", baseExtLinkPath, *link.ID), link)
	if err != nil {
		return err
	}

	// Clear out the id to prevent re-submission
	*link.ID = ""
	return nil
}

func (e ExternalLinks) crudExtLinks(method, path string, extLink *ExternalLink) error {
	payload, err := json.Marshal(extLink)
	if err != nil {
		return err
	}

	request, err := e.client.NewRequest(method, path, nil, payload)
	if err != nil {
		return err
	}

	resp, err := e.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Close()

	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &struct {
		Response *ExternalLink `json:"response"`
	}{
		Response: extLink,
	})
}

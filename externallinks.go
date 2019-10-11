package wavefront

import (
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

type ExternalLinks struct {
	client Wavefronter
}

func (c *Client) ExternalLinks() *ExternalLinks {
	return &ExternalLinks{client: c}
}

func (e ExternalLinks) Find(searchConditions []*SearchCondition) ([]*ExternalLinks, error) {
	return nil, fmt.Errorf("not yet implemented")
}

func (e ExternalLinks) Get(link *ExternalLink) error {
	if *link.ID == "" {
		return fmt.Errorf("id must be specified")
	}
	return fmt.Errorf("not yet implemented")
}

func (e ExternalLinks) Create(link *ExternalLink) error {
	return fmt.Errorf("not yet implemented")
}

func (e ExternalLinks) Update(link *ExternalLink) error {
	if *link.ID == "" {
		return fmt.Errorf("id must be specified")
	}
	return fmt.Errorf("not yet implemented")
}

func (e ExternalLinks) Delete(link *ExternalLink) error {
	*link.ID = ""
	return fmt.Errorf("not yet implemented")
}

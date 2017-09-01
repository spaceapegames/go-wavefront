package wavefront

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Dashboard represents a single Wavefront Dashboard
type Dashboard struct {
	// Name is the name given to an Dashboard
	Name string `json:"name"`

	// ID is the Wavefront-assigned ID of an existing Dashboard
	ID string `json:"id"`

	// Tags are the tags applied to the Dashboard
	Tags []string

	// Description is a description given to the Dashboard
	Description string `json:"description"`

	// Url is the relative url to access the dashboard by on a cluster
	Url string `json:"url"`

	// Sections is an array of Section that split up the dashboard
	Sections []Section `json:"sections"`

	// ParameterDetails sets variables that can be used within queries
	ParameterDetails map[string]ParameterDetail `json:"parameterDetails"`
}

// ParameterDetail represents a parameter to dashboard that can be consumed in queries
type ParameterDetail struct {
	// Label represents the name of the variable
	Label string `json:"label"`

	// DefaultValue maps to keys in the map ValuesToReadableStrings
	DefaultValue string `json:"defaultValue"`

	// HideFromView Whether to hide from the view of the user viewing the Dashboard
	HideFromView bool `json:"hideFromView"`

	// ParameterType
	ParameterType string `json:"parameterType"`

	// ValuesToReadableStrings
	ValuesToReadableStrings map[string]string `json:"valuesToReadableStrings"`
}

// Section Represents a Single section within a Dashboard
type Section struct {
	// Name is the name given to this section
	Name string `json:"name"`

	// Rows is an array of Rows in this section
	Rows []Row `json:"rows"`
}

// Row represents a single Row withing a Section of a Wavefront Dashboard
type Row struct {
	// Name represents the display name of the Row
	Name string `json:"name"`

	// Charts is an array of Chart that this row contains
	Charts []Chart `json:"charts"`
}

// Chart represents a single Chart, on a single Row with in Section of a Wavefront Dashboard
type Chart struct {
	// Name is the name of a chart
	Name string `json:"name"`

	// Description is the description of the chart
	Description string `json:"description"`

	// Sources is an Array of Source
	Sources []Source `json:"sources"`

	// Units are the units to use for the y axis
	Units string `json:"units,omitempty"`
}

// Source represents a single Source for a Chart
type Source struct {
	// Name is the name given to the source
	Name string `json:"name"`

	// Query is a wavefront Query
	Query string `json:"query"`
}

// Dashboards is used to perform Dashboard-related operations against the Wavefront API
type Dashboards struct {
	// client is the Wavefront client used to perform Dashboard-related operations
	client Wavefronter
}

const baseDashboardPath = "/api/v2/dashboard"

// UnmarshalJSON is a custom JSON unmarshaller for an Dashboard, used in order to
// populate the Tags field in a more intuitive fashion
func (a *Dashboard) UnmarshalJSON(b []byte) error {
	type dashboard Dashboard
	temp := struct {
		Tags map[string][]string `json:"tags,omitempty"`
		*dashboard
	}{
		dashboard: (*dashboard)(a),
	}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}
	a.Tags = temp.Tags["customerTags"]
	return nil
}

func (a *Dashboard) MarshalJSON() ([]byte, error) {
	type dashboard Dashboard
	return json.Marshal(&struct {
		Tags map[string][]string `json:"tags,omitempty"`
		*dashboard
	}{
		Tags: map[string][]string{
			"customerTags": a.Tags,
		},
		dashboard: (*dashboard)(a),
	})
}

// Dashboards is used to return a client for Dashboard-related operations
func (c *Client) Dashboards() *Dashboards {
	return &Dashboards{client: c}
}

// Find returns all Dashboards filtered by the given search conditions.
// If filter is nil, all Dashboards are returned.
func (a Dashboards) Find(filter []*SearchCondition) ([]*Dashboard, error) {
	search := &Search{
		client: a.client,
		Type:   "dashboard",
		Params: &SearchParams{
			Conditions: filter,
		},
	}

	var results []*Dashboard
	moreItems := true
	for moreItems == true {
		resp, err := search.Execute()
		if err != nil {
			return nil, err
		}
		var tmpres []*Dashboard
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

// Create is used to create an Dashboard in Wavefront.
// If successful, the ID field of the Dashboard will be populated.
func (a Dashboards) Create(dashboard *Dashboard) error {
	return a.crudDashboard("POST", baseDashboardPath, dashboard)
}

// Update is used to update an existing Dashboard.
// The ID field of the Dashboard must be populated
func (a Dashboards) Update(dashboard *Dashboard) error {
	if dashboard.ID == "" {
		return fmt.Errorf("Dashboard id field not set")
	}

	return a.crudDashboard("PUT", fmt.Sprintf("%s/%s", baseDashboardPath, dashboard.ID), dashboard)

}

// Get is used to retrieve an existing Dashboard by ID.
// The ID field must be provided
func (a Dashboards) Get(dashboard *Dashboard) error {
	if dashboard.ID == "" {
		return fmt.Errorf("Dashboard id field is not set")
	}

	return a.crudDashboard("GET", fmt.Sprintf("%s/%s", baseDashboardPath, dashboard.ID), dashboard)
}

// Delete is used to delete an existing Dashboard.
// The ID field of the Dashboard must be populated
func (a Dashboards) Delete(dashboard *Dashboard) error {
	if dashboard.ID == "" {
		return fmt.Errorf("Dashboard id field not set")
	}

	err := a.crudDashboard("DELETE", fmt.Sprintf("%s/%s", baseDashboardPath, dashboard.ID), dashboard)
	if err != nil {
		return err
	}

	//reset the ID field so deletion is not attempted again
	dashboard.ID = ""
	return nil

}

func (a Dashboards) crudDashboard(method, path string, dashboard *Dashboard) error {
	payload, err := json.Marshal(dashboard)
	if err != nil {
		return err
	}
	req, err := a.client.NewRequest(method, path, nil, payload)
	if err != nil {
		return err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Close()

	body, err := ioutil.ReadAll(resp)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &struct {
		Response *Dashboard `json:"response"`
	}{
		Response: dashboard,
	})

}

package wavefront

import "fmt"

const (
	maintenanceWindowEndpoint = "/api/v2/maintenancewindow"
)

// MaintenanceWindow represents a maintenance window in Wavefront
type MaintenanceWindow struct {
	ID                              string   `json:"id"`
	RunningState                    string   `json:"runningState"`
	SortAttr                        int      `json:"sortAttr"`
	Reason                          string   `json:"reason"`
	CustomerId                      string   `json:"customerId"`
	RelevantCustomerTags            []string `json:"relevantCustomerTags"`
	Title                           string   `json:"title"`
	StartTimeInSeconds              int64    `json:"startTimeInSeconds"`
	EndTimeInSeconds                int64    `json:"endTimeInSeconds"`
	RelevantHostTags                []string `json:"relevantHostTags"`
	RelevantHostNames               []string `json:"relevantHostNames"`
	CreatorId                       string   `json:"creatorId"`
	UpdaterId                       string   `json:"updaterId"`
	CreatedEpochMillis              int64    `json:"createdEpochMillis"`
	UpdatedEpochMillis              int64    `json:"updatedEpochMillis"`
	RelevantHostTagsAnded           bool     `json:"relevantHostTagsAnded"`
	HostTagGroupHostNamesGroupAnded bool     `json:"hostTagGroupHostNamesGroupAnded"`
	EventName                       string   `json:"eventName"`
}

// Options returns the current options for this maintenance window
func (m *MaintenanceWindow) Options() *MaintenanceWindowOptions {
	return &MaintenanceWindowOptions{
		Reason:                          m.Reason,
		Title:                           m.Title,
		StartTimeInSeconds:              m.StartTimeInSeconds,
		EndTimeInSeconds:                m.EndTimeInSeconds,
		RelevantCustomerTags:            m.RelevantCustomerTags,
		RelevantHostTags:                m.RelevantHostTags,
		RelevantHostNames:               m.RelevantHostNames,
		RelevantHostTagsAnded:           m.RelevantHostTagsAnded,
		HostTagGroupHostNamesGroupAnded: m.HostTagGroupHostNamesGroupAnded,
	}
}

// MaintenanceWindowOptions represents the configurable options for a
// maintenance window.
type MaintenanceWindowOptions struct {
	// Required. The reason for the maintenance window.
	Reason string `json:"reason,omitempty"`

	// Required. The title of the maintenance window.
	Title string `json:"title,omitempty"`

	// Required. The start time of the maintenance window in seconds since 1 Jan 1970
	StartTimeInSeconds int64 `json:"startTimeInSeconds,omitempty"`

	// Required. The end time of the maintenance window in seconds since 1 Jan 1970
	EndTimeInSeconds int64 `json:"endTimeInSeconds,omitempty"`

	RelevantCustomerTags []string `json:"relevantCustomerTags"`
	RelevantHostTags     []string `json:"relevantHostTags,omitempty"`
	RelevantHostNames    []string `json:"relevantHostNames,omitempty"`

	RelevantHostTagsAnded           bool `json:"relevantHostTagsAnded,omitempty"`
	HostTagGroupHostNamesGroupAnded bool `json:"hostTagGroupHostNamesGroupAnded,omitempty"`
}

// MaintenanceWindows is used to perform maintenance window related operations
// against the Wavefront API
type MaintenanceWindows struct {
	// client is the Wavefront client used to perform target-related operations
	client Wavefronter
}

// MaintenanceWindows is used to return a client for maintenance window
// related operations
func (c *Client) MaintenanceWindows() *MaintenanceWindows {
	return &MaintenanceWindows{client: c}
}

// Find returns all maintenance windows filtered by the given search conditions.
// If filter is nil, all maintenance windows are returned.
func (m *MaintenanceWindows) Find(filter []*SearchCondition) (
	results []*MaintenanceWindow, err error) {
	err = doSearch(filter, "maintenancewindow", m.client, &results)
	return
}

// GetByID returns the MaintenanceWindow with given ID. If no such
// MaintenanceWindow exists, GetByID returns an error. The caller can call
// NotFound on err to determine whether or not the error is because the
// MaintenanceWindow doesn't exist.
func (m *MaintenanceWindows) GetByID(id string) (
	maintenanceWindow *MaintenanceWindow, err error) {
	var result MaintenanceWindow
	err = doRest(
		"GET",
		fmt.Sprintf("%s/%s", maintenanceWindowEndpoint, id),
		m.client,
		doResponse(&result))
	if err != nil {
		return
	}
	return &result, nil
}

// Create creates a MaintenanceWindow according to options and returns the
// newly created MaintenanceWindow.
func (m *MaintenanceWindows) Create(options *MaintenanceWindowOptions) (
	maintenanceWindow *MaintenanceWindow, err error) {
	var result MaintenanceWindow
	err = doRest(
		"POST",
		maintenanceWindowEndpoint,
		m.client,
		doPayload(options),
		doResponse(&result))
	if err != nil {
		return
	}
	return &result, nil
}

// Update updates a MaintenanceWindow according to options and returns the
// updated MaintenanceWindow.
func (m *MaintenanceWindows) Update(
	id string, options *MaintenanceWindowOptions) (
	maintenanceWindow *MaintenanceWindow, err error) {
	var result MaintenanceWindow
	err = doRest(
		"PUT",
		fmt.Sprintf("%s/%s", maintenanceWindowEndpoint, id),
		m.client,
		doPayload(conformOptionsForWavefrontAPI(*options)),
		doResponse(&result))
	if err != nil {
		return
	}
	return &result, nil
}

// DeleteByID deletes the MaintenanceWindow with given id.
func (m *MaintenanceWindows) DeleteByID(id string) error {
	return doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", maintenanceWindowEndpoint, id),
		m.client)
}

func conformOptionsForWavefrontAPI(
	options MaintenanceWindowOptions) *MaintenanceWindowOptions {
	if options.RelevantCustomerTags == nil {
		options.RelevantCustomerTags = []string{}
	}
	return &options
}

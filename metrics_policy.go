package wavefront

// MetricsPolicy represents the global metrics policy for a given Wavefront domain
type MetricsPolicy struct {
	PolicyRules        []PolicyRule `json:"policyRules,omitempty"`
	Customer           string       `json:"customer,omitempty"`
	UpdaterId          string       `json:"updaterId,omitempty"`
	UpdatedEpochMillis int          `json:"updatedEpochMillis,omitempty"`
}

type PolicyRule struct {
	Accounts    []PolicyUser      `json:"accounts,omitempty"`
	UserGroups  []PolicyUserGroup `json:"userGroups,omitempty"`
	Roles       []Role            `json:"roles,omitempty"`
	Name        string            `json:"name,omitempty"`
	Tags        []PolicyTag       `json:"tags,omitempty"`
	Description string            `json:"description,omitempty"`
	Prefixes    []string          `json:"prefixes,omitempty"`
	TagsAnded   bool              `json:"tagsAnded,omitempty"`
	AccessType  string            `json:"accessType,omitempty"`
}

type UpdateMetricsPolicyRequest struct {
	PolicyRules []PolicyRuleRequest `json:"policyRules,omitempty"`
}

type PolicyRuleRequest struct {
	AccountIds   []string    `json:"accounts,omitempty"`
	UserGroupIds []string    `json:"userGroups,omitempty"`
	RoleIds      []string    `json:"roles,omitempty"`
	Name         string      `json:"name,omitempty"`
	Tags         []PolicyTag `json:"tags,omitempty"`
	Description  string      `json:"description,omitempty"`
	Prefixes     []string    `json:"prefixes,omitempty"`
	TagsAnded    bool        `json:"tagsAnded,omitempty"`
	AccessType   string      `json:"accessType,omitempty"`
}

type PolicyTag struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type PolicyUser struct {
	// Unique ID for the user
	ID string `json:"id,omitempty"`
	// Name of the user
	Name string `json:"name,omitempty"`
}

type PolicyUserGroup struct {
	// Unique ID for the user group
	ID string `json:"id,omitempty"`
	// Name of the user group
	Name string `json:"name,omitempty"`
	// Description of the Group purpose
	Description string `json:"description,omitempty"`
}

// MetricsPolicyAPI is used to perform MetricsPolicy-related operations against the Wavefront API
type MetricsPolicyAPI struct {
	// client is the Wavefront client used to perform Dashboard-related operations
	client Wavefronter
}

const baseMetricsPolicyPath = "/api/v2/metricspolicy"

// MetricsPolicyAPI is used to return a client for MetricsPolicy-related operations
func (c *Client) MetricsPolicyAPI() *MetricsPolicyAPI {
	return &MetricsPolicyAPI{client: c}
}

func (m *MetricsPolicyAPI) Get() (*MetricsPolicy, error) {
	metricsPolicy := MetricsPolicy{}
	err := doRest(
		"GET",
		baseMetricsPolicyPath,
		m.client,
		doResponse(&metricsPolicy),
	)
	return &metricsPolicy, err
}

func (m *MetricsPolicyAPI) Update(policyRules *UpdateMetricsPolicyRequest) (*MetricsPolicy, error) {
	metricsPolicy := MetricsPolicy{}
	err := doRest(
		"PUT",
		baseMetricsPolicyPath,
		m.client,
		doPayload(policyRules),
		doResponse(&metricsPolicy),
	)
	return &metricsPolicy, err
}

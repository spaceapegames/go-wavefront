package wavefront

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockMetricsPolicyClient struct {
	Client
	T *testing.T
}

func (m *MockMetricsPolicyClient) Do(req *http.Request) (io.ReadCloser, error) {
	switch req.Method {
	case "GET":
		return testDo(m.T, req, "./fixtures/crud-metrics-policy-default-response.json", "GET", &MetricsPolicy{})

	case "PUT":
		return testDo(m.T, req, "./fixtures/crud-metrics-policy-response.json", "PUT", &UpdateMetricsPolicyRequest{})

	default:
		return nil, fmt.Errorf("unimplemented METHOD %s", req.Method)
	}
}

func TestMetricsPolicy_Get(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	m := &MetricsPolicyAPI{
		client: &MockMetricsPolicyClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
	id := "8bcffe68-5fcb-47fa-b935-ba7bc102b9a7"
	resp, err := m.Get()
	assert.Nil(t, err)
	assert.Equal(t, &MetricsPolicy{
		PolicyRules: []PolicyRule{{
			Accounts:    []PolicyUser{},
			UserGroups:  []PolicyUserGroup{{ID: id, Name: "Everyone", Description: "System group which contains all users"}},
			Roles:       []Role{},
			Name:        "Allow All Metrics",
			Tags:        []PolicyTag{},
			Description: "Predefined policy rule. Allows access to all metrics (timeseries, histograms, and counters) for all accounts. If this rule is removed, all accounts can access all metrics if there are no matching blocking rules.",
			Prefixes:    []string{"*"},
			TagsAnded:   false,
			AccessType:  "ALLOW",
		}},
		Customer:           "example",
		UpdaterId:          "system",
		UpdatedEpochMillis: 1603762170831,
	}, resp)
}

func TestMetricsPolicy_Put(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	m := &MetricsPolicyAPI{
		client: &MockMetricsPolicyClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
	id := "8bcffe68-5fcb-47fa-b935-ba7bc102b9a7"
	id2 := "7y6ffe68-5fcb-47fa-b935-ba7bc102b9a7"
	resp, err := m.Update(&UpdateMetricsPolicyRequest{PolicyRules: []PolicyRuleRequest{{
		AccountIds:   []string{},
		UserGroupIds: []string{id},
		RoleIds:      []string{},
		Name:         "Allow All Metrics",
		Tags:         []PolicyTag{},
		Description:  "Predefined policy rule. Allows access to all metrics (timeseries, histograms, and counters) for all accounts. If this rule is removed, all accounts can access all metrics if there are no matching blocking rules.",
		Prefixes:     []string{"*"},
		TagsAnded:    true,
		AccessType:   "ALLOW",
	},
		{
			AccountIds:   []string{},
			UserGroupIds: []string{},
			RoleIds:      []string{"abc123", "poi567"},
			Name:         "BLOCK Some Metrics by role",
			Tags:         []PolicyTag{{Key: "Custom", Value: "Value"}},
			Description:  "Scoped filter for roles.",
			Prefixes:     []string{"aa.*", "bb.*"},
			TagsAnded:    true,
			AccessType:   "BLOCK",
		},
		{
			AccountIds:   []string{id2},
			UserGroupIds: []string{},
			RoleIds:      []string{},
			Name:         "Allow Some Metrics by accounts",
			Tags:         []PolicyTag{{Key: "env", Value: "prod"}},
			Description:  "Scoped filter for users.",
			Prefixes:     []string{"*"},
			TagsAnded:    false,
			AccessType:   "ALLOW",
		},
	}})

	assert.Nil(t, err)
	assert.Equal(t, &MetricsPolicy{
		PolicyRules: []PolicyRule{{
			Accounts: []PolicyUser{},
			UserGroups: []PolicyUserGroup{{
				Name:        "Everyone",
				ID:          id,
				Description: "System group which contains all users",
			}},
			Roles:       []Role{},
			Name:        "Allow All Metrics",
			Tags:        []PolicyTag{},
			Description: "Predefined policy rule. Allows access to all metrics (timeseries, histograms, and counters) for all accounts. If this rule is removed, all accounts can access all metrics if there are no matching blocking rules.",
			Prefixes:    []string{"*"},
			TagsAnded:   false,
			AccessType:  "ALLOW",
		},
			{
				Accounts:   []PolicyUser{},
				UserGroups: []PolicyUserGroup{},
				Roles: []Role{
					{Name: "test-role1", ID: "abc123", Description: "misc"},
					{Name: "test-role2", ID: "poi567", Description: ""},
				},
				Name:        "BLOCK Some Metrics by role",
				Tags:        []PolicyTag{{Key: "Custom", Value: "Value"}},
				Description: "Scoped filter for roles.",
				Prefixes:    []string{"aa.*", "bb.*"},
				TagsAnded:   true,
				AccessType:  "BLOCK",
			},
			{
				Accounts: []PolicyUser{
					{ID: "test1@example.com", Name: "test1@example.com"},
					{ID: "test2@example.com", Name: "test2@example.com"}},
				UserGroups:  []PolicyUserGroup{},
				Roles:       []Role{},
				Name:        "Allow Some Metrics by account",
				Tags:        []PolicyTag{{Key: "env", Value: "prod"}},
				Description: "Scoped filter for users.",
				Prefixes:    []string{"*"},
				TagsAnded:   true,
				AccessType:  "BLOCK",
			},
		},
		Customer:           "example",
		UpdaterId:          "john.doe@example.com",
		UpdatedEpochMillis: 2603766170831,
	},
		resp)
}

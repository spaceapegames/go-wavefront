package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockRoleClient struct {
	Client
	T *testing.T
}

type MockCrudRoleClient struct {
	Client
	T      *testing.T
	method string
}

func (m *MockRoleClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/search-role-response.json", "POST", &SearchParams{})
}

func (m *MockCrudRoleClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/crud-role-response.json", m.method, &User{})
}

func TestRole_Find(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	r := &Roles{
		client: &MockRoleClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	roles, err := r.Find(nil)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 1, len(roles))
	assertEqual(t, "12345678-1234-cccc-dddd-123456789101", roles[0].ID)
	assertEqual(t, "Example Role", roles[0].Name)
	assertEqual(t, "test", roles[0].Customer)

	permissions := []string{
		AGENT_MANAGEMENT,
		ALERTS_MANAGEMENT,
		DERIVED_METRICS_MANAGEMENT,
	}

	foundPermission := false
	for _, v := range permissions {
		for _, p := range (*roles[0]).Permissions {
			if v == p {
				foundPermission = true
				break
			}
		}
		if !foundPermission {
			t.Errorf("expected to find %s permission on user", v)
		}
	}
}

func TestRole_CreateUpdateDelete(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	r := &Roles{
		client: &MockCrudRoleClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T:      t,
			method: "POST",
		},
	}

	role := &Role{
		Name:        "Example Role",
		Description: "Example Role Description",
		Permissions: []string{
			AGENT_MANAGEMENT,
			ALERTS_MANAGEMENT,
			DERIVED_METRICS_MANAGEMENT,
		},
	}

	assertEqual(t, "Example Role", role.Name)
	assertEqual(t, "Example Role Description", role.Description)

	foundPermission := false
	for _, v := range []string{AGENT_MANAGEMENT, ALERTS_MANAGEMENT, DERIVED_METRICS_MANAGEMENT} {
		for _, p := range role.Permissions {
			if v == p {
				foundPermission = true
				break
			}
		}
		if !foundPermission {
			t.Errorf("expected to find %s permission on user", v)
		}
	}

	r.client.(*MockCrudRoleClient).method = "PUT"
	var _ = r.Update(role)
	assertEqual(t, 3, len(role.Permissions))

	r.client.(*MockCrudRoleClient).method = "DELETE"
	var _ = r.Delete(role)
	assertEqual(t, "", role.ID)
}

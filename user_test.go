package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockUserClient struct {
	Client
	T *testing.T
}

type MockCrudUserClient struct {
	Client
	T      *testing.T
	method string
}

func (m *MockUserClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/search-user-response.json", "POST", &SearchParams{})
}

func TestUsers_Find(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	u := &Users{
		client: &MockUserClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	users, err := u.Find(nil)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 2, len(users))
	assertEqual(t, "thatoneperson@example.com", *users[0].ID)

	permissions := []string{
		AGENT_MANAGEMENT,
		ALERTS_MANAGEMENT,
		DASHBOARD_MANAGEMENT,
		EMBEDDED_CHARTS_MANAGEMENT,
		EVENTS_MANAGEMENT,
		EXTERNAL_LINKS_MANAGEMENT,
		HOST_TAG_MANAGEMENT,
		METRICS_MANAGEMENT,
		USER_MANAGEMENT,
		INTEGRATIONS_MANAGEMENT,
		DIRECT_INGESTION,
		BATCH_QUERY_PRIORITY,
		DERIVED_METRICS_MANAGEMENT,
	}

	foundPermission := false
	for _, v := range permissions {
		for _, p := range (*users[1]).Permissions {
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

func (m *MockCrudUserClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/crud-user-response.json", m.method, &User{})
}

func TestUsers_CreateUpdateDelete(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	u := &Users{
		client: &MockCrudUserClient{
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

	emailAddress := "someone+testing@example.com"
	newUser := &NewUserRequest{
		Permissions: []string{},
		Groups:      UserGroupsWrapper{},
	}
	user := &User{
		Customer:    "test",
		Permissions: []string{DERIVED_METRICS_MANAGEMENT},
		Groups:      UserGroupsWrapper{},
	}

	if err := u.Create(newUser, user, true); err == nil {
		t.Errorf("expected error user mising emailAddress")
	}

	newUser.EmailAddress = emailAddress
	var _ = u.Create(newUser, user, true)
	assertEqual(t, emailAddress, *user.ID)

	u.client.(*MockCrudUserClient).method = "PUT"
	user.ID = &emailAddress
	var _ = u.Update(user)
	assertEqual(t, 1, len(user.Permissions))

	u.client.(*MockCrudUserClient).method = "DELETE"
	var _ = u.Delete(user)
	assertEqual(t, "", *user.ID)
}

package wavefront

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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
	body, _ := ioutil.ReadAll(req.Body)
	search := SearchParams{}
	err := json.Unmarshal(body, &search)
	if err != nil {
		m.T.Fatal(err)
	}

	response, err := ioutil.ReadFile("./fixtures/search-user-response.json")
	if err != nil {
		m.T.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
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

	if len(users) != 2 {
		t.Errorf("expected two users returned, got %d", len(users))
	}

	if *users[0].ID != "thatoneperson@example.com" {
		t.Errorf("expected first user to be someone@example.com, got %s", *users[0].ID)
	}

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
	response, err := ioutil.ReadFile("./fixtures/crud-user-response.json")
	if err != nil {
		m.T.Fatal(err)
	}
	if req.Method != m.method {
		m.T.Errorf("request method expected '%s' got '%s'", m.method, req.Method)
	}
	// the delete call is sent with zero body in the request
	if req.Body != nil {
		body, _ := ioutil.ReadAll(req.Body)
		user := User{}
		err = json.Unmarshal(body, &user)
		if err != nil {
			m.T.Fatal(err)
		}
	}
	return ioutil.NopCloser(bytes.NewReader(response)), nil
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
	if *user.ID != emailAddress {
		t.Errorf("expected ID of %s, got %s", emailAddress, *user.ID)
	}

	u.client.(*MockCrudUserClient).method = "PUT"
	user.ID = &emailAddress
	var _ = u.Update(user)
	if len(user.Permissions) != 1 {
		t.Errorf("expected only a single permission on user, got %d", len(user.Permissions))
	}

	u.client.(*MockCrudUserClient).method = "DELETE"
	var _ = u.Delete(user)
	if *user.ID != "" {
		t.Errorf("expected user ID to be blank got %s", *user.ID)
	}
}

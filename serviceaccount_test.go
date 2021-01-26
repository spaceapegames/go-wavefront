package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockServiceAccountsClient struct {
	Client
	T *testing.T
}

type MockCrudServiceAccountsClient struct {
	Client
	T      *testing.T
	method string
}

func (client MockServiceAccountsClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(client.T, req, "./fixtures/search-serviceaccount-response.json", "POST", &SearchParams{})
}

func (client MockCrudServiceAccountsClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(client.T, req, "./fixtures/crud-serviceaccount-response.json", client.method, &ServiceAccount{})
}

func TestServiceAccounts_Find(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")

	testClient := &ServiceAccounts{
		client: &MockServiceAccountsClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	serviceAccounts, err := testClient.Find(nil)

	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, 1, len(serviceAccounts))
	assertEqual(t, "sa::someid", serviceAccounts[0].ID)
	assertEqual(t, "some-policy", serviceAccounts[0].IngestionPolicy.Name)
	assertEqual(t, "some-policy-1579800000000", serviceAccounts[0].IngestionPolicy.ID)
}

func testClient(t *testing.T) *ServiceAccounts {
	baseurl, _ := url.Parse("http://testing.wavefront.com")

	return &ServiceAccounts{
		client: &MockCrudServiceAccountsClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}
}

func TestServiceAccountsCreate(t *testing.T) {
	testClient := testClient(t)

	testClient.client.(*MockCrudServiceAccountsClient).method = "POST"
	options := &ServiceAccountOptions{}

	_, err := testClient.Create(options)

	if err != nil {
		t.Errorf("expected to receive error for missing fields")
	}

	options = &ServiceAccountOptions{}
	options.ID = "sa::tester"
	options.Description = "someDescription"

	_, err = testClient.Create(options)

	if err != nil {
		t.Fatal(err)
	}
}

func TestServiceAccountGetByID(t *testing.T) {
	testClient := testClient(t)
	testClient.client.(*MockCrudServiceAccountsClient).method = "GET"
	res, err := testClient.GetByID("sa::some_account")

	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "sa::some_account", res.ID)
}

func TestServiceAccountPut(t *testing.T) {
	testClient := testClient(t)
	testClient.client.(*MockCrudServiceAccountsClient).method = "PUT"

	options := &ServiceAccountOptions{}
	options.ID = "sa::tester"
	options.Description = "new description"

	_, err := testClient.Update(options)

	if err != nil {
		t.Fatal(err)
	}
}

func TestServiceAccountDeleteByID(t *testing.T) {
	testClient := testClient(t)
	testClient.client.(*MockCrudServiceAccountsClient).method = "DELETE"
	err := testClient.DeleteByID("sa::some_account")

	if err != nil {
		t.Fatal(err)
	}
}

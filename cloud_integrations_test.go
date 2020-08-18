package wavefront

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

type MockCloudIntegrationClient struct {
	Client
	T *testing.T
}

type MockCrudCloudIntegrationClient struct {
	Client
	method string
	T      *testing.T
}

type MockCloudIntegrationClientExtId struct {
	Client
	method string
	T      *testing.T
}

func (m *MockCloudIntegrationClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/search-cloud-integrations-response.json", "POST", &SearchParams{})
}

func (m *MockCloudIntegrationClientExtId) Do(req *http.Request) (io.ReadCloser, error) {
	var testType string
	return testDo(m.T, req, "./fixtures/aws-ext-id-cloud-integrations-response.json", m.method, &testType)
}

func (m *MockCrudCloudIntegrationClient) Do(req *http.Request) (io.ReadCloser, error) {
	return testDo(m.T, req, "./fixtures/crud-cloud-integrations-response.json", m.method, &CloudIntegration{})
}

func TestCloudIntegration_Search(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	c := &CloudIntegrations{
		client: &MockCloudIntegrationClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	cloudIntegrations, err := c.Find(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(cloudIntegrations) != 4 {
		t.Errorf("Expected to get 4 integrations got, %d", len(cloudIntegrations))
	}

	cloudWatch := cloudIntegrations[1]
	if cloudWatch.Service != "CLOUDWATCH" {
		t.Errorf("Expected Cloudwatch Integration, got %s", cloudWatch.Service)
	}

	if cloudWatch.LastErrorEvent == nil {
		t.Errorf("Expected error on cloudwatch integration")
	}
}

func TestCloudIntegrations_AwsExternalID(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	c := &CloudIntegrations{
		client: &MockCloudIntegrationClientExtId{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	c.client.(*MockCloudIntegrationClientExtId).method = "POST"
	extId, err := c.CreateAwsExternalID()
	if err != nil {
		t.Fatal(err)
	}

	if extId != "EXTERNAL-ID" {
		t.Errorf("Expected ext-id of EXTERNAL-ID, got %s", extId)
	}

	c.client.(*MockCloudIntegrationClientExtId).method = "GET"
	err = c.VerifyAwsExternalID("EXTERNAL-ID")
	if err != nil {
		t.Fatal(err)
	}

	c.client.(*MockCloudIntegrationClientExtId).method = "DELETE"
	err = c.DeleteAwsExternalID(&extId)

	if err != nil {
		t.Fatal(err)
	}

	if extId != "" {
		t.Errorf("Expected ext-id would be blank, got %s", extId)
	}
}

func TestCloudIntegrations_CRUD(t *testing.T) {
	baseurl, _ := url.Parse("http://testing.wavefront.com")
	c := &CloudIntegrations{
		client: &MockCrudCloudIntegrationClient{
			Client: Client{
				Config:     &Config{Token: "1234-5678-9977"},
				BaseURL:    baseurl,
				httpClient: http.DefaultClient,
				debug:      true,
			},
			T: t,
		},
	}

	c.client.(*MockCrudCloudIntegrationClient).method = "GET"
	cloudIntegration := &CloudIntegration{}
	err := c.Get(cloudIntegration)

	if err == nil {
		t.Errorf("Expected error missing CloudIntegration ID")
	}

	cloudIntegration.Id = "12345678-1234-5678-9101-12345678910111"
	if err = c.Get(cloudIntegration); err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "CLOUDWATCH", cloudIntegration.Service)
	assertEqual(t, "12345678-1234-5678-9101-12345678910111", cloudIntegration.Id)
	assertEqual(t, false, cloudIntegration.InTrash)
	assertEqual(t, "^aws.(ecs|rds|billing|instance|efs|datasync).*)$",
		cloudIntegration.CloudWatch.MetricFilterRegex)
	assertEqual(t, "arn:aws:iam::123456789012:role/example-arn",
		cloudIntegration.CloudWatch.BaseCredentials.RoleARN)

	c.client.(*MockCrudCloudIntegrationClient).method = "PUT"
	cloudIntegration.Name = "TESTING INTEGRATION"
	err = c.Update(cloudIntegration)
	if err != nil {
		t.Fatal(err)
	}

	c.client.(*MockCrudCloudIntegrationClient).method = "DELETE"
	err = c.Delete(cloudIntegration, false)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "", cloudIntegration.Id)
}

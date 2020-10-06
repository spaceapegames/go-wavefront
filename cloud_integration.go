package wavefront

import (
	"fmt"
	"strconv"
)

type CloudIntegration struct {
	ForceSave bool   `json:"forceSave,omitempty"`
	Id        string `json:"id,omitempty"`

	// The human-readable name of this integration
	Name string `json:"name"`

	// A value denoting which cloud service this integration integrates with
	Service string `json:"service"`

	// A value denoting whether or not this cloud integration is in the trashcan
	InTrash bool `json:"inTrash,omitempty"`

	// The creator of the cloud integration
	CreatorId string `json:"creatorId,omitempty"`

	// The user of the last person to update this integration
	UpdaterId string `json:"updaterId,omitempty"`

	// Event of the last error encountered by Wavefront servers when fetching data using this integration
	LastErrorEvent *Event `json:"lastErrorEvent,omitempty"`

	// A list of point tag key-values to add to every point ingested using this integration
	AdditionalTags map[string]string `json:"additionalTags,omitempty"`

	// Time that this integration last received a data point, in epoch millis
	LastReceivedDataPointMs int `json:"lastReceivedDataPointMs,omitempty"`

	// Number of metrics / events ingested by this integration the last time it ran
	LastMetricCount int `json:"lastMetricCount,omitempty"`

	CloudWatch       *CloudWatchConfiguration       `json:"cloudWatch,omitempty"`
	CloudTrail       *CloudTrailConfiguration       `json:"cloudTrail,omitempty"`
	EC2              *EC2Configuration              `json:"ec2,omitempty"`
	GCP              *GCPConfiguration              `json:"gcp,omitempty"`
	GCPBilling       *GCPBillingConfiguration       `json:"gcpBilling,omitempty"`
	NewRelic         *NewRelicConfiguration         `json:"newRelic,omitempty"`
	AppDynamics      *AppDynamicsConfiguration      `json:"appDynamics,omitempty"`
	Tesla            *TeslaConfiguration            `json:"tesla,omitempty"`
	Azure            *AzureConfiguration            `json:"azure,omitempty"`
	AzureActivityLog *AzureActivityLogConfiguration `json:"azureActivityLog,omitempty"`

	// Time, in epoch millis, of the last error encountered by Wavefront servers when fetching data using this integration.
	LastErrorMs int `json:"lastErrorMs,omitempty"`

	// True when an aws credential failed to authenticate
	Disabled bool `json:"disabled,omitempty"`

	// Opaque Id of the last Wavefront Integrations service to act on this integration
	LastProcessorId string `json:"lastProcessorId,omitempty"`

	// Time, in epoch millis, that this integration was last processed
	LastProcessingTimestamp int `json:"lastProcessingTimestamp,omitempty"`

	CreatedEpochMillis int `json:"createdEpochMillis,omitempty"`
	UpdatedEpochMillis int `json:"updatedEpochMillis,omitempty"`

	// Service refresh rate in minutes
	ServiceRefreshRateInMins int `json:"serviceRefreshRateInMins,omitempty"`

	Deleted bool `json:"deleted,omitempty"`
}

type CloudWatchConfiguration struct {
	// A regular expression that a CloudWatch metric name must match (case-insensitively) in order to be ingested
	MetricFilterRegex string `json:"metricFilterRegex,omitempty"`

	// A list of namespace that limit what we query from cloudwatch
	Namespaces []string `json:"namespaces,omitempty"`

	BaseCredentials *AWSBaseCredentials `json:"baseCredentials,omitempty"`

	// A string->string map of white list of AWS instance tag-value pairs (in AWS).
	// If the instance's AWS tags match this whitelist, CloudWatch data about this instance is ingested.
	// Multiple entries are OR'ed
	InstanceSelectionTags map[string]string `json:"instanceSelectionTags,omitempty"`

	// A string->string map of white list of AWS volume tag-value pairs (in AWS).
	// If the volume's AWS tags match this whitelist, CloudWatch data about this volume is ingested.
	// Multiple entries are OR'ed
	VolumeSelectionTags map[string]string `json:"volumeSelectionTags,omitempty"`

	// A regular expression that AWS tag key name must match (case-insensitively) in order to be ingested
	PointTagFilterRegex string `json:"pointTagFilterRegex,omitempty"`
}
type CloudTrailConfiguration struct {
	// The AWS Region of the S3 bucket where CloudTrail logs are stored
	Region string `json:"region"`

	// The common prefix, if any, appended to all CloudTrail log files
	Prefix string `json:"prefix,omitempty"`

	BaseCredentials *AWSBaseCredentials `json:"baseCredentials,omitempty"`

	// Name of the S3 bucket where CloudTrail logs are stored
	BucketName string `json:"bucketName"`

	// Rule to filter cloud trail log event
	FilterRule string `json:"filterRule,omitempty"`
}
type EC2Configuration struct {
	BaseCredentials *AWSBaseCredentials `json:"baseCredentials,omitempty"`

	// A list of AWS instance tags that, when found, will be used as the "source" name in a series.
	// Default: hostname, host, name
	// If no tag in this list is found, the series source is set to the instance id.
	HostNameTags []string `json:"hostNameTags"`
}
type GCPConfiguration struct {
	// A regular expression that a CloudWatch metric name must match (case-insensitively) in order to be ingested
	MetricFilterRegex string `json:"metricFilterRegex,omitempty"`

	// The Google Cloud Platform (GCP) project id.
	ProjectId string `json:"projectId"`

	// Private key for a Google Cloud Platform (GCP) service account within your project. The account must at least be
	// granted Monitoring Viewer permissions.  This key must be in the JSON format generated by GCP.
	// Use `{"project_id":"%s"}' to retain the existing key when updating
	GcpJSONKey string `json:"gcpJsonKey"`

	// A list of Google Cloud Platform (GCP) services (Such as ComputeEngine, PUbSub...etc) from which to pull metrics.
	CategoriesToFetch []string `json:"categoriesToFetch,omitempty"`
}
type GCPBillingConfiguration struct {
	// The Google Cloud Platform (GCP) project id.
	ProjectId string `json:"projectId"`

	// API key for Google Cloud Platform (GCP). Use 'saved_api_key' to retain existing API key when updating
	GcpApiKey string `json:"gcpApiKey"`

	// Private key for a Google Cloud Platform (GCP) service account within your project. The account must at least be
	// granted Monitoring Viewer permissions.  This key must be in the JSON format generated by GCP.
	// Use `{"project_id":"%s"}' to retain the existing key when updating
	GcpJSONKey string `json:"gcpJsonKey"`
}
type NewRelicConfiguration struct {
	// New Relic REST API Key
	ApiKey string `json:"apiKey"`

	// A regular expression that a application name must match (case-insensitively) in order to collect metrics
	AppFilterRegex string `json:"appFilterRegex,omitempty"`

	// A regular expression that a host name must match (case-insensitively) in order to collect metrics
	HostFilterRegex string `json:"hostFilterRegex,omitempty"`

	NewRelicMetricFilters []*NewRelicMetricFilters `json:"newRelicMetricFilters,omitempty"`
}
type AppDynamicsConfiguration struct {
	// Username is a combination of username and the account name
	UserName string `json:"userName"`

	// Name of the SaaS controller
	ControllerName string `json:"controllerName"`

	// Password for the AppDynamics user
	EncryptedPassword string `json:"encryptedPassword"`

	// Set this to 'false' to get separate results for all values with in the time range, by default it is true
	EnableRollup bool `json:"enableRollup,omitempty"`

	// Flag to control Error metric injections
	EnableErrorMetrics bool `json:"enableErrorMetrics,omitempty"`

	// Flag to control Business Transaction Metric injection
	EnableBusinessTrxMetrics bool `json:"enableBusinessTrxMetrics,omitempty"`

	// flag to control Backend metric injection
	EnableBackendMetrics bool `json:"enableBackendMetrics,omitempty"`

	// Flag to control Overall Performance metric injection
	EnableOverallPerfMetrics bool `json:"enableOverallPerfMetrics,omitempty"`

	// Flag to control individual Node metric injection
	EnableIndividualNodeMetrics bool `json:"enableIndividualNodeMetrics,omitempty"`

	// Flag to control Application Infrastructure metric injection
	EnableAppInfraMetrics bool `json:"enableAppInfraMetrics,omitempty"`

	// Flag to control Service End point metric injection
	EnableServiceEndpointMetrics bool `json:"enableServiceEndpointMetrics,omitempty"`

	// List of regular expressions that a application name must match (case-insensitively) in order to be ingested
	AppFilterRegex []string `json:"appFilterRegex"`
}
type TeslaConfiguration struct {
	// Email address for the Tesla account login
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AzureConfiguration struct {
	// A regular expression that a CloudWatch metric name must match (case-insensitively) in order to be ingested
	MetricFilterRegex string `json:"metricFilterRegex,omitempty"`

	BaseCredentials *AzureBaseCredentials `json:"baseCredentials,omitempty"`

	// A list of Azure services (such as Microsoft.Compute/virtualMachines,Microsoft.Cache/redis etc) from which to pull
	// metrics
	CategoryFilter []string `json:"categoryFilter,omitempty"`

	// A list of Azure resource groups from which to pull metrics
	ResourceGroupFilter []string `json:"resourceGroupFilter,omitempty"`
}
type AzureActivityLogConfiguration struct {
	BaseCredentials *AzureBaseCredentials `json:"baseCredentials,omitempty"`

	// A list of Azure ActivityLog categories to pull events for
	CategoryFilter []string `json:"categoryFilter,omitempty"`
}
type AWSBaseCredentials struct {
	// The Role ARN that the customer has created in AWS IAM to allow access to Wavefront
	RoleARN string `json:"roleArn"`

	// The external id corresponding to the Role ARN
	ExternalID string `json:"externalId"`
}
type NewRelicMetricFilters struct {
	AppName           string `json:"appName"`
	MetricFilterRegex string `json:"metricFilterRegex"`
}
type AzureBaseCredentials struct {
	// Client Id for an Azure service account within your project
	ClientID string `json:"clientId"`

	// Client Secret for an Azure service account within your project.  use `saved_secret` to retain the client secret
	// when updating
	ClientSecret string `json:"clientSecret"`

	// Tenant Id for an Azure service account within your project
	Tenant string `json:"tenant"`
}

const baseCloudIntegrationPath = "/api/v2/cloudintegration"

type CloudIntegrations struct {
	client Wavefronter
}

func (c *Client) CloudIntegrations() *CloudIntegrations {
	return &CloudIntegrations{client: c}
}

func (ci CloudIntegrations) Find(filter []*SearchCondition) (
	results []*CloudIntegration, err error) {
	err = doSearch(filter, "cloudintegration", ci.client, &results)
	return
}

// Get a CloudIntegration for a given ID
// ID must be specified
func (ci CloudIntegrations) Get(cloudIntegration *CloudIntegration) error {
	if cloudIntegration.Id == "" {
		return fmt.Errorf("cloud integration id must be specified")
	}
	return doRest(
		"GET",
		fmt.Sprintf("%s/%s", baseCloudIntegrationPath, cloudIntegration.Id),
		ci.client,
		doResponse(cloudIntegration))
}

// Deletes a given CloudIntegration and sets the ID of the object to ""
// ID must be specified
func (ci CloudIntegrations) Delete(cloudIntegration *CloudIntegration, skipTrash bool) error {
	if cloudIntegration.Id == "" {
		return fmt.Errorf("cloud integration id must be specified")
	}

	params := map[string]string{
		"skipTrash": strconv.FormatBool(skipTrash),
	}

	err := doRest(
		"DELETE",
		fmt.Sprintf("%s/%s", baseCloudIntegrationPath, cloudIntegration.Id),
		ci.client,
		doParams(params))
	if err == nil {
		cloudIntegration.Id = ""
	}
	return err
}

// Updates a given CloudIntegration in Wavefront
func (ci CloudIntegrations) Update(cloudIntegration *CloudIntegration) error {
	if cloudIntegration.Id == "" {
		return fmt.Errorf("cloud integration id must be specified")
	}
	return doRest(
		"PUT",
		fmt.Sprintf("%s/%s", baseCloudIntegrationPath, cloudIntegration.Id),
		ci.client,
		doPayload(cloudIntegration),
		doResponse(cloudIntegration))
}

// Creates a CloudIntegration in Wavefront
// If successful, the ID field will be populated
func (ci CloudIntegrations) Create(cloudIntegration *CloudIntegration) error {
	return doRest(
		"POST",
		baseCloudIntegrationPath,
		ci.client,
		doPayload(cloudIntegration),
		doResponse(cloudIntegration))
}

// Creates an AWS ExternalID for use in AWS IAM Roles
func (ci CloudIntegrations) CreateAwsExternalID() (string, error) {
	externalId := ""
	err := doRest(
		"POST",
		fmt.Sprintf("%s/awsExternalId", baseCloudIntegrationPath),
		ci.client,
		doResponse(&externalId))
	return externalId, err
}

// Deletes an AWS ExternalID
func (ci CloudIntegrations) DeleteAwsExternalID(externalId *string) error {
	err := doRest(
		"DELETE",
		fmt.Sprintf("%s/awsExternalId/%s", baseCloudIntegrationPath, *externalId),
		ci.client)
	if err == nil {
		*externalId = ""
	}
	return err
}

// Verifies an AWS ExternalID exists
func (ci CloudIntegrations) VerifyAwsExternalID(externalId string) error {
	return doRest(
		"GET",
		fmt.Sprintf("%s/awsExternalId/%s", baseCloudIntegrationPath, externalId),
		ci.client)
}

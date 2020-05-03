package main

import (
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"log"
)

func main() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	client, err := wavefront.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	cloudIntegrations := client.CloudIntegrations()

	// Create a new AWS CloudWatch Integration
	externalId, err := cloudIntegrations.CreateAwsExternalID()
	if err != nil {
		log.Fatal(err)
	}

	cloudWatchIntegration := wavefront.CloudIntegration{
		Name:    "example-cloud-watch-integration",
		Service: "CLOUDWATCH",
		CloudWatch: &wavefront.CloudWatchConfiguration{
			MetricFilterRegex: "^(ec2|elb).*$",
			BaseCredentials: &wavefront.AWSBaseCredentials{
				RoleARN:    "arn:aws:iam:1234567890:role/example-cloud-watch-integration",
				ExternalID: externalId,
			},
		},
	}

	err = cloudIntegrations.Create(&cloudWatchIntegration)
	if err != nil {
		log.Fatal(err)
	}

	// Get an integration by id
	err = cloudIntegrations.Get(&cloudWatchIntegration)
	if err != nil {
		log.Fatal(err)
	}

	// Update an existing integration
	cloudWatchIntegration.CloudWatch.InstanceSelectionTags = map[string]string{
		"env":  "prod",
		"role": "app",
	}

	err = cloudIntegrations.Update(&cloudWatchIntegration)

	if err != nil {
		log.Fatal(err)
	}

	// Delete an existing integration and bypass the trashcan
	err = cloudIntegrations.Delete(&cloudWatchIntegration, true)

	if err != nil {
		log.Fatal(err)
	}
}

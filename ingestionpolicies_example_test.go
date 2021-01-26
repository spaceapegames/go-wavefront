package wavefront_test

import (
	"fmt"
	"log"
	"time"

	"github.com/WavefrontHQ/go-wavefront-management-api"
)

func ExampleIngestionPolicies() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}

	client, err := wavefront.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	ingestionPolicies := client.IngestionPolicies()

	policy := &wavefront.IngestionPolicy{
		Name:        "test ingestion policy",
		Description: "an ingestion policy created by the Go SDK test suite",
	}

	err = ingestionPolicies.Create(policy)

	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete the policy
	fmt.Println("policy ID is", policy.ID)

	// Change the description
	policy.Description = "an ingestion policy updated by the Go SDK test suite"

	err = ingestionPolicies.Update(policy)

	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 60)

	err = ingestionPolicies.Delete(policy)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("policy deleted")
}

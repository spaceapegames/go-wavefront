package wavefront_test

import (
	"fmt"
	"log"
	"time"

	"github.com/WavefrontHQ/go-wavefront-management-api"
)

func ExampleExternalLinks() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	client, err := wavefront.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	externalLinks := client.ExternalLinks()

	el := &wavefront.ExternalLink{
		Name:              "Some External System",
		Description:       "A link to an external System",
		Template:          "https://example.com/source={{{source}}}&startTime={{startEpochMillis}}",
		MetricFilterRegex: "prod",
	}

	// Create the ExternalLink on Wavefront
	err = externalLinks.Create(el)
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete the ExternalLink
	fmt.Println("external link ID is", *el.ID)

	// Change to 10 minutes
	el.PointTagFilterRegexes = map[string]string{
		"region": "us-west-2",
	}

	// Update the ExternalLink
	err = externalLinks.Update(el)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 60)

	// Delete the ExternalLink
	err = externalLinks.Delete(el)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("external link deleted")

}

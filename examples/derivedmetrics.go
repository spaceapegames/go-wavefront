package main

import (
	"fmt"
	"log"
	"time"

	"github.com/WavefrontHQ/go-wavefront-management-api"
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

	derivedMetrics := client.DerivedMetrics()

	dm := &wavefront.DerivedMetric{
		Name:    "Some Derived Metric",
		Query:   "ts(someQuery)",
		Minutes: 5,
		Tags: wavefront.WFTags{
			CustomerTags: []string{"prod"},
		},
	}

	// Create the DerivedMetric on Wavefront
	err = derivedMetrics.Create(dm)
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete the DerivedMetric
	fmt.Println("derived metric ID is", *dm.ID)

	// Change to 10 minutes
	dm.Minutes = 10

	// Update the DerivedMetric
	err = derivedMetrics.Update(dm)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 60)

	// Delete the DerivedMetric
	err = derivedMetrics.Delete(dm, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("derived metric deleted")

}

package main

import (
	"fmt"
	"log"

	"github.com/spaceapegames/go-wavefront"
	"time"
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

	alerts := client.Alerts()

	a := &wavefront.Alert{
		Name:                "My First Alert",
		Target:              "test@example.com",
		Condition:           "ts(servers.cpu.usage, dc=dc2) > 10 * 10",
		DisplayExpression:   "ts(servers.cpu.usage, dc=dc2)",
		Minutes:             2,
		ResolveAfterMinutes: 2,
		Severity:            "WARN",
		Tags:                []string{"dc1", "synergy"},
	}

	// Create the alert on Wavefront
	err = alerts.Create(a)
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete the Alert
	fmt.Println("alert ID is", *a.ID)

	// Alternatively we could search for the Alert
	err = alerts.Get(&wavefront.Alert{
		ID: a.ID,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Update the Alert
	a.Target = "test@example.com,bob@example.com"
	err = alerts.Update(a)
	if err != nil {
		log.Fatal(err)
	}

	// Delete the Alert
	err = alerts.Delete(a)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("alert deleted")

}

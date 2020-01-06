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

	userGroups := client.UserGroups()
	// Which populates a User object upon success
	ug := &wavefront.UserGroup{
		Name: "Alert Users",
		Permissions: []string{
			wavefront.ALERTS_MANAGEMENT,
		},
		Users: []string{
			"someone@example.com",
		},
		Description: "Users can access & edit Wavefront Alerts",
	}

	// Create the UserGroup on Wavefront
	err = userGroups.Create(ug)
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete the UserGroup
	fmt.Println("user group ID is", *ug.ID)

	// Change to 10 minutes
	ug.Permissions = append(ug.Permissions, wavefront.DERIVED_METRICS_MANAGEMENT)

	// Update the User
	err = userGroups.Update(ug)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 60)

	// Delete the User
	err = userGroups.Delete(ug)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user group deleted")

}

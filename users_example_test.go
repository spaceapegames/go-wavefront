package wavefront_test

import (
	"fmt"
	"log"
	"time"

	"github.com/WavefrontHQ/go-wavefront-management-api"
)

func ExampleUsers() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	client, err := wavefront.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	users := client.Users()

	// Users are created with a NewUserRequest
	newUserRequest := &wavefront.NewUserRequest{
		EmailAddress: "someone@example.com",
		Permissions: []string{
			wavefront.ALERTS_MANAGEMENT,
		},
		Groups: wavefront.UserGroupsWrapper{},
	}

	// Which populates a User object upon success
	u := &wavefront.User{}

	// Create the User on Wavefront and don't send an email
	err = users.Create(newUserRequest, u, false)
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete the User
	fmt.Println("user ID is", *u.ID)

	// Change to 10 minutes
	u.Permissions = append(u.Permissions, wavefront.DERIVED_METRICS_MANAGEMENT)

	// Update the User
	err = users.Update(u)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 60)

	// Delete the User
	err = users.Delete(u)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user deleted")

}

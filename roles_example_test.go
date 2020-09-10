package wavefront_test

import (
	"log"

	"github.com/WavefrontHQ/go-wavefront-management-api"
)

func ExampleRoles() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	client, err := wavefront.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	client.Debug(true)

	roles := client.Roles()

	role := wavefront.Role{
		Name: "test role",
		Permissions: []string{
			wavefront.AGENT_MANAGEMENT,
			wavefront.ALERTS_MANAGEMENT,
			wavefront.DASHBOARD_MANAGEMENT,
		},
		Description: "testing something",
	}

	// Create the role on Wavefront
	err = roles.Create(&role)
	if err != nil {
		log.Fatal(err)
	}

	/**
	The following Add/Remove Assignees will return an error if the assignee does not exist in wavefront
	An assignee is either a UserGroup or User
	*/
	// Add an assignee
	err = roles.AddAssignees([]string{"user@example.com"}, &role)
	if err != nil {
		log.Fatal(err)
	}
	// Remove an assignee
	err = roles.RemoveAssignees([]string{"user@example.com"}, &role)

	// Revoke a permission
	err = roles.RevokePermission(wavefront.ALERTS_MANAGEMENT, []*wavefront.Role{&role})
	if err != nil {
		log.Fatal(err)
	}

	// Grant a permission
	err = roles.GrantPermission(wavefront.EVENTS_MANAGEMENT, []*wavefront.Role{&role})
	if err != nil {
		log.Fatal(err)
	}

	// Get an target by ID
	err = roles.Get(&wavefront.Role{ID: role.ID})
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update or delete the Target
	role.Description = "new description"
	err = roles.Update(&role)
	if err != nil {
		log.Fatal(err)
	}

	err = roles.Delete(&role)
	if err != nil {
		log.Fatal(err)
	}

}

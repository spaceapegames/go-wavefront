package wavefront_test

import (
	"io/ioutil"
	"log"

	"github.com/WavefrontHQ/go-wavefront-management-api"
)

func ExampleTargets() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	client, err := wavefront.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	client.Debug(true)

	targets := client.Targets()

	tmpl, _ := ioutil.ReadFile("./target-template.tmpl")

	target := wavefront.Target{
		Title:       "test target",
		Description: "testing something",
		Method:      "WEBHOOK",
		Recipient:   "https://hooks.slack.com/services/test/me",
		ContentType: "application/json",
		CustomHeaders: map[string]string{
			"Testing": "true",
		},
		Triggers: []string{"ALERT_OPENED", "ALERT_RESOLVED"},
		Template: string(tmpl),
	}

	// Create the target on Wavefront
	err = targets.Create(&target)
	if err != nil {
		log.Fatal(err)
	}

	// Get an target by ID
	err = targets.Get(&wavefront.Target{ID: target.ID})
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update or delete the Target
	target.Description = "new description"
	err = targets.Update(&target)
	if err != nil {
		log.Fatal(err)
	}

	err = targets.Delete(&target)
	if err != nil {
		log.Fatal(err)
	}

}

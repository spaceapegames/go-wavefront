package wavefront_test

import (
	"fmt"
	"log"
	"time"

	"github.com/WavefrontHQ/go-wavefront-management-api"
)

func ExampleServiceAccounts() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}

	wfClient, err := wavefront.NewClient(config)

	if err != nil {
		log.Fatal(err)
	}

	client := wfClient.ServiceAccounts()

	opts := &wavefront.ServiceAccountOptions{
		ID:                "sa::example",
		Description:       "an example service account",
		IngestionPolicyID: "example-policy-1579802191862",
	}

	obj, err := client.Create(opts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID is '%s' and description is '%s'.\n", obj.ID, obj.Description)
	time.Sleep(time.Second * 10)

	opts.Description = "a brand new description"
	obj, err = client.Update(opts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID is '%s' and description is now '%s'.\n", obj.ID, obj.Description)
	time.Sleep(time.Second * 10)

	err = client.DeleteByID(obj.ID)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Object deleted.")
}

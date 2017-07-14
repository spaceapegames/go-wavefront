package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spaceapegames/go-wavefront"
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

	events := client.Events()

	e := &wavefront.Event{
		Name:    "Deploying Stuff",
		Type:    "deployment",
		Details: "Stuff being deployed.",
		Tags:    []string{"prod"},
	}

	// Create the event on Wavefront
	err = events.Create(e)
	if err != nil {
		log.Fatal(err)
	}

	// The ID field is now set, so we can update/delete/close the Event
	fmt.Println("event ID is", *e.ID)

	// Close the Event
	err = events.Close(e)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 60)

	// Delete the Event
	err = events.Delete(e)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("event deleted")

}

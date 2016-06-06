package cmd

import (
	_ "fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var instant, end, delete bool
var hosts, evtStart string

// eventsCmd represents the events command
var cmdEvents = &cobra.Command{
	Use:   "events",
	Short: "Create Wavefront Events",
	Long: `Create, Close and Delete Wavefront Events.
For non-instantaneous Events, take note of the start time returned by the initial creation.
This should be passed with --start-time when closing or deleting the event`,
	Run: eventsCmd,
}

func init() {
	RootCmd.AddCommand(cmdEvents)

	cmdEvents.Flags().BoolVarP(&instant, "instant", "i", false, "Create an instantaneous Event")
	cmdEvents.Flags().BoolVarP(&delete, "delete", "r", false, "Delete an Event.")
	cmdEvents.Flags().BoolVarP(&end, "end", "e", false, "End the specified Event")
	cmdEvents.Flags().StringVarP(&hosts, "hosts", "a", "", "Comma-separated list of affected hosts.")
	cmdEvents.Flags().StringVarP(&evtStart, "start-time", "s", "", "Start time of event to Delete or End.")
}

func eventsCmd(cmd *cobra.Command, args []string) {
	createClient()

	if len(args) != 1 {
		log.Fatal("Please provide a name for the Event (or --help for usage)")
	}

	name := args[0]
	event := client.Events.NewEvent(name)

	if end == true || delete == true {
		if evtStart == "" {
			log.Fatal("Please provide a start time for deleting or ending events")
		}
		start, err := strconv.Atoi(evtStart)
		if err != nil {
			log.Fatal(err)
		}
		event.SetStartTime(int64(start))
	}

	if end == true && delete == true {
		log.Fatal("Please choose either --end-event or --delete")
	}

	for _, host := range strings.Split(hosts, ",") {
		event.AddAffectedHost(host)
	}

	var resp []byte
	var err error
	if instant == true {
		// create an Instantaneous Event
		if resp, err = event.Instant(); err != nil {
			log.Fatal(err)
		}
	} else if end == true {
		// end an event
		if resp, err = event.End(); err != nil {
			log.Fatal(err)
		}
	} else if delete == true {
		// delete an event
		if resp, err = event.Delete(); err != nil {
			log.Fatal(err)
		}
	} else {
		// start an event
		if resp, err = event.Create(); err != nil {
			log.Fatal(err)
		}
	}

	prettyPrint(&resp)
}

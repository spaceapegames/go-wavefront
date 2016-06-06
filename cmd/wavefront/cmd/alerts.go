package cmd

import (
	"fmt"
	"log"

	"github.com/buger/goterm"
	"github.com/spaceapegames/go-wavefront/api"
	"github.com/spf13/cobra"
)

var customerTag, userTag string

// cmdAlerts represents the alerts command
var cmdAlerts = &cobra.Command{
	Use:   "alerts [alert type]",
	Short: "Retrieve Wavefront Alerts",
	Long: `Poll Alerts endpoint for current Wavefront alerts.
Valid alert types are: all, snoozed, firing, invalid, affected_by_maintenance`,
	Run: alertsCmd,
}

func init() {
	RootCmd.AddCommand(cmdAlerts)

	cmdAlerts.Flags().StringVarP(&customerTag, "customer-tag", "c", "", "Filter by customer tag")
	cmdAlerts.Flags().StringVarP(&userTag, "user-tag", "u", "", "Filter by user tag")
}

func alertsCmd(cmd *cobra.Command, args []string) {
	createClient()

	params := wavefront.QueryParams{}

	if customerTag != "" {
		params["customerTag"] = customerTag
	}

	if userTag != "" {
		params["userTag"] = userTag
	}

	// parse the args to ascertain which function to use
	var f func(*wavefront.QueryParams) ([]*wavefront.Alert, error)
	if len(args) == 0 || args[0] == "all" {
		f = client.Alerts.All
	} else {
		switch args[0] {
		case "all":
			f = client.Alerts.All
		case "snoozed":
			f = client.Alerts.Snoozed
		case "firing":
			f = client.Alerts.Active
		case "invalid":
			f = client.Alerts.Invalid
		case "affected_by_maintenance":
			f = client.Alerts.AffectedByMaintenance
		default:
			log.Fatal("Please provide a valid alert type.")
		}
	}

	alerts, err := f(&params)
	if err != nil {
		log.Fatal(err)
	}

	if rawResponse == true {
		prettyPrint(&client.Alerts.RawResponse)
	} else {

		table := goterm.NewTable(0, 10, 5, ' ', 0)
		fmt.Fprintf(table, "Name\tSeverity\n")
		for _, a := range alerts {
			fmt.Fprintf(table, "%s\t%s\n", a.Name, a.Severity)
		}
		goterm.Println(table)
		goterm.Flush()
	}
}

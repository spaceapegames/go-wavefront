package wavefront_test

import (
	"fmt"
	"log"

	"github.com/WavefrontHQ/go-wavefront-management-api"
)

func ExampleDashboards() {
	config := &wavefront.Config{
		Address: "test.wavefront.com",
		Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
	}
	client, err := wavefront.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	dashboards := client.Dashboards()

	// Create the components of the Dashboard (Sections, Rows, Charts and Sources).

	sources := []wavefront.Source{
		{
			Name:  "Source 1",
			Query: "ts()",
		},
	}

	charts := []wavefront.Chart{
		{
			Name:        "Chart 1",
			Description: "Chart 1 shows ...",
			Sources:     sources,
			Units:       "units per time",
		},
	}

	rows := []wavefront.Row{
		{
			Name:   "Row 1",
			Charts: charts,
		},
	}

	sections := []wavefront.Section{
		{
			Name: "Section 1",
			Rows: rows,
		},
	}

	params := map[string]wavefront.ParameterDetail{
		"param": {
			Label:                   "test",
			DefaultValue:            "Label",
			HideFromView:            false,
			ParameterType:           "SIMPLE",
			ValuesToReadableStrings: map[string]string{"Label": "test"},
		},
	}

	d := &wavefront.Dashboard{
		Name:             "My First Dashboard",
		ID:               "dashboard1",
		Description:      "A Dashboard to show things",
		Url:              "dashboard1",
		Sections:         sections,
		Tags:             []string{"dc1", "synergy"},
		ParameterDetails: params,
	}

	// Create the dashboard on Wavefront
	err = dashboards.Create(d)
	if err != nil {
		log.Fatal(err)
	}

	// We can update/delete the Dashboard
	fmt.Println("dashboard ID is", d.ID)

	// Alternatively we could search for the Dashboard
	results, err := dashboards.Find(
		[]*wavefront.SearchCondition{
			{
				Key:            "name",
				Value:          "My First Dashboard",
				MatchingMethod: "EXACT",
			},
		})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("found dashboard with ID", results[0].ID)

	// Update the Dashboard
	d.Sections[0].Rows[0].Name = "Updated Chart Name"
	err = dashboards.Update(d)
	if err != nil {
		log.Fatal(err)
	}

	// Delete the Dashboard
	err = dashboards.Delete(d, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("dashboard deleted")

}

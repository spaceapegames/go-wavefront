package wavefront_test

import (
	"log"
	"os"
	"time"

	wavefront "github.com/WavefrontHQ/go-wavefront-management-api/writer"
)

func ExampleWriter() {
	// create a Writer with the source = hostname and
	// with all metrics exhibiting the tag environment=staging
	tags := []*wavefront.PointTag{
		{
			Key:   "environment",
			Value: "staging",
		},
	}

	source, _ := os.Hostname()

	wf, _ := wavefront.NewWriter(
		"wavefront-proxy.example.com",
		2878,
		source,
		tags,
	)
	defer wf.Close()

	// write a simple metric (timestamp now)
	err := wf.Write(wavefront.NewMetric("something.very.good.count", 33))
	if err != nil {
		log.Fatal(err)
	}

	// for more control over the metric (e.g. setting the timestamp and decimal places)
	m := wavefront.Metric{
		Name:      "something.very.good.count",
		Timestamp: time.Now().Unix() - 60*1000,
		Precision: 2,
	}
	err = wf.Write(m.Update(35.07))
	if err != nil {
		log.Fatal(err)
	}
}

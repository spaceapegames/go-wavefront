package main

import (
	"log"
	"os"

	wavefront "github.com/spaceapegames/go-wavefront/writer"
)

func main() {
	source, _ := os.Hostname()

	wf, err := wavefront.NewWriter("wavefront-proxy.example.com", 2878, source, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer wf.Close()

	wf.Write(wavefront.NewMetric("something.very.good.count", 33))
}

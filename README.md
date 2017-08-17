# Golang Wavefront Client [![GoDoc](https://godoc.org/github.com/spaceapegames/go-wavefront?status.svg)](https://godoc.org/github.com/spaceapegames/go-wavefront) [![Build Status](https://travis-ci.org/spaceapegames/go-wavefront.svg?branch=master)](https://travis-ci.org/spaceapegames/go-wavefront)

Golang SDK for interacting with the Wavefront v2 API, and sending metrics through a Wavefront proxy. 

## Usage 

### API Client
 
Presently support for querying, searching, alerts and events.

Please see the [examples](examples) directory for an example on how to use each, or check out the [documentation](https://godoc.org/github.com/spaceapegames/go-wavefront).

```Go
package main

import (
    "log"

    wavefront "github.com/spaceapegames/go-wavefront"
)

func main() {
    client, err := wavefront.NewClient{
        &wavefront.Config{
            Address: "test.wavefront.com",
            Token:   "xxxx-xxxx-xxxx-xxxx-xxxx",
        },
    }

    query := client.NewQuery(
        wavefront.NewQueryParams(`ts("cpu.load.1m.avg", dc=dc1)`),
    )

    if result, err := query.Execute(); err != nil {
        log.Fatal(err)
    }

    fmt.Println(result.TimeSeries[0].Label)
    fmt.Println(result.TimeSeries[0].DataPoints[0])
}
```

### Writer

Writer has full support for metric tagging etc.

Again, see [examples](examples) for a more detailed explanation.

```Go
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
```

## Contributing

Pull requests are welcomed. 

If you'd like to contribute to this project, please raise an issue and indicate that you'd like to take on the work prior to submitting a pull request. 

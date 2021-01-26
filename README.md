# Golang Wavefront Client

Golang SDK for interacting with the Wavefront v2 API, and sending metrics through a Wavefront proxy.

## Usage

### API Client

Supports most public API paths. Refer to the
[documentation](https://pkg.go.dev/github.com/WavefrontHQ/go-wavefront-management-api)
for further information.

```Go
package main

import (
    "fmt"
    "log"
    "os"

    wavefront "github.com/WavefrontHQ/go-wavefront-management-api"
)

func main() {
    client, err := wavefront.NewClient(
        &wavefront.Config{
            Address: os.Getenv("WAVEFRONT_ADDRESS"),
            Token:     os.Getenv("WAVEFRONT_TOKEN"),
        },
    )

    query := client.NewQuery(
        wavefront.NewQueryParams(`ts("cpu.load.1m.avg", dc=dc1)`),
    )

    result, err := query.Execute()

    if err != nil {
        log.Fatal(err)
    }

    if len(result.TimeSeries) > 0 {
        fmt.Println(result.TimeSeries[0].Label)
        fmt.Println(result.TimeSeries[0].DataPoints[0])
    } else {
        fmt.Println("No matching data.")
    }
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

    wavefront "github.com/WavefrontHQ/go-wavefront-management-api/writer"
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

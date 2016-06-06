package wavefront

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Metric is a representation of a single metric
type Metric struct {
	// Name is the full metric name (e.g. some.cool.count)
	Name string
	// Value is the numerical value of the metric
	Value string
	// Timestamp is the Unix epoch time at which the metric was measured.
	// If omitted, will take the time at point of writing metric
	Timestamp int64
}

type Writer struct {
	// A slice of type Metric
	Metrics []Metric
	// conn will hold the open connection to a Wavefront proxy agent
	conn net.Conn
	// Hostname of the sender. If omitted will default to os.Hostname()
	Hostname string
}

// NewWriter returns a Writer object
// Connection should be closed explicitly with Close()
func NewWriter(address string, port int) (*Writer, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	w := Writer{conn: conn, Metrics: []Metric{}, Hostname: hostname}
	return &w, nil
}

// AddMetric appends a metric to the list of those to be written.
func (w *Writer) AddMetric(m *Metric) {
	if m.Timestamp == 0 {
		m.Timestamp = time.Now().Unix()
	}
	w.Metrics = append(w.Metrics, *m)
}

// WriteMetrics flushes the list of Metrics, writing to the Wavefront proxy agent
func (w *Writer) WriteMetrics() {
	for _, m := range w.Metrics {
		fmt.Fprintf(w.conn, "%s %s %d host=%s\n", m.Name, m.Value, m.Timestamp, w.Hostname)
	}
	//reset Metrics
	w.Metrics = []Metric{}
}

// Write adds a metric and flushes all in one go
func (w *Writer) Write(m *Metric) {
	w.AddMetric(m)
	w.WriteMetrics()
}

func (w *Writer) Close() {
	w.conn.Close()
}

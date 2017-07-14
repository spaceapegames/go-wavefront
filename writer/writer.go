package wavefront

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

// Metric is a representation of a single metric
type Metric struct {
	// Name is the full metric name (e.g. some.cool.count)
	Name string

	// Value is the numerical value of the metric.
	Value float64

	// Precision is the number of decimal places to keep in sending the metric.
	// Defaults to 0
	Precision int

	// Timestamp is the Unix epoch seconds at which the metric was measured.
	// If omitted, will take the time at point at which the metric was received
	// by the Wavefront proxy
	Timestamp int64
}

// PointTag represents a metric point-tag, which is simply a key/value pair
type PointTag struct {
	// Key is the name of the point tag
	Key string

	// Value is the value of the point tag
	Value string
}

// Writer is used to write metrics to a Wavefront proxy.
// It will maintain an open connection until Close() is explicitly called.
// All metrics will be sent with the source and point-tags associated with
// the Writer.
type Writer struct {
	// conn will hold the open connection to a Wavefront proxy agent
	// It will generally be a net.Conn but is abstracted for testing
	conn io.WriteCloser

	// Source - metrics written by this Writer will have the 'source' set to this value
	source string

	// PointTags - metrics written by this Writer will have these tags
	pointTags []*PointTag

	// suffix is the suffix that will be applied to all metrics sent with this
	// Writer. It will be a combination of source and pointTags
	suffix string
}

// NewMetric returns a metric with the given name and value.
// The timestamp will default to the current time and the Precision to 0
// The value of the Metric can be updated with Update()
func NewMetric(name string, value float64) *Metric {
	return &Metric{
		Name:  name,
		Value: value,
	}
}

// Update updates the value of a metric
func (m *Metric) Update(value float64) *Metric {
	m.Value = value
	return m
}

// NewWriter returns a Writer object configured to send metrics to the address
// and port given.
// The Source of the Writer will be set to 'source'
// PointTags will be configured, pass nil if none required
// The Writer should be closed with Close() when no longer required
func NewWriter(address string, port int, source string, tags []*PointTag) (*Writer, error) {
	if source == "" {
		return nil, fmt.Errorf("source is required")
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}

	return &Writer{
		conn:      conn,
		source:    source,
		pointTags: tags,
		suffix:    metricSuffix(source, tags),
	}, nil
}

func metricSuffix(source string, tags []*PointTag) string {
	suffix := fmt.Sprintf("source=%s", source)
	for _, t := range tags {
		suffix = fmt.Sprintf("%s %s=%s", suffix, t.Key, t.Value)
	}
	return suffix + "\n"
}

// SetPointTags sets or updates the point tags with which metrics will be
// sent from this Writer
func (w *Writer) SetPointTags(tags []*PointTag) {
	w.pointTags = tags
	w.suffix = metricSuffix(w.source, tags)
}

// SetSource sets or updates the source with which metrics will be
// sent from this Writer
func (w *Writer) SetSource(source string) {
	w.source = source
	w.suffix = metricSuffix(source, w.pointTags)
}

// Write writes a metric to a Wavefront proxy
func (w *Writer) Write(m *Metric) error {
	format := "%s %." + strconv.Itoa(m.Precision) + "f"
	metric := fmt.Sprintf(format, m.Name, m.Value)
	if m.Timestamp != 0 {
		metric = fmt.Sprintf("%s %d", metric, m.Timestamp)
	}
	_, err := fmt.Fprintf(w.conn, "%s %s", metric, w.suffix)
	if err != nil {
		return err
	}

	return nil
}

// Close is used to close the connection to the Wavefront proxy
func (w *Writer) Close() {
	w.conn.Close()
}

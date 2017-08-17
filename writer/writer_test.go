package wavefront

import (
	"bytes"
	"fmt"
	"testing"
)

type dummyConn struct{}

var b bytes.Buffer

func (d dummyConn) Write(p []byte) (int, error) {
	fmt.Fprintf(&b, string(p))
	return 0, nil
}

func (d dummyConn) Close() error {
	return nil
}

func TestWriteMetrics(t *testing.T) {

	w := &Writer{
		conn:      dummyConn{},
		source:    "myHost1",
		pointTags: nil,
		suffix:    metricSuffix("myHost1", nil),
	}
	defer w.Close()

	updated := NewMetric("my.cool.test", 6969)
	updated.Update(7000)

	testMetrics := []struct {
		metric *Metric
		expect string
	}{
		{NewMetric("my.cool.test", 6969), "my.cool.test 6969 source=myHost1"},
		{updated, "my.cool.test 7000 source=myHost1"},
		{&Metric{"my.cool.test", 6969.00, 0, 0}, "my.cool.test 6969 source=myHost1"},
		{&Metric{"my.cool.test", 6969.00, 2, 1499695112}, "my.cool.test 6969.00 1499695112 source=myHost1"},
	}

	for _, tm := range testMetrics {
		var out []byte
		err := w.Write(tm.metric)
		if err != nil {
			t.Error(err)
		}
		out, err = b.ReadBytes('\n')
		if err != nil {
			t.Error(err)
		}
		if string(out) != tm.expect+"\n" {
			t.Errorf("metric, expected %s, got %s", tm.expect, string(out))
		}
	}

	w.SetPointTags([]*PointTag{
		&PointTag{Key: "some",
			Value: "tag",
		},
	})

	w.Write(NewMetric("my.cool.test", 7077))
	out, err := b.ReadBytes('\n')
	if err != nil {
		t.Error(err)
	}
	expect := "my.cool.test 7077 source=myHost1 some=tag"
	if string(out) != expect+"\n" {
		t.Errorf("point tags, expected %s, got %s", expect, string(out))
	}

	w.SetSource("anotherHost")
	w.Write(NewMetric("my.cool.test", 7077))
	out, err = b.ReadBytes('\n')
	if err != nil {
		t.Error(err)
	}
	expect = "my.cool.test 7077 source=anotherHost some=tag"
	if string(out) != expect+"\n" {
		t.Errorf("set source, expected %s, got %s", expect, string(out))
	}
}

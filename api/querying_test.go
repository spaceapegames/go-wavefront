package wavefront

import (
	_ "fmt"
	"strconv"
	"testing"
	"time"
)

func TestNewQuery(t *testing.T) {
	query := Querying{client: c}
	q := query.NewQuery("ts(dummy-query)", map[string]string{"s": "1461939366"})

	params := q.GetParams()

	if params["strict"] != "true" {
		t.Errorf("Default strict mode, expected true, got %s", params["strict"])
	}

	if params["s"] != "1461939366" {
		t.Errorf("Override start time, expected 1461939366, got %s", params["s"])
	}
}

func TestSetStartEndTime(t *testing.T) {
	query := Querying{client: c}
	q := query.NewQuery("ts(dummy-query)")

	i, _ := strconv.ParseInt("1461939366", 10, 64)
	tm := time.Unix(i, 0)

	q.SetEndTime(tm)
	params := q.GetParams()
	if params["e"] != "1461939366" {
		t.Errorf("SetEndTime, expected 1461939366, got %s", params["e"])
	}

	q.SetStartTime(LAST_HOUR)
	params = q.GetParams()
	if params["s"] != strconv.Itoa(1461939366-3600) {
		t.Errorf("SetStartTime, expected %s, got %s", strconv.Itoa(1461939366-3600), params["s"])
	}
}

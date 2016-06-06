package wavefront

import (
	_ "fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// test creating an Instant event, with a custom start time.
func TestCreateEvent(t *testing.T) {
	b := []byte(`{"Name":"Testing"}`)
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(b)

		if r.Method != "POST" {
			t.Fatalf("Expected POST request, got %s", r.Method)
		}

		r.ParseForm()
		name := r.Form.Get("n")
		if name != "Testing" {
			t.Errorf("Event name, expected Testing, got %s", name)
		}

		instant := r.Form.Get("c")
		if instant != "true" {
			t.Errorf("Event, expected instantaneous, got normal event")
		}

		startTime := r.Form.Get("s")
		if startTime != "1463501640308" {
			t.Errorf("Event start time, expected 1463501640308, got %s", startTime)
		}
	}))
	defer ts.Close()

	conf := &Config{Address: strings.TrimLeft(ts.URL, "https://"),
		Token:         "123456789",
		SkipTLSVerify: true}

	if tc, err := NewClient(conf); err != nil {
		t.Fatal(err)
	} else {
		event := Event{client: tc}
		e := event.NewEvent("Testing")
		e.SetStartTime(1463501640308)
		if resp, err := e.Instant(); err != nil {
			t.Fatalf("Creating new Event, failed with: %s", err)
		} else {
			if string(resp) != string(b) {
				t.Errorf("Event Create return value, expected %s, got %s", b, resp)
			}
		}
	}
}

func TestDeleteEvent(t *testing.T) {
	b := []byte(`{"Name":"Testing"}`)
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(b)

		if r.Method != "DELETE" {
			t.Fatalf("Expected DELETE request, got %s", r.Method)
		}
	}))
	defer ts.Close()

	conf := &Config{Address: strings.TrimLeft(ts.URL, "https://"),
		Token:         "123456789",
		SkipTLSVerify: true}

	if tc, err := NewClient(conf); err == nil {
		event := Event{client: tc}
		e := event.NewEvent("Testing Delete")
		e.SetStartTime(1463501640308)
		if _, err := e.Delete(); err != nil {
			t.Fatalf("Deleting Event failed with: %s, err")
		}
	}
}

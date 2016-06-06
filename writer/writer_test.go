package wavefront

import (
	"log"
	"net"
	"strings"
	"testing"
)

var ln = mockSocket()

func mockSocket() net.Listener {
	l, err := net.Listen("tcp", "localhost:31245")
	if err != nil {
		log.Fatalf("Unable to create listener: %s", err)
	}
	return l
}

func TestWriteMetrics(t *testing.T) {
	w, err := NewWriter("localhost", 31245)
	if err != nil {
		t.Fatalf("NewWriter failed with %s", err)
	}

	connChannel := make(chan string)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %s, err")
		}
		defer conn.Close()
		buf := make([]byte, 1000)
		if _, err := conn.Read(buf); err != nil && err.Error() != "EOF" {
			log.Fatalf("Error reading from socket: %s", err)
		} else {
			connChannel <- string(buf)
		}

	}()

	w.Write(&Metric{Name: "my.cool.test", Value: "6969"})
	output := <-connChannel
	if strings.Split(output, " ")[0] != "my.cool.test" {
		t.Errorf("metric name expected my.cool.test, got %s", strings.Split(output, " ")[0])
	}

	if strings.Split(output, " ")[1] != "6969" {
		t.Errorf("metric value expected 6969, got %s", strings.Split(output, " ")[1])
	}

}

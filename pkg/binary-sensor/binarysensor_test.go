package binarysensor

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

var (
	devAddr    = flag.String("device_addr", "6F:6F:6F:6F:01", "The esb device pipeline address")
	serverAddr = flag.String("server_addr", "localhost", "The server address")
	serverPort = flag.Uint("server_port", 9815, "The server port")
)

var b BinarySensor

func setup() {
	err := b.Open(*devAddr, *serverAddr, *serverPort)
	if err != nil {
		log.Fatalf("Setup: Connection Error: %v", err)
	}
}
func teardown() {
	err := b.Close()
	if err != nil {
		fmt.Printf("Error while disconnection: %v)", err)

	}
}

// TestGetValue tests reading a sensor channel
func TestGetValue(t *testing.T) {

	testChannel := uint8(0)

	val, err := b.GetValue(testChannel)

	if err != nil {
		t.Fatalf("%v", err)
	}

	fmt.Printf("Got Channel value: %v\n", val)
}

// TestSetValue tests setting a sensor channel by setting it and reading it back
func TestSetValue(t *testing.T) {

	testChannel := uint8(0)

	err := b.SetValue(testChannel, false)
	if err != nil {
		t.Fatalf("%v", err)
	}

	valExpectedFalse, _ := b.GetValue(testChannel)

	err = b.SetValue(testChannel, true)
	if err != nil {
		t.Fatalf("%v", err)
	}

	valExpectedTrue, _ := b.GetValue(testChannel)

	if valExpectedFalse || !valExpectedTrue {
		t.Fatalf("Unexpected values after set, got: %v (should be false), %v (should be true)", valExpectedFalse, valExpectedTrue)
	}
}

func TestSubscribe(t *testing.T) {
	s, err := b.Subscribe(0)
	if err != nil {
		t.Fatalf("%v", err)
	}

	select {
	case state := <-s.rxChan:
		fmt.Printf("Incoming Sensor state for channel %v: %v", s.channel, state)
	case <-time.After(10 * time.Second):
		t.Fatal("Timeout", err)
	}

}
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

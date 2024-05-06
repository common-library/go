package client_test

import (
	"os"
	"testing"
	"time"

	"github.com/common-library/go/database/prometheus/client"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

var addressForHttp = "http://" + os.Getenv("PROMETHEUS_ADDRESS")
var addressForHttps = "https://" + os.Getenv("PROMETHEUS_ADDRESS")

func stop() bool {
	if len(os.Getenv("PROMETHEUS_ADDRESS")) == 0 {
		return true
	} else {
		return false
	}
}

func TestNewClient(t *testing.T) {
	if stop() {
		return
	}

	if _, err := client.NewClient(addressForHttp); err != nil {
		t.Fatal(err)
	}
}

func TestNewClientWithBasicAuth(t *testing.T) {
	if stop() {
		return
	}

	if _, err := client.NewClientWithBasicAuth(addressForHttps, "username", "password"); err != nil {
		t.Fatal(err)
	}
}

func TestNewClientWithBearerToken(t *testing.T) {
	if stop() {
		return
	}

	if _, err := client.NewClientWithBearerToken(addressForHttps, "token"); err != nil {
		t.Fatal(err)
	}
}

func TestQuery(t *testing.T) {
	if stop() {
		return
	}

	if c, err := client.NewClient(addressForHttp); err != nil {
		t.Fatal(err)
	} else if value, warnings, err := c.Query("up", time.Now(), 10*time.Second); err != nil {
		t.Fatal(err)
	} else {
		t.Log(value)
		t.Log(warnings)
	}
}

func TestQueryRange(t *testing.T) {
	if stop() {
		return
	}

	r := v1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}

	if c, err := client.NewClient(addressForHttp); err != nil {
		t.Fatal(err)
	} else if value, warnings, err := c.QueryRange("rate(process_cpu_seconds_total[5m])", r, 10*time.Second); err != nil {
		t.Fatal(err)
	} else {
		t.Log(value)
		t.Log(warnings)
	}
}

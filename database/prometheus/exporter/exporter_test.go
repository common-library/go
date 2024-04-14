package exporter_test

import (
	"math/rand/v2"
	"strconv"
	"testing"
	"time"

	"github.com/common-library/go/database/prometheus/exporter"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNew(t *testing.T) {
	sample01Collector := exporter.New([]exporter.Metric{&metric01{}})

	if _, err := testutil.CollectAndLint(sample01Collector); err != nil {
		t.Fatal(err)
	}

	if count := testutil.CollectAndCount(sample01Collector); count != 3 {
		t.Fatal("invalid -", count)
	}

	if err := testutil.CollectAndCompare(sample01Collector, (&metric01{}).getExpected()); err != nil {
		t.Fatal(err)
	}
}

func TestRegister(t *testing.T) {
	sample01Collector := exporter.New([]exporter.Metric{&metric01{}})

	if err := exporter.Register(sample01Collector); err != nil {
		t.Fatal(err)
	}

	if exporter.UnRegister(sample01Collector) == false {
		t.Fatal("UnRegister false")
	}
}

func TestUnRegister(t *testing.T) {
	TestRegister(t)
}

func TestStart(t *testing.T) {
	address := ":" + strconv.Itoa(10000+rand.IntN(1000))
	path := "/metrics"

	sample01Collector := exporter.New([]exporter.Metric{&metric01{}})

	if err := exporter.Register(sample01Collector); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if exporter.UnRegister(sample01Collector) == false {
			t.Fatal("UnRegister false")
		}
	}()

	listenAndServeFailureFunc := func(err error) { t.Fatal(err) }
	if err := exporter.Start(address, path, listenAndServeFailureFunc); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := exporter.Stop(60); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	if err := testutil.ScrapeAndCompare("http://"+address+"/metrics", (&metric01{}).getExpected(), "sample01_metric01"); err != nil {
		t.Fatal(err)
	}
}

func TestStop(t *testing.T) {
	TestStart(t)
}

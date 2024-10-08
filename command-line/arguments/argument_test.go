package arguments_test

import (
	"flag"
	"os"
	"testing"

	"github.com/common-library/go/command-line/arguments"
)

func TestMain(m *testing.M) {
	setUp := func() {
		os.Args = []string{"test", "a1", "b2"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}

	tearDown := func() {}

	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func TestGet(t *testing.T) {
	answer := []string{"test", "a1", "b2"}

	for index, value := range answer {
		if arguments.Get(index) != value {
			t.Fatal(index, ",", value)
		}
	}
}

func TestGetAll(t *testing.T) {
	if args := arguments.GetAll(); args[0] != "test" || args[1] != "a1" || args[2] != "b2" {
		t.Fatal(arguments.GetAll())
	}
}

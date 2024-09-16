package flags_test

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/common-library/go/command-line/flags"
)

func set() error {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flagInfos := []flags.FlagInfo{
		{FlagName: "bool", Usage: "bool usage", DefaultValue: true},
		{FlagName: "time.Duration", Usage: "time.Duration usage (default 0h0m0s0ms0us0ns)", DefaultValue: time.Duration(0) * time.Second},
		{FlagName: "float64", Usage: "float64 usage (default 0)", DefaultValue: float64(0)},
		{FlagName: "int64", Usage: "int64 usage (default 0)", DefaultValue: int64(0)},
		{FlagName: "int", Usage: "int usage (default 0)", DefaultValue: int(0)},
		{FlagName: "string", Usage: "string usage (default \"\")", DefaultValue: string("")},
		{FlagName: "uint64", Usage: "uint64 usage (default 0)", DefaultValue: uint64(0)},
		{FlagName: "uint", Usage: "uint usage (default 0)", DefaultValue: uint(0)},
	}

	return flags.Parse(flagInfos)
}

func TestParse(t *testing.T) {
	os.Args = []string{"test"}
	if err := set(); err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"test"}

	flagInfos := []flags.FlagInfo{
		{FlagName: "invalid", Usage: "invalid usage", DefaultValue: int32(0)},
	}

	if err := flags.Parse(flagInfos); err.Error() != `this data type is not supported. - (int32)` {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	os.Args = []string{"test", "-bool=true", "-time.Duration=1h2m3s4ms5us6ns", "-float64=1", "-int64=2", "-int=3", "-string=a", "-uint64=4", "-uint=5"}

	if err := set(); err != nil {
		t.Fatal(err)
	}

	if value := flags.Get[bool]("bool"); value != true {
		t.Fatal(value)
	}

	if duration, err := time.ParseDuration("1h2m3s4ms5us6ns"); err != nil {
		t.Fatal(err)
	} else if value := flags.Get[time.Duration]("time.Duration"); value != duration {
		t.Fatal(value)
	}

	if value := flags.Get[float64]("float64"); value != 1 {
		t.Fatal(value)
	}

	if value := flags.Get[int64]("int64"); value != 2 {
		t.Fatal(value)
	}

	if value := flags.Get[int]("int"); value != 3 {
		t.Fatal(value)
	}

	if value := flags.Get[string]("string"); value != "a" {
		t.Fatal(value)
	}

	if value := flags.Get[uint64]("uint64"); value != 4 {
		t.Fatal(value)
	}

	if value := flags.Get[uint]("uint"); value != 5 {
		t.Fatal(value)
	}
}

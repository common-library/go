package flag_test

import (
	"flag"
	"os"
	"testing"
	"time"

	command_line_flag "github.com/heaven-chp/common-library-go/command-line/flag"
)

func set() error {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flagInfos := []command_line_flag.FlagInfo{
		{FlagName: "bool", Usage: "bool usage", DefaultValue: true},
		{FlagName: "time.Duration", Usage: "time.Duration usage (default 0h0m0s0ms0us0ns)", DefaultValue: time.Duration(0) * time.Second},
		{FlagName: "float64", Usage: "float64 usage (default 0)", DefaultValue: float64(0)},
		{FlagName: "int64", Usage: "int64 usage (default 0)", DefaultValue: int64(0)},
		{FlagName: "int", Usage: "int usage (default 0)", DefaultValue: int(0)},
		{FlagName: "string", Usage: "string usage (default \"\")", DefaultValue: string("")},
		{FlagName: "uint64", Usage: "uint64 usage (default 0)", DefaultValue: uint64(0)},
		{FlagName: "uint", Usage: "uint usage (default 0)", DefaultValue: uint(0)},
	}

	return command_line_flag.Parse(flagInfos)
}

func TestParse(t *testing.T) {
	os.Args = []string{"test"}
	if err := set(); err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"test"}

	flagInfos := []command_line_flag.FlagInfo{
		{FlagName: "invalid", Usage: "invalid usage", DefaultValue: int32(0)},
	}

	if err := command_line_flag.Parse(flagInfos); err.Error() != `this data type is not supported. - (int32)` {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	os.Args = []string{"test", "-bool=true", "-time.Duration=1h2m3s4ms5us6ns", "-float64=1", "-int64=2", "-int=3", "-string=a", "-uint64=4", "-uint=5"}

	if err := set(); err != nil {
		t.Fatal(err)
	}

	{
		value := command_line_flag.Get[bool]("bool")
		if value != true {
			t.Errorf("invalid value - (%t)", value)
		}
	}

	{
		value := command_line_flag.Get[time.Duration]("time.Duration")
		if duration, err := time.ParseDuration("1h2m3s4ms5us6ns"); err != nil {
			t.Fatal(err)
		} else if value != duration {
			t.Errorf("invalid value - (%#v)", value)
		}
	}

	{
		value := command_line_flag.Get[float64]("float64")
		if value != 1 {
			t.Errorf("invalid value - (%f)", value)
		}
	}

	{
		value := command_line_flag.Get[int64]("int64")
		if value != 2 {
			t.Errorf("invalid value - (%d)", value)
		}
	}

	{
		value := command_line_flag.Get[int]("int")
		if value != 3 {
			t.Errorf("invalid value - (%d)", value)
		}
	}

	{
		value := command_line_flag.Get[string]("string")
		if value != "a" {
			t.Errorf("invalid value - (%s)", value)
		}
	}

	{
		value := command_line_flag.Get[uint64]("uint64")
		if value != 4 {
			t.Errorf("invalid value - (%d)", value)
		}
	}

	{
		value := command_line_flag.Get[uint]("uint")
		if value != 5 {
			t.Errorf("invalid value - (%d)", value)
		}
	}
}

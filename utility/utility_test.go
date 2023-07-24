package utility_test

import (
	"testing"

	"github.com/heaven-chp/common-library-go/utility"
)

func TestTypeName(t *testing.T) {
	{
		typeName := utility.GetTypeName(int(1))
		if typeName != "int" {
			t.Errorf("invalid type name - (%s)", typeName)
		}
	}

	{
		typeName := utility.GetTypeName(string("test"))
		if typeName != "string" {
			t.Errorf("invalid type name - (%s)", typeName)
		}
	}

	{
		value := int(1)
		typeName := utility.GetTypeName(&value)
		if typeName != "*int" {
			t.Errorf("invalid type name - (%s)", typeName)
		}
	}
}

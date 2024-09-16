package utility_test

import (
	"testing"

	"github.com/common-library/go/utility"
)

func TestGetTypeName(t *testing.T) {
	if typeName := utility.GetTypeName(int(1)); typeName != "int" {
		t.Fatal(typeName)
	}

	if typeName := utility.GetTypeName(string("test")); typeName != "string" {
		t.Fatal(typeName)
	}

	value := int(1)
	if typeName := utility.GetTypeName(&value); typeName != "*int" {
		t.Fatal(typeName)
	}
}

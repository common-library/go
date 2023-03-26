package json_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/heaven-chp/common-library-go/file"
	"github.com/heaven-chp/common-library-go/json"
)

type test1Struct struct {
	Bool   bool   `json:"bool"`
	Int    int    `json:"int"`
	String string `json:"string"`
}

type test2Struct struct {
	Test1Struct test1Struct `json:"test1Struct"`

	ArrayString      []string      `json:"arrayString"`
	ArrayTest1Struct []test1Struct `json:"arrayTest1Struct"`
}

type sampleStruct struct {
	Test1Struct test1Struct `json:"test1"`
	Test2Struct test2Struct `json:"test2"`
}

func getTestData() sampleStruct {
	var sampleStruct sampleStruct

	sampleStruct.Test1Struct.Bool = true
	sampleStruct.Test1Struct.Int = 111
	sampleStruct.Test1Struct.String = "aaa"

	sampleStruct.Test2Struct.Test1Struct.Bool = true
	sampleStruct.Test2Struct.Test1Struct.Int = 222
	sampleStruct.Test2Struct.Test1Struct.String = "bbb"

	sampleStruct.Test2Struct.ArrayString = make([]string, 3)
	sampleStruct.Test2Struct.ArrayString[0] = "abc"
	sampleStruct.Test2Struct.ArrayString[1] = "def"
	sampleStruct.Test2Struct.ArrayString[2] = "ghi"

	sampleStruct.Test2Struct.ArrayTest1Struct = make([]test1Struct, 2)
	sampleStruct.Test2Struct.ArrayTest1Struct[0].Bool = true
	sampleStruct.Test2Struct.ArrayTest1Struct[0].Int = 333
	sampleStruct.Test2Struct.ArrayTest1Struct[0].String = "ccc"
	sampleStruct.Test2Struct.ArrayTest1Struct[1].Bool = false
	sampleStruct.Test2Struct.ArrayTest1Struct[1].Int = 444
	sampleStruct.Test2Struct.ArrayTest1Struct[1].String = "ddd"

	return sampleStruct
}

func checkTestData(sampleStruct sampleStruct) error {
	if sampleStruct.Test1Struct.Bool != true ||
		sampleStruct.Test1Struct.Int != 111 ||
		sampleStruct.Test1Struct.String != "aaa" {
		return errors.New(fmt.Sprintf("invalid data - data : (%#v)", sampleStruct.Test1Struct))
	}

	if sampleStruct.Test2Struct.Test1Struct.Bool != true ||
		sampleStruct.Test2Struct.Test1Struct.Int != 222 ||
		sampleStruct.Test2Struct.Test1Struct.String != "bbb" {
		return errors.New(fmt.Sprintf("invalid data - data : (%#v)", sampleStruct.Test2Struct.Test1Struct))
	}

	if len(sampleStruct.Test2Struct.ArrayString) != 3 ||
		sampleStruct.Test2Struct.ArrayString[0] != "abc" ||
		sampleStruct.Test2Struct.ArrayString[1] != "def" ||
		sampleStruct.Test2Struct.ArrayString[2] != "ghi" {
		return errors.New(fmt.Sprintf("invalid data - size : (%d), data : (%#v)", len(sampleStruct.Test2Struct.ArrayString), sampleStruct.Test2Struct.ArrayString))
	}

	if len(sampleStruct.Test2Struct.ArrayTest1Struct) != 2 ||
		sampleStruct.Test2Struct.ArrayTest1Struct[0].Bool != true ||
		sampleStruct.Test2Struct.ArrayTest1Struct[0].Int != 333 ||
		sampleStruct.Test2Struct.ArrayTest1Struct[0].String != "ccc" ||
		sampleStruct.Test2Struct.ArrayTest1Struct[1].Bool != false ||
		sampleStruct.Test2Struct.ArrayTest1Struct[1].Int != 444 ||
		sampleStruct.Test2Struct.ArrayTest1Struct[1].String != "ddd" {
		return errors.New(fmt.Sprintf("invalid data - size : (%d), data : (%#v)", len(sampleStruct.Test2Struct.ArrayTest1Struct), sampleStruct.Test2Struct.ArrayTest1Struct))
	}

	return nil
}

func TestToString(t *testing.T) {
	sampleStruct := getTestData()
	output, err := json.ToString(sampleStruct)
	if err != nil {
		t.Error(err)
	}

	jsonData, err := file.Read("./sample.json")
	if err != nil {
		t.Error(err)
	}

	compare := strings.Join(jsonData, "")
	compare = strings.Replace(compare, " ", "", -1)
	compare = strings.Replace(compare, "\t", "", -1)

	if output != compare {
		t.Errorf("invalid data - output : (%s), jsonData : (%s)", output, compare)
	}
}

func TestToStringIndent(t *testing.T) {
	sampleStruct := getTestData()

	output, err := json.ToStringIndent(sampleStruct, "", "\t")
	if err != nil {
		t.Error(err)
	}

	jsonData, err := file.Read("./sample.json")
	if err != nil {
		t.Error(err)
	}

	compare := strings.Join(jsonData, "\n")
	if output != compare {
		t.Errorf("invalid data - output : (%s), jsonData : (%s)", output, compare)
	}
}

func TestToStructFromString(t *testing.T) {
	jsonData, err := file.Read("./sample.json")
	if err != nil {
		t.Error(err)
	}

	var sampleStruct sampleStruct
	err = json.ToStructFromString(strings.Join(jsonData, ""), &sampleStruct)
	if err != nil {
		t.Error(err)
	}

	err = checkTestData(sampleStruct)
	if err != nil {
		t.Error(err)
	}
}

func TestToStructFromFile(t *testing.T) {
	var sampleStruct sampleStruct

	err := json.ToStructFromFile("./no_such_file", &sampleStruct)
	if err.Error() != "no such file - (./no_such_file)" {
		t.Error(err)
	}

	err = json.ToStructFromFile("./sample.json", &sampleStruct)
	if err != nil {
		t.Error(err)
	}

	err = checkTestData(sampleStruct)
	if err != nil {
		t.Error(err)
	}
}

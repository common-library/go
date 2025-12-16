package json_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/common-library/go/file"
	"github.com/common-library/go/json"
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

func getSampleStruct() sampleStruct {
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
		return fmt.Errorf("invalid data - data : (%#v)", sampleStruct.Test1Struct)
	}

	if sampleStruct.Test2Struct.Test1Struct.Bool != true ||
		sampleStruct.Test2Struct.Test1Struct.Int != 222 ||
		sampleStruct.Test2Struct.Test1Struct.String != "bbb" {
		return fmt.Errorf("invalid data - data : (%#v)", sampleStruct.Test2Struct.Test1Struct)
	}

	if len(sampleStruct.Test2Struct.ArrayString) != 3 ||
		sampleStruct.Test2Struct.ArrayString[0] != "abc" ||
		sampleStruct.Test2Struct.ArrayString[1] != "def" ||
		sampleStruct.Test2Struct.ArrayString[2] != "ghi" {
		return fmt.Errorf("invalid data - size : (%d), data : (%#v)", len(sampleStruct.Test2Struct.ArrayString), sampleStruct.Test2Struct.ArrayString)
	}

	if len(sampleStruct.Test2Struct.ArrayTest1Struct) != 2 ||
		sampleStruct.Test2Struct.ArrayTest1Struct[0].Bool != true ||
		sampleStruct.Test2Struct.ArrayTest1Struct[0].Int != 333 ||
		sampleStruct.Test2Struct.ArrayTest1Struct[0].String != "ccc" ||
		sampleStruct.Test2Struct.ArrayTest1Struct[1].Bool != false ||
		sampleStruct.Test2Struct.ArrayTest1Struct[1].Int != 444 ||
		sampleStruct.Test2Struct.ArrayTest1Struct[1].String != "ddd" {
		return fmt.Errorf("invalid data - size : (%d), data : (%#v)", len(sampleStruct.Test2Struct.ArrayTest1Struct), sampleStruct.Test2Struct.ArrayTest1Struct)
	}

	return nil
}

func TestToString(t *testing.T) {
	answer := `{"test1":{"bool":true,"int":111,"string":"aaa"},"test2":{"test1Struct":{"bool":true,"int":222,"string":"bbb"},"arrayString":["abc","def","ghi"],"arrayTest1Struct":[{"bool":true,"int":333,"string":"ccc"},{"bool":false,"int":444,"string":"ddd"}]}}`

	if data, err := json.ToString(getSampleStruct()); err != nil {
		t.Fatal(err)
	} else if data != answer {
		t.Fatal(data)
	}
}

func TestToStringIndent(t *testing.T) {
	answer := `{
    "test1": {
        "bool": true,
        "int": 111,
        "string": "aaa"
    },
    "test2": {
        "test1Struct": {
            "bool": true,
            "int": 222,
            "string": "bbb"
        },
        "arrayString": [
            "abc",
            "def",
            "ghi"
        ],
        "arrayTest1Struct": [
            {
                "bool": true,
                "int": 333,
                "string": "ccc"
            },
            {
                "bool": false,
                "int": 444,
                "string": "ddd"
            }
        ]
    }
}`

	if data, err := json.ToStringIndent(getSampleStruct(), "", "    "); err != nil {
		t.Fatal(err)
	} else if data != answer {
		t.Fatal(data)
	}
}

func TestConvertFromFile(t *testing.T) {
	jsonFile := t.Name() + ".json"
	defer os.Remove(jsonFile)

	if data, err := json.ToString(getSampleStruct()); err != nil {
		t.Fatal(err)
	} else if err := file.Write(jsonFile, data, 0600); err != nil {
		t.Fatal(err)
	} else if _, err := json.ConvertFromFile[sampleStruct]("./no_such_file"); err.Error() != "open ./no_such_file: no such file or directory" {
		t.Fatal(err)
	} else if sampleStruct, err := json.ConvertFromFile[sampleStruct](jsonFile); err != nil {
		t.Fatal(err)
	} else if err := checkTestData(sampleStruct); err != nil {
		t.Fatal(err)
	}
}

func TestConvertFromString(t *testing.T) {
	if data, err := json.ToString(getSampleStruct()); err != nil {
		t.Fatal(err)
	} else if sampleStruct, err := json.ConvertFromString[sampleStruct](strings.Join([]string{data}, "")); err != nil {
		t.Fatal(err)
	} else if err := checkTestData(sampleStruct); err != nil {
		t.Fatal(err)
	}
}

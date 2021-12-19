// Package json provides a conversion between json and struct.
package json

import (
	"encoding/json"
	"github.com/heaven-chp/common-library-go/file"
	"strings"
)

// ToString is convert struct to json.
//
// ex) output, err := ToString(testStruct)
func ToString(input interface{}) (string, error) {
	output, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// ToStringIndent is convert the structto json by apply prefix and indent.
//
// ex) output, err := ToStringIndent(testStruct, "", "\t")
func ToStringIndent(input interface{}, prefix string, indent string) (string, error) {
	output, err := json.MarshalIndent(input, prefix, indent)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// ToStructFromString is convert json string to struct.
//
// ex) err := ToStructFromString("{\"test\" : 111}", &testStruct)
func ToStructFromString(data string, result interface{}) error {
	return json.Unmarshal([]byte(data), result)
}

// ToStructFromFile is convert json file to struct.
//
// ex) err := ToStructFromFile("./sample.config", &sampleConfig)
func ToStructFromFile(fileName string, result interface{}) error {
	data, err := file.Read(fileName)
	if err != nil {
		return err
	}

	return ToStructFromString(strings.Join(data, ""), result)
}

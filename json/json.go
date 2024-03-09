// Package json provides a conversion between json and struct.
package json

import (
	"encoding/json"

	"github.com/heaven-chp/common-library-go/file"
)

// ToString is convert struct to json.
//
// ex) output, err := ToString(testStruct)
func ToString(input any) (string, error) {
	output, err := json.Marshal(input)
	return string(output), err
}

// ToStringIndent is convert the structto json by apply prefix and indent.
//
// ex) output, err := ToStringIndent(testStruct, "", "\t")
func ToStringIndent(input any, prefix string, indent string) (string, error) {
	output, err := json.MarshalIndent(input, prefix, indent)
	return string(output), err
}

// ConvertFromFile is convert json file to T.
//
// ex) t, err := json.ConvertFromFile[T]("./test.json")
func ConvertFromFile[T any](fileName string) (T, error) {
	if data, err := file.Read(fileName); err != nil {
		var t T
		return t, err
	} else {
		return ConvertFromString[T](data)
	}
}

// ConvertFromString is convert json string to T.
//
// ex) t, err := json.ConvertFromString[T](data)
func ConvertFromString[T any](data string) (T, error) {
	var t T

	err := json.Unmarshal([]byte(data), &t)

	return t, err
}

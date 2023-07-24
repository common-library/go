package utility

import "reflect"

// GetTypeName is get the type name.
//
// ex) typeName := utility.GetTypeName(int(1))
func GetTypeName(value interface{}) string {
	return reflect.TypeOf(value).String()
}

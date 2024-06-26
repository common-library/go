// Package utility provides utility implementations.
package utility

import "reflect"

// GetTypeName is get the type name.
//
// ex) typeName := utility.GetTypeName(int(1))
func GetTypeName(value any) string {
	return reflect.TypeOf(value).String()
}

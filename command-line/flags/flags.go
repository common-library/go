// Package flag provides command line flag
package flags

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/heaven-chp/common-library-go/utility"
)

// FlagInfo is a struct that has command line flags information.
type FlagInfo struct {
	FlagName     string
	Usage        string
	DefaultValue any

	valueOriginal any
	value         any
}

var result map[string]*FlagInfo

// Parse is parse the command line flags.
//
//	ex) err := flags.Parse([]flags.FlagInfo{
//			{FlagName: "bool", Usage: "bool usage", DefaultValue: true},
//			{FlagName: "time.Duration", Usage: "time.Duration usage (default 0h0m0s0ms0us0ns)", DefaultValue: time.Duration(0) * time.Second},
//			{FlagName: "float64", Usage: "float64 usage (default 0)", DefaultValue: float64(0)},
//			{FlagName: "int64", Usage: "int64 usage (default 0)", DefaultValue: int64(0)},
//			{FlagName: "int", Usage: "int usage (default 0)", DefaultValue: int(0)},
//			{FlagName: "string", Usage: "string usage (default \"\")", DefaultValue: string("")},
//			{FlagName: "uint64", Usage: "uint64 usage (default 0)", DefaultValue: uint64(0)},
//			{FlagName: "uint", Usage: "uint usage (default 0)", DefaultValue: uint(0)},
//		})
func Parse(flagInfos []FlagInfo) error {
	result = make(map[string]*FlagInfo)

	for index, flagInfo := range flagInfos {
		switch flagInfo.DefaultValue.(type) {
		case bool:
			flagInfos[index].valueOriginal = flag.Bool(flagInfo.FlagName, flagInfo.DefaultValue.(bool), flagInfo.Usage)
		case time.Duration:
			flagInfos[index].valueOriginal = flag.Duration(flagInfo.FlagName, flagInfo.DefaultValue.(time.Duration), flagInfo.Usage)
		case float64:
			flagInfos[index].valueOriginal = flag.Float64(flagInfo.FlagName, flagInfo.DefaultValue.(float64), flagInfo.Usage)
		case int64:
			flagInfos[index].valueOriginal = flag.Int64(flagInfo.FlagName, flagInfo.DefaultValue.(int64), flagInfo.Usage)
		case int:
			flagInfos[index].valueOriginal = flag.Int(flagInfo.FlagName, flagInfo.DefaultValue.(int), flagInfo.Usage)
		case string:
			flagInfos[index].valueOriginal = flag.String(flagInfo.FlagName, flagInfo.DefaultValue.(string), flagInfo.Usage)
		case uint64:
			flagInfos[index].valueOriginal = flag.Uint64(flagInfo.FlagName, flagInfo.DefaultValue.(uint64), flagInfo.Usage)
		case uint:
			flagInfos[index].valueOriginal = flag.Uint(flagInfo.FlagName, flagInfo.DefaultValue.(uint), flagInfo.Usage)
		default:
			return errors.New(fmt.Sprintf("this data type is not supported. - (%s)", utility.GetTypeName(flagInfo.DefaultValue)))
		}

		result[flagInfo.FlagName] = &flagInfos[index]
	}

	flag.Parse()

	for _, flagInfo := range result {
		switch flagInfo.DefaultValue.(type) {
		case bool:
			flagInfo.value = *flagInfo.valueOriginal.(*bool)
		case time.Duration:
			flagInfo.value = *flagInfo.valueOriginal.(*time.Duration)
		case float64:
			flagInfo.value = *flagInfo.valueOriginal.(*float64)
		case int64:
			flagInfo.value = *flagInfo.valueOriginal.(*int64)
		case int:
			flagInfo.value = *flagInfo.valueOriginal.(*int)
		case string:
			flagInfo.value = *flagInfo.valueOriginal.(*string)
		case uint64:
			flagInfo.value = *flagInfo.valueOriginal.(*uint64)
		case uint:
			flagInfo.value = *flagInfo.valueOriginal.(*uint)
		}
	}

	return nil

}

// Get is get the command line flags value.
//
// ex) value := flags.Get[int]("int")
func Get[T any](flagName string) T {
	return result[flagName].value.(T)
}

// Package command_line_argument provides command line argument
package command_line_argument

import (
	"errors"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/heaven-chp/common-library-go/utility"
)

var once sync.Once
var instance *commandLineArgument

func singleton() *commandLineArgument {
	once.Do(func() {
		instance = &commandLineArgument{}
	})

	return instance
}

// Get is get the command line argument.
//
// ex) value := command_line_argument.Get("int").(int)
func Get(flagName string) interface{} {
	return singleton().get(flagName)
}

// Set is set the command line arguments.
//
//	ex) err := command_line_argument.Set([]command_line_argument.CommandLineArgumentInfo{
//			{FlagName: "bool", Usage: "bool usage", DefaultValue: true},
//			{FlagName: "time.Duration", Usage: "time.Duration usage (default 0h0m0s0ms0us0ns)", DefaultValue: time.Duration(0) * time.Second},
//			{FlagName: "float64", Usage: "float64 usage (default 0)", DefaultValue: float64(0)},
//			{FlagName: "int64", Usage: "int64 usage (default 0)", DefaultValue: int64(0)},
//			{FlagName: "int", Usage: "int usage (default 0)", DefaultValue: int(0)},
//			{FlagName: "string", Usage: "string usage (default \"\")", DefaultValue: string("")},
//			{FlagName: "uint64", Usage: "uint64 usage (default 0)", DefaultValue: uint64(0)},
//			{FlagName: "uint", Usage: "uint usage (default 0)", DefaultValue: uint(0)},
//		})
func Set(infos []CommandLineArgumentInfo) error {
	return singleton().set(infos)
}

type commandLineArgument struct {
	infos map[string]*CommandLineArgumentInfo
}

func (this *commandLineArgument) get(flagName string) interface{} {
	return this.infos[flagName].value
}

func (this *commandLineArgument) set(infos []CommandLineArgumentInfo) error {
	this.infos = make(map[string]*CommandLineArgumentInfo)

	for index, info := range infos {
		switch info.DefaultValue.(type) {
		case bool:
			infos[index].valueOriginal = flag.Bool(info.FlagName, info.DefaultValue.(bool), info.Usage)
		case time.Duration:
			infos[index].valueOriginal = flag.Duration(info.FlagName, info.DefaultValue.(time.Duration), info.Usage)
		case float64:
			infos[index].valueOriginal = flag.Float64(info.FlagName, info.DefaultValue.(float64), info.Usage)
		case int64:
			infos[index].valueOriginal = flag.Int64(info.FlagName, info.DefaultValue.(int64), info.Usage)
		case int:
			infos[index].valueOriginal = flag.Int(info.FlagName, info.DefaultValue.(int), info.Usage)
		case string:
			infos[index].valueOriginal = flag.String(info.FlagName, info.DefaultValue.(string), info.Usage)
		case uint64:
			infos[index].valueOriginal = flag.Uint64(info.FlagName, info.DefaultValue.(uint64), info.Usage)
		case uint:
			infos[index].valueOriginal = flag.Uint(info.FlagName, info.DefaultValue.(uint), info.Usage)
		default:
			return errors.New(fmt.Sprintf("this data type is not supported. - (%s)", utility.GetTypeName(info.DefaultValue)))
		}

		this.infos[info.FlagName] = &infos[index]
	}

	flag.Parse()

	for _, info := range this.infos {
		switch info.DefaultValue.(type) {
		case bool:
			info.value = *info.valueOriginal.(*bool)
		case time.Duration:
			info.value = *info.valueOriginal.(*time.Duration)
		case float64:
			info.value = *info.valueOriginal.(*float64)
		case int64:
			info.value = *info.valueOriginal.(*int64)
		case int:
			info.value = *info.valueOriginal.(*int)
		case string:
			info.value = *info.valueOriginal.(*string)
		case uint64:
			info.value = *info.valueOriginal.(*uint64)
		case uint:
			info.value = *info.valueOriginal.(*uint)
		}
	}

	return nil
}

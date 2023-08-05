package utility

import (
	"errors"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// CallerInfo is a structure that has caller information.
type CallerInfo struct {
	PackageName  string
	FileName     string
	FunctionName string
	Line         int
	GoroutineID  int
}

// GetCallerInfo is get the caller information.
//
// ex) callerInfo, err := utility.GetCallerInfo()
func GetCallerInfo(numberOfStackFramesToAscend int) (CallerInfo, error) {
	pc, file, line, ok := runtime.Caller(numberOfStackFramesToAscend)
	if ok == false {
		return CallerInfo{}, errors.New("runtime.Caller() call fail")
	}

	_, fileName := path.Split(file)

	split := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	splitLen := len(split)

	packageName := ""
	functionName := split[splitLen-1]
	if split[splitLen-2][0] == '(' {
		functionName = split[splitLen-2] + "." + functionName
		packageName = strings.Join(split[0:splitLen-2], ".")
	} else {
		packageName = strings.Join(split[0:splitLen-1], ".")
	}

	var buffer [64]byte
	n := runtime.Stack(buffer[:], false)
	field := strings.Fields(strings.TrimPrefix(string(buffer[:n]), "goroutine "))[0]
	goroutineID, err := strconv.Atoi(field)
	if err != nil {
		return CallerInfo{}, err
	}

	return CallerInfo{PackageName: packageName, FileName: fileName, FunctionName: functionName, Line: line, GoroutineID: goroutineID}, nil
}

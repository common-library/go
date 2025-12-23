// Package utility provides general-purpose utility functions.
//
// This package offers various helper functions for common tasks including
// runtime introspection, type information, and network utilities.
//
// Features:
//   - Caller information retrieval (file, line, function, goroutine ID)
//   - Type name extraction
//   - CIDR network utilities
//
// Example:
//
//	callerInfo, _ := utility.GetCallerInfo(1)
//	fmt.Printf("Called from %s:%d\n", callerInfo.FileName, callerInfo.Line)
package utility

import (
	"errors"
	"path"
	"runtime"
	"strconv"
	"strings"
)

// CallerInfo is a GetCallerInfo that has caller information.
type CallerInfo struct {
	PackageName  string
	FileName     string
	FunctionName string
	Line         int
	GoroutineID  int
}

// GetCallerInfo retrieves information about the calling function.
//
// This function uses runtime introspection to gather details about the
// function that called it, including package, file, function name, line
// number, and goroutine ID.
//
// Parameters:
//   - numberOfStackFramesToAscend: Number of stack frames to skip
//     0 = GetCallerInfo itself
//     1 = Direct caller of GetCallerInfo
//     2 = Caller's caller, etc.
//
// Returns:
//   - CallerInfo: Struct containing caller details
//   - error: Error if stack frame retrieval fails
//
// The returned CallerInfo contains:
//   - PackageName: Full package path (e.g., "github.com/user/project/pkg")
//   - FileName: Base file name (e.g., "main.go")
//   - FunctionName: Function name (e.g., "main" or "(*Type).Method")
//   - Line: Line number
//   - GoroutineID: ID of the current goroutine
//
// Example:
//
//	func myFunction() {
//	    callerInfo, err := utility.GetCallerInfo(1)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("Called from: %s\n", callerInfo.FileName)
//	    fmt.Printf("Line: %d\n", callerInfo.Line)
//	    fmt.Printf("Function: %s\n", callerInfo.FunctionName)
//	    fmt.Printf("Goroutine: %d\n", callerInfo.GoroutineID)
//	}
//
// Example in logging:
//
//	func logWithCaller(message string) {
//	    info, _ := utility.GetCallerInfo(1)
//	    log.Printf("[%s:%d] %s", info.FileName, info.Line, message)
//	}
//
// Example for debugging:
//
//	func debugStack() {
//	    for i := 0; i < 5; i++ {
//	        info, err := utility.GetCallerInfo(i)
//	        if err != nil {
//	            break
//	        }
//	        fmt.Printf("#%d %s:%d %s\n", i, info.FileName, info.Line, info.FunctionName)
//	    }
//	}
func GetCallerInfo(numberOfStackFramesToAscend int) (CallerInfo, error) {
	pc, file, line, ok := runtime.Caller(numberOfStackFramesToAscend)
	if !ok {
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

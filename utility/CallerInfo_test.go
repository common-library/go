package utility_test

import (
	"sync"
	"testing"

	"github.com/common-library/go/utility"
)

func TestGetCallerInfo(t *testing.T) {
	wg := new(sync.WaitGroup)
	goroutineID := 0
	errorChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		if callerInfo, err := utility.GetCallerInfo(1); err != nil {
			errorChan <- err
			return
		} else {
			goroutineID = callerInfo.GoroutineID
		}
	}()
	wg.Wait()

	select {
	case err := <-errorChan:
		t.Fatal(err)
	default:
	}

	if callerInfo, err := utility.GetCallerInfo(1); err != nil {
		t.Fatal(err)
	} else if callerInfo.PackageName != "github.com/common-library/go/utility_test" {
		t.Fatal(callerInfo.PackageName)
	} else if callerInfo.FileName != "CallerInfo_test.go" {
		t.Fatal(callerInfo.FileName)
	} else if callerInfo.FunctionName != "TestGetCallerInfo" {
		t.Fatal(callerInfo.FunctionName)
	} else if callerInfo.Line != 34 {
		t.Fatal(callerInfo.Line)
	} else if callerInfo.GoroutineID == goroutineID {
		t.Fatal(callerInfo.GoroutineID, goroutineID)
	} else if callerInfo2, err := utility.GetCallerInfo(1); err != nil {
		t.Fatal(err)
	} else if callerInfo.GoroutineID != callerInfo2.GoroutineID {
		t.Fatal(callerInfo.GoroutineID, ",", callerInfo2.GoroutineID)
	}
}

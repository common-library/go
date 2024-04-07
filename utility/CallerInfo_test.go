package utility_test

import (
	"sync"
	"testing"

	"github.com/common-library/go/utility"
)

func TestGetCallerInfo(t *testing.T) {
	callerInfo, err := utility.GetCallerInfo(1)
	if err != nil {
		t.Fatal(err)
	}

	if callerInfo.PackageName != "github.com/common-library/go/utility_test" {
		t.Errorf("invalid package name - (%s)", callerInfo.PackageName)
	}

	if callerInfo.FileName != "CallerInfo_test.go" {
		t.Errorf("invalid file name - (%s)", callerInfo.FileName)
	}

	if callerInfo.FunctionName != "TestGetCallerInfo" {
		t.Errorf("invalid function name - (%s)", callerInfo.FunctionName)
	}

	if callerInfo.Line != 11 {
		t.Errorf("invalid line - (%d)", callerInfo.Line)
	}

	{
		callerInfo2, err := utility.GetCallerInfo(1)
		if err != nil {
			t.Fatal(err)
		}

		if callerInfo.GoroutineID != callerInfo2.GoroutineID {
			t.Errorf("invalid goroutine id - (%d)", callerInfo.GoroutineID)
		}
	}

	{
		wg := new(sync.WaitGroup)
		goroutineID := 0
		wg.Add(1)
		go func() {
			defer wg.Done()

			callerInfo, err := utility.GetCallerInfo(1)
			if err != nil {
				t.Fatal(err)
			}

			goroutineID = callerInfo.GoroutineID
		}()
		wg.Wait()

		if callerInfo.GoroutineID == goroutineID {
			t.Errorf("invalid goroutine id - (%d)(%d)", callerInfo.GoroutineID, goroutineID)
		}
	}
}

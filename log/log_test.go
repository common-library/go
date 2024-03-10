package log_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/heaven-chp/common-library-go/file"
	"github.com/heaven-chp/common-library-go/json"
	"github.com/heaven-chp/common-library-go/log"
)

func test(t *testing.T, level log.Level) {
	repeat := 1000
	count := map[log.Level]int{
		log.LevelTrace: 6,
		log.LevelDebug: 5,
		log.LevelInfo:  4,
		log.LevelWarn:  3,
		log.LevelError: 2,
		log.LevelFatal: 1,
	}
	answer := map[string]map[string]any{
		"TRACE": map[string]any{"msg": "message-01", "key-01": "value-01", "key-02": float64(1), "CallerInfo": map[string]any{"PackageName": "github.com/heaven-chp/common-library-go/log_test.test", "FileName": "log_test.go", "FunctionName": "func1", "Line": float64(52)}},
		"DEBUG": map[string]any{"msg": "message-02", "key-01": "value-02", "key-02": float64(2), "CallerInfo": map[string]any{"PackageName": "github.com/heaven-chp/common-library-go/log_test.test", "FileName": "log_test.go", "FunctionName": "func1", "Line": float64(53)}},
		"INFO":  map[string]any{"msg": "message-03", "key-01": "value-03", "key-02": float64(3), "CallerInfo": map[string]any{"PackageName": "github.com/heaven-chp/common-library-go/log_test.test", "FileName": "log_test.go", "FunctionName": "func1", "Line": float64(54)}},
		"WARN":  map[string]any{"msg": "message-04", "key-01": "value-04", "key-02": float64(4), "CallerInfo": map[string]any{"PackageName": "github.com/heaven-chp/common-library-go/log_test.test", "FileName": "log_test.go", "FunctionName": "func1", "Line": float64(55)}},
		"ERROR": map[string]any{"msg": "message-05", "key-01": "value-05", "key-02": float64(5), "CallerInfo": map[string]any{"PackageName": "github.com/heaven-chp/common-library-go/log_test.test", "FileName": "log_test.go", "FunctionName": "func1", "Line": float64(56)}},
		"FATAL": map[string]any{"msg": "message-06", "key-01": "value-06", "key-02": float64(6), "CallerInfo": map[string]any{"PackageName": "github.com/heaven-chp/common-library-go/log_test.test", "FileName": "log_test.go", "FunctionName": "func1", "Line": float64(57)}},
	}

	fileName := uuid.New().String()
	fileExtensionName := "log"
	fullName := fileName + "." + fileExtensionName
	defer file.Remove(fullName)

	testLog := log.Log{}
	defer testLog.Flush()

	testLog.SetOutputToFile(fileName, fileExtensionName, false)

	testLog.SetLevel(level)
	testLog.SetWithCallerInfo(true)

	wg := new(sync.WaitGroup)
	for i := 0; i < repeat; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			testLog.Trace("message-01", "key-01", "value-01", "key-02", 1)
			testLog.Debug("message-02", "key-01", "value-02", "key-02", 2)
			testLog.Info("message-03", "key-01", "value-03", "key-02", 3)
			testLog.Warn("message-04", "key-01", "value-04", "key-02", 4)
			testLog.Error("message-05", "key-01", "value-05", "key-02", 5)
			testLog.Fatal("message-06", "key-01", "value-06", "key-02", 6)
		}()
	}
	wg.Wait()
	testLog.Flush()

	if data, err := file.Read(fullName); err != nil {
		t.Fatal(err)
	} else {
		lines := strings.Split(data, "\n")
		lines = lines[0 : len(lines)-1]
		if len(lines) != count[level]*repeat {
			t.Fatal("invalid -", len(lines))
		}

		for _, line := range lines {
			if result, err := json.ConvertFromString[map[string]any](line); err != nil {
				t.Log(line)
				t.Fatal(err)
			} else {
				for key, value := range answer[result["level"].(string)] {
					if key != "CallerInfo" {
						if result[key] != value {
							t.Fatal("invalid -", key, ",", result[key], ",", value)
						}
					} else {
						for key1, value1 := range value.(map[string]any) {
							if result[key].(map[string]any)[key1] != value1 {
								t.Fatal("invalid -", key1, ",", result[key].(map[string]any)[key1], ",", value1)
							}
						}
					}
				}
			}
		}
	}
}

func TestTrace(t *testing.T) {
	test(t, log.LevelTrace)
}

func TestDebug(t *testing.T) {
	test(t, log.LevelDebug)
}

func TestInfo(t *testing.T) {
	test(t, log.LevelInfo)
}

func TestWarn(t *testing.T) {
	test(t, log.LevelWarn)
}

func TestError(t *testing.T) {
	test(t, log.LevelError)
}

func TestFatal(t *testing.T) {
	test(t, log.LevelFatal)
}

func TestFlush(t *testing.T) {
	testLog := log.Log{}
	defer testLog.Flush()
}

func TestGetLevel(t *testing.T) {
	testLog := log.Log{}

	testLog.SetLevel(log.LevelTrace)
	testLog.Flush()
	if testLog.GetLevel() != log.LevelTrace {
		t.Fatal("invalid -", testLog.GetLevel())
	}
}

func TestSetLevel(t *testing.T) {
	testLog := log.Log{}

	testLog.SetLevel(log.LevelTrace)
	testLog.Flush()
	if testLog.GetLevel() != log.LevelTrace {
		t.Fatal("invalid -", testLog.GetLevel())
	}
}

func TestSetOutputToStdout(t *testing.T) {
	testLog := log.Log{}
	defer testLog.Flush()

	testLog.SetOutputToStdout()
	testLog.Flush()

	testLog.Info("test")
}

func TestSetOutputToStderr(t *testing.T) {
	testLog := log.Log{}
	defer testLog.Flush()

	testLog.SetOutputToStderr()
	testLog.Flush()

	testLog.Info("test")
}

func TestSetOutputToFile(t *testing.T) {
	TestInfo(t)
}

func TestSetWithCallerInfo(t *testing.T) {
	testLog := log.Log{}
	defer testLog.Flush()

	testLog.SetWithCallerInfo(true)

	testLog.Info("test")
}

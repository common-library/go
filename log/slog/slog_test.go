package slog_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/common-library/go/file"
	"github.com/common-library/go/json"
	"github.com/common-library/go/log/slog"
)

func test(t *testing.T, level slog.Level) {
	t.Parallel()

	repeat := 100
	count := map[slog.Level]int{
		slog.LevelTrace: 6,
		slog.LevelDebug: 5,
		slog.LevelInfo:  4,
		slog.LevelWarn:  3,
		slog.LevelError: 2,
		slog.LevelFatal: 1,
	}
	answer := map[string]map[string]any{
		"TRACE": map[string]any{"msg": "message-01", "key-01": "value-01", "key-02": float64(1), "CallerInfo": map[string]any{"PackageName": "github.com/common-library/go/log/slog_test.test", "FileName": "slog_test.go", "FunctionName": "func1", "Line": float64(53)}},
		"DEBUG": map[string]any{"msg": "message-02", "key-01": "value-02", "key-02": float64(2), "CallerInfo": map[string]any{"PackageName": "github.com/common-library/go/log/slog_test.test", "FileName": "slog_test.go", "FunctionName": "func1", "Line": float64(54)}},
		"INFO":  map[string]any{"msg": "message-03", "key-01": "value-03", "key-02": float64(3), "CallerInfo": map[string]any{"PackageName": "github.com/common-library/go/log/slog_test.test", "FileName": "slog_test.go", "FunctionName": "func1", "Line": float64(55)}},
		"WARN":  map[string]any{"msg": "message-04", "key-01": "value-04", "key-02": float64(4), "CallerInfo": map[string]any{"PackageName": "github.com/common-library/go/log/slog_test.test", "FileName": "slog_test.go", "FunctionName": "func1", "Line": float64(56)}},
		"ERROR": map[string]any{"msg": "message-05", "key-01": "value-05", "key-02": float64(5), "CallerInfo": map[string]any{"PackageName": "github.com/common-library/go/log/slog_test.test", "FileName": "slog_test.go", "FunctionName": "func1", "Line": float64(57)}},
		"FATAL": map[string]any{"msg": "message-06", "key-01": "value-06", "key-02": float64(6), "CallerInfo": map[string]any{"PackageName": "github.com/common-library/go/log/slog_test.test", "FileName": "slog_test.go", "FunctionName": "func1", "Line": float64(58)}},
	}

	fileName := t.Name()
	fileExtensionName := "log"
	fullName := fileName + "." + fileExtensionName
	defer file.Remove(fullName)

	testLog := slog.Log{}
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
							t.Fatal(key, ",", result[key], ",", value)
						}
					} else {
						for key1, value1 := range value.(map[string]any) {
							if result[key].(map[string]any)[key1] != value1 {
								t.Fatal(key1, ",", result[key].(map[string]any)[key1], ",", value1)
							}
						}
					}
				}
			}
		}
	}
}

func TestTrace(t *testing.T) {
	test(t, slog.LevelTrace)
}

func TestDebug(t *testing.T) {
	test(t, slog.LevelDebug)
}

func TestInfo(t *testing.T) {
	test(t, slog.LevelInfo)
}

func TestWarn(t *testing.T) {
	test(t, slog.LevelWarn)
}

func TestError(t *testing.T) {
	test(t, slog.LevelError)
}

func TestFatal(t *testing.T) {
	test(t, slog.LevelFatal)
}

func TestFlush(t *testing.T) {
	t.Parallel()

	testLog := slog.Log{}
	defer testLog.Flush()
}

func TestGetLevel(t *testing.T) {
	t.Parallel()

	testLog := slog.Log{}

	testLog.SetLevel(slog.LevelTrace)
	testLog.Flush()
	if testLog.GetLevel() != slog.LevelTrace {
		t.Fatal(testLog.GetLevel())
	}
}

func TestSetLevel(t *testing.T) {
	t.Parallel()

	testLog := slog.Log{}

	testLog.SetLevel(slog.LevelTrace)
	testLog.Flush()
	if testLog.GetLevel() != slog.LevelTrace {
		t.Fatal(testLog.GetLevel())
	}
}

func TestSetOutputToStdout(t *testing.T) {
	t.Parallel()

	testLog := slog.Log{}
	defer testLog.Flush()

	testLog.SetOutputToStdout()
	testLog.Flush()

	testLog.Info("test")
}

func TestSetOutputToStderr(t *testing.T) {
	t.Parallel()

	testLog := slog.Log{}
	defer testLog.Flush()

	testLog.SetOutputToStderr()
	testLog.Flush()

	testLog.Info("test")
}

func TestSetOutputToFile(t *testing.T) {
	TestInfo(t)
}

func TestSetWithCallerInfo(t *testing.T) {
	t.Parallel()

	testLog := slog.Log{}
	defer testLog.Flush()

	testLog.SetWithCallerInfo(true)

	testLog.Info("test")
}

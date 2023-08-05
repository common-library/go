package log_test

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/heaven-chp/common-library-go/file"
	log "github.com/heaven-chp/common-library-go/log/file"
)

var once sync.Once
var fileLog *log.FileLog

func instance() *log.FileLog {
	once.Do(func() {
		fileLog = &log.FileLog{}
	})

	return fileLog
}

func logging(sequential bool, numberOfIterations int) {
	wg := new(sync.WaitGroup)

	job1 := func() {
		defer wg.Done()

		if sequential == false {
			instance().Flush()
			instance().SetChannelSize(rand.Intn(5) + 1)
		}

		instance().Critical(1, "a")
		instance().Error(2, "b")
		instance().Warning(3, "c")
		instance().Info(4, "d")
		instance().Debug(5, 5)
	}

	job2 := func() {
		defer wg.Done()

		if sequential == false {
			instance().Flush()
			instance().SetChannelSize(rand.Intn(5) + 1)
		}

		instance().Criticalf("(%d) (%s)", 1, "a")
		instance().Errorf("(%d) (%s)", 2, "b")
		instance().Warningf("(%d) (%s)", 3, "c")
		instance().Infof("(%d) (%s)", 4, "d")
		instance().Debugf("(%d) (%s)", 5, "e")
	}

	run := func(printCallerInfo bool, job func()) {
		instance().SetPrintCallerInfo(printCallerInfo)

		for i := 0; i < numberOfIterations; i++ {
			instance().Flush()
			instance().SetChannelSize(rand.Intn(5) + 1)

			wg.Add(1)
			if sequential {
				job()
			} else {
				go job()
			}
		}

		wg.Wait()
	}

	run(true, job1)
	run(false, job2)

	instance().Flush()
}

func resultCheck(t *testing.T, sequential bool, numberOfIterations int) {
	getCorrectAnswer := func(level string, functionName string, startLine int) []string {
		var levelInfo = map[string]int{
			"CRITICAL": 1,
			"ERROR":    2,
			"WARNING":  3,
			"INFO":     4,
			"DEBUG":    5}

		prefix := "FileLog_test.go:" + functionName + ":"
		correctAnswer1 := []string{
			prefix + strconv.Itoa(startLine) + "][CRITICAL] : 1 a",
			prefix + strconv.Itoa(startLine+1) + "][ERROR] : 2 b",
			prefix + strconv.Itoa(startLine+2) + "][WARNING] : 3 c",
			prefix + strconv.Itoa(startLine+3) + "][INFO] : 4 d",
			prefix + strconv.Itoa(startLine+4) + "][DEBUG] : 5 5"}

		correctAnswer2 := []string{
			"[CRITICAL] : (1) (a)",
			"[ERROR] : (2) (b)",
			"[WARNING] : (3) (c)",
			"[INFO] : (4) (d)",
			"[DEBUG] : (5) (e)"}

		all := []string{}
		for i := 0; i < numberOfIterations; i++ {
			all = append(all, correctAnswer1...)
		}
		for i := 0; i < numberOfIterations; i++ {
			all = append(all, correctAnswer2...)
		}

		correctAnswer := []string{}
		for index, value := range all {
			if index%(levelInfo["DEBUG"]) > levelInfo[level]-1 {
				continue
			}
			correctAnswer = append(correctAnswer, value)
		}

		return correctAnswer
	}

	getResult := func() ([]string, map[string]int) {
		setting := instance().GetSetting()

		temp, err := file.Read(setting.OutputPath + "/" + setting.FileNamePrefix + "_" + time.Now().Format("20060102") + ".log")
		if err != nil {
			t.Fatal(err)
		}

		sequentialResult := []string{}
		parallelResult := make(map[string]int)
		for _, value := range temp {
			split1 := strings.SplitN(value, "]", 2)[1]
			split2 := split1
			if strings.Contains(split2, "FileLog_test.go") {
				split2 = strings.SplitN(split1, ":", 2)[1]
			}

			sequentialResult = append(sequentialResult, split2)
			parallelResult[split2]++
		}

		return sequentialResult, parallelResult
	}

	setting := instance().GetSetting()
	correctAnswer := getCorrectAnswer(setting.Level, "func1", 38)

	sequentialResult, parallelResult := getResult()

	if sequential {
		if len(sequentialResult) != len(correctAnswer) {
			t.Fatalf("invalid len - (%d)(%d)", len(sequentialResult), len(correctAnswer))
		}

		for index, value := range correctAnswer {
			if sequentialResult[index] != value {
				t.Fatalf("invalid log - (%s)(%s)", sequentialResult[index], value)
			}
		}
	} else {
		totalCount := 0

		for _, value := range correctAnswer {
			totalCount += parallelResult[value]

			if parallelResult[value] != numberOfIterations {
				t.Log(parallelResult)
				t.Fatalf("invalid count - (%s)(%d)(%d)", value, parallelResult[value], numberOfIterations)
			}
		}

		if totalCount/numberOfIterations != len(correctAnswer) {
			t.Fatalf("invalid len - (%d)(%d)", totalCount/numberOfIterations, len(correctAnswer))
		}
	}
}

func test(t *testing.T, level string) {
	job := func(sequential bool) {
		instance().SetOutputPath(t.TempDir())

		const numberOfIterations = 100
		logging(sequential, numberOfIterations)
		resultCheck(t, sequential, numberOfIterations)
	}

	instance().SetLevel(level)
	instance().SetFileNamePrefix("test_" + level)
	instance().SetChannelSize(rand.Intn(1024) + 1)

	job(true)
	job(false)
}

func setUp() {
	setting := log.Setting{
		Level:           "DEBUG",
		OutputPath:      "/tmp/log/",
		FileNamePrefix:  "test",
		PrintCallerInfo: true,
		ChannelSize:     1024}
	if err := instance().Initialize(setting); err != nil {
		panic(err)
	}
}

func tearDown() {
	if err := instance().Finalize(); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setUp()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func TestInitialize(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fileLog := log.FileLog{}

		err := fileLog.Initialize(log.Setting{Level: "invalid"})
		if err.Error() != "please select one of CRITICAL, ERROR, WARNING, INFO, DEBUG" {
			t.Fatalf("invalid error - (%s)", err.Error())
		}

		err = fileLog.Initialize(log.Setting{
			Level:           "DEBUG",
			OutputPath:      t.TempDir(),
			FileNamePrefix:  "",
			PrintCallerInfo: false,
			ChannelSize:     1024})
		if err != nil {
			t.Fatal(err)
		}

		if err := fileLog.Finalize(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestFinalize(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fileLog := log.FileLog{}

		if err := fileLog.Finalize(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCritical(t *testing.T) {
	test(t, "CRITICAL")
}

func TestCriticalf(t *testing.T) {
	test(t, "CRITICAL")
}

func TestError(t *testing.T) {
	test(t, "ERROR")
}

func TestErrorf(t *testing.T) {
	test(t, "ERROR")
}

func TestWarning(t *testing.T) {
	test(t, "WARNING")
}

func TestWarningf(t *testing.T) {
	test(t, "WARNING")
}

func TestInfo(t *testing.T) {
	test(t, "INFO")
}

func TestInfof(t *testing.T) {
	test(t, "INFO")
}

func TestDebug(t *testing.T) {
	test(t, "DEBUG")
}

func TestDebugf(t *testing.T) {
	test(t, "DEBUG")
}

func TestFlush(t *testing.T) {
	fileLog := log.FileLog{}
	fileLog.Flush()
}

func TestGetSetting(t *testing.T) {
	setting := instance().GetSetting()
	instance().SetSetting(setting)

	if setting != instance().GetSetting() {
		t.Log(setting)
		t.Log(instance().GetSetting())

		t.Fatal("invalid Setting")
	}
}

func TestSetSetting(t *testing.T) {
	setting := instance().GetSetting()
	instance().SetSetting(setting)

	if setting != instance().GetSetting() {
		t.Log(setting)
		t.Log(instance().GetSetting())

		t.Fatal("invalid Setting")
	}
}

func TestGetLevel(t *testing.T) {
	level := instance().GetLevel()
	instance().SetLevel(level)

	if level != instance().GetLevel() {
		t.Log(level)
		t.Log(instance().GetLevel())

		t.Fatal("invalid Level")
	}
}

func TestSetLevel(t *testing.T) {
	level := instance().GetLevel()
	instance().SetLevel(level)

	if level != instance().GetLevel() {
		t.Log(level)
		t.Log(instance().GetLevel())

		t.Fatal("invalid Level")
	}
}

func TestGetOutputPath(t *testing.T) {
	outputPath := instance().GetOutputPath()
	instance().SetOutputPath(outputPath)

	if outputPath != instance().GetOutputPath() {
		t.Log(outputPath)
		t.Log(instance().GetOutputPath())

		t.Fatal("invalid OutputPath")
	}
}

func TestSetOutputPath(t *testing.T) {
	outputPath := instance().GetOutputPath()
	instance().SetOutputPath(outputPath)

	if outputPath != instance().GetOutputPath() {
		t.Log(outputPath)
		t.Log(instance().GetOutputPath())

		t.Fatal("invalid OutputPath")
	}
}

func TestGetFileNamePrefix(t *testing.T) {
	fileNamePrefix := instance().GetFileNamePrefix()
	instance().SetFileNamePrefix(fileNamePrefix)

	if fileNamePrefix != instance().GetFileNamePrefix() {
		t.Log(fileNamePrefix)
		t.Log(instance().GetFileNamePrefix())

		t.Fatal("invalid FileNamePrefix")
	}
}

func TestSetFileNamePrefix(t *testing.T) {
	fileNamePrefix := instance().GetFileNamePrefix()
	instance().SetFileNamePrefix(fileNamePrefix)

	if fileNamePrefix != instance().GetFileNamePrefix() {
		t.Log(fileNamePrefix)
		t.Log(instance().GetFileNamePrefix())

		t.Fatal("invalid FileNamePrefix")
	}
}

func TestGetPrintCallerInfo(t *testing.T) {
	printCallerInfo := instance().GetPrintCallerInfo()
	instance().SetPrintCallerInfo(printCallerInfo)

	if printCallerInfo != instance().GetPrintCallerInfo() {
		t.Log(printCallerInfo)
		t.Log(instance().GetPrintCallerInfo())

		t.Fatal("invalid PrintCallerInfo")
	}
}

func TestSetPrintCallerInfo(t *testing.T) {
	printCallerInfo := instance().GetPrintCallerInfo()
	instance().SetPrintCallerInfo(printCallerInfo)

	if printCallerInfo != instance().GetPrintCallerInfo() {
		t.Log(printCallerInfo)
		t.Log(instance().GetPrintCallerInfo())

		t.Fatal("invalid PrintCallerInfo")
	}
}

func TestGetChannelSize(t *testing.T) {
	for i := 1; i < 1000; i++ {
		instance().SetChannelSize(i)
		if instance().GetChannelSize() != i {
			t.Fatalf("invalid ChannelSize - (%d)(%d)", instance().GetChannelSize(), i)
		}
	}
}

func TestSetChannelSize(t *testing.T) {
	for i := 1; i < 1000; i++ {
		instance().SetChannelSize(i)
		if instance().GetChannelSize() != i {
			t.Fatalf("invalid ChannelSize - (%d)(%d)", instance().GetChannelSize(), i)
		}
	}
}

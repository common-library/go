// Package logger provides a file logging by level.
package log

import (
	"errors"
	"fmt"
	"sync"
)

var once sync.Once
var instance *fileLog

var loggingWaitGroup sync.WaitGroup

var mutex = new(sync.Mutex)
var channel chan logInfo = make(chan logInfo, 1024)

const (
	CRITICAL = 0 + iota
	ERROR
	WARNING
	INFO
	DEBUG
)

var logLevelInfo = map[int]string{
	CRITICAL: "CRITICAL",
	ERROR:    "ERROR",
	WARNING:  "WARNING",
	INFO:     "INFO",
	DEBUG:    "DEBUG",
}

type logInfo struct {
	level  int
	format string
	value  []interface{}
}

func singleton() *fileLog {
	once.Do(func() {
		instance = &fileLog{}
	})

	return instance
}

func logging() {
	mutex.Lock()
	defer mutex.Unlock()

	loggingWaitGroup.Add(1)
	defer loggingWaitGroup.Done()

	logInfo := <-channel

	singleton().logging(logInfo.level, logInfo.format, logInfo.value...)
}

// Initialize is initialize. If there is no outputPath, standard output.
//
// log level priority : CRITICAL < ERROR < WARNING < INFO < DEBUG
//
// ex 1)
//   log.Initialize(log.DEBUG, "./log", "")
//   // filename : ./log/20200630.log
//
// ex 2)
//    log.Initialize(log.DEBUG, "./log", "abc")
//    // filename : ./log/abc_20200630.log
//
// ex 3)
//   log.Initialize(log.DEBUG, "", "")
//   // standard output
func Initialize(level int, outputPath string, fileNamePrefix string) error {
	mutex.Lock()
	defer mutex.Unlock()

	return singleton().initialize(level, outputPath, fileNamePrefix)
}

func Finalize() error {
	Flush()

	return singleton().finalize()
}

// Critical is critical logging.
//
// ex) log.Critical("(%d) (%s)", 1, "a")
//
// output) [07:37:49] [CRITICAL] : (1) (a)
func Critical(format string, value ...interface{}) {
	channel <- logInfo{CRITICAL, format, value}

	go logging()
}

// Error is error logging.
//
// ex) log.Error("(%d) (%s)", 2, "b")
//
// output) [07:37:49] [ERROR] : (2) (b)
func Error(format string, value ...interface{}) {
	channel <- logInfo{ERROR, format, value}

	go logging()
}

// Warning is warning logging.
//
// ex) log.Warning("(%d) (%s)", 3, "c")
//
// output) [07:37:49] [WARNING] : (3) (c)
func Warning(format string, value ...interface{}) {
	channel <- logInfo{WARNING, format, value}

	go logging()
}

// Info is info logging.
//
// ex) log.Info("(%d) (%s)", 4, "d")
//
// output) [07:37:49] [INFO] : (4) (d)
func Info(format string, value ...interface{}) {
	channel <- logInfo{INFO, format, value}

	go logging()
}

// Debug is debug logging.
//
// ex) log.Debug("(%d) (%s)", 5, "e")
//
// output) [07:37:49] [DEBUG] : (5) (e)
func Debug(format string, value ...interface{}) {
	channel <- logInfo{DEBUG, format, value}

	go logging()
}

// Flush waits until all logs have been logging.
//
// ex) Flush()
func Flush() {
	for len(channel) != 0 {
	}

	loggingWaitGroup.Wait()
}

// ToIntLevel is change the log level of string type to integer type
//
// ex) level, err := ToIntLevel("DEBUG")
func ToIntLevel(level string) (int, error) {
	for key, value := range logLevelInfo {
		if value == level {
			return key, nil
		}
	}

	return -1, errors.New(fmt.Sprintf("invalid level - level : (%s)", level))
}

// GetLevel get the log level
//
// ex) level := GetLevel()
func GetLevel() int {
	return singleton().getLevel()
}

// SetLevel set the log level
//
// ex) SetLevel(log.DEBUG
func SetLevel(level int) {
	mutex.Lock()
	defer mutex.Unlock()

	singleton().setLevel(level)
}

// GetFileName get the file name
//
// ex) fileName := GetFileName()
func GetFileName() string {
	return singleton().getFileName()
}

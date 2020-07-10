// Package logger provides a file logging by level.
package log

import (
	"sync"
)

var once sync.Once
var instance *fileLog

var loggingWaitGroup sync.WaitGroup

var mutex = new(sync.Mutex)
var channel chan contentInfo = make(chan contentInfo, 1024)

const (
	CRITICAL = 0 + iota
	ERROR
	WARNING
	INFO
	DEBUG
)

var log_level_string = [...]string{
	"CRITICAL",
	"ERROR",
	"WARNING",
	"INFO",
	"DEBUG",
}

type contentInfo struct {
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

	contentInfo := <-channel

	singleton().logging(contentInfo.level, contentInfo.format, contentInfo.value...)
}

// Initialize is initialize. If there is no outputPath, standard output.
// log level priority : CRITICAL < ERROR < WARNING < INFO < DEBUG
//  ex 1) log.Initialize(log.DEBUG, "./log", "")
//   // filename : ./log/20200630.log
//  ex 2) log.Initialize(log.DEBUG, "./log", "abc")
//   // filename : ./log/abc_20200630.log
//  ex 2) log.Initialize(log.DEBUG, "", "")
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
//  ex) log.Critical("(%d) (%s)", 1, "a")
//  output) [07:37:49] [CRITICAL] : (1) (a)
func Critical(format string, value ...interface{}) {
	channel <- contentInfo{CRITICAL, format, value}

	go logging()
}

// Error is error logging.
//  ex) log.Error("(%d) (%s)", 2, "b")
//  output) [07:37:49] [ERROR] : (2) (b)
func Error(format string, value ...interface{}) {
	channel <- contentInfo{ERROR, format, value}

	go logging()
}

// Warning is warning logging.
//  ex) log.Warning("(%d) (%s)", 3, "c")
//  output) [07:37:49] [WARNING] : (3) (c)
func Warning(format string, value ...interface{}) {
	channel <- contentInfo{WARNING, format, value}

	go logging()
}

// Info is info logging.
//  ex) log.Info("(%d) (%s)", 4, "d")
//  output) [07:37:49] [INFO] : (4) (d)
func Info(format string, value ...interface{}) {
	channel <- contentInfo{INFO, format, value}

	go logging()
}

// Debug is debug logging.
//  ex) log.Debug("(%d) (%s)", 5, "e")
//  output) [07:37:49] [DEBUG] : (5) (e)
func Debug(format string, value ...interface{}) {
	channel <- contentInfo{DEBUG, format, value}

	go logging()
}

// Flush waits until all logs have been logging.
func Flush() {
	for len(channel) != 0 {
	}

	loggingWaitGroup.Wait()
}

// GetLevel get the log level
func GetLevel() int {
	return singleton().getLevel()
}

// SetLevel set the log level
func SetLevel(level int) {
	mutex.Lock()
	defer mutex.Unlock()

	singleton().setLevel(level)
}

// GetFileName get the file name
func GetFileName() string {
	return singleton().getFileName()
}

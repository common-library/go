// Package log provides logging.
package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	data_structure "github.com/heaven-chp/common-library-go/data-structure"
	"github.com/heaven-chp/common-library-go/lock"
	"github.com/heaven-chp/common-library-go/utility"
)

type Level int

const (
	LevelTrace = Level(-8)
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
	LevelFatal = Level(12)
)

var levelNames = map[slog.Leveler]string{
	slog.Level(LevelTrace): "TRACE",
	slog.Level(LevelFatal): "FATAL",
}

type outPut int

const (
	outPutStdout = outPut(0)
	outPutStderr = outPut(1)
	outPutFile   = outPut(2)
)

// Log is struct that provides log related methods.
type Log struct {
	logger *slog.Logger

	level Level

	outputLocation    outPut
	fileName          string
	fileExtensionName string
	addDate           bool
	lastDay           int
	withCallerInfo    bool

	mutexForLogger  lock.Mutex
	mutexForLogging lock.Mutex

	queueForTime    data_structure.Queue[time.Time]
	queueForLogging data_structure.Queue[func()]
}

// Trace means recording trace level logs.
//
// ex) testLog.Trace("message-01", "key-01", "value-01", "key-02", 1)
func (this *Log) Trace(message string, arguments ...any) {
	this.producer(LevelTrace, message, arguments...)
}

// Debug means recording debug level logs.
//
// ex) testLog.Debug("message-02", "key-01", "value-02", "key-02", 2)
func (this *Log) Debug(message string, arguments ...any) {
	this.producer(LevelDebug, message, arguments...)
}

// Info means recording info level logs.
//
// ex) testLog.Info("message-03", "key-01", "value-03", "key-02", 3)
func (this *Log) Info(message string, arguments ...any) {
	this.producer(LevelInfo, message, arguments...)
}

// Warn means recording warn level logs.
//
// ex) testLog.Warn("message-04", "key-01", "value-04", "key-02", 4)
func (this *Log) Warn(message string, arguments ...any) {
	this.producer(LevelWarn, message, arguments...)
}

// Error means recording error level logs.
//
// ex) testLog.Error("message-05", "key-01", "value-05", "key-02", 5)
func (this *Log) Error(message string, arguments ...any) {
	this.producer(LevelError, message, arguments...)
}

// Fatal means recording fatal level logs.
//
// ex) testLog.Fatal("message-06", "key-01", "value-06", "key-02", 6)
func (this *Log) Fatal(message string, arguments ...any) {
	this.producer(LevelFatal, message, arguments...)
}

// Flush waits to record the logs accumulated up to the time it was called.
//
// ex) testLog.Flush()
func (this *Log) Flush() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	defer wg.Wait()

	this.queueForLogging.Push(func() { wg.Done() })
	go this.cosumer()
}

// GetLevel gets the level.
//
// ex) level := testLog.GetLevel()
func (this *Log) GetLevel() Level {
	return this.level
}

// SetLevel sets the level.
//
// ex) testLog.SetLevel(log.LevelInfo)
func (this *Log) SetLevel(level Level) {
	this.queueForLogging.Push(func() {
		this.setLogger(level, this.outputLocation, this.fileName, this.fileExtensionName, this.addDate)
	})
}

// SetOutputToStdout sets the output to standard output.
//
// ex) testLog.SetOutputToStdout()
func (this *Log) SetOutputToStdout() {
	this.queueForLogging.Push(func() { this.setLogger(this.level, outPutStdout, "", "", this.addDate) })
}

// SetOutputToStderr sets the output to standard error.
//
// ex) testLog.SetOutputToStderr()
func (this *Log) SetOutputToStderr() {
	this.queueForLogging.Push(func() { this.setLogger(this.level, outPutStderr, "", "", this.addDate) })
}

// SetOutputToFile sets the output to file.
//
// ex) testLog.SetOutputToFile(fileName, fileExtensionName, true)
func (this *Log) SetOutputToFile(fileName, fileExtensionName string, addDate bool) {
	this.queueForLogging.Push(func() {
		this.setLogger(this.level, outPutFile, fileName, fileExtensionName, addDate)
	})
}

// SetWithCallerInfo also records caller information.
//
// ex) testLog.SetWithCallerInfo(true)
func (this *Log) SetWithCallerInfo(withCallerInfo bool) {
	this.queueForLogging.Push(func() { this.withCallerInfo = withCallerInfo })
}

func (this *Log) getLogger() *slog.Logger {
	this.mutexForLogger.Lock()
	defer this.mutexForLogger.Unlock()

	if this.logger == nil {
		this.logger = slog.Default()
	}

	return this.logger
}

func (this *Log) setLogger(level Level, outputLocation outPut, fileName, fileExtensionName string, addDate bool) {
	this.mutexForLogger.Lock()
	defer this.mutexForLogger.Unlock()

	this.level = level
	this.outputLocation = outputLocation
	this.fileName = fileName
	this.fileExtensionName = fileExtensionName
	this.addDate = addDate

	opts := &slog.HandlerOptions{
		Level: slog.Level(level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				if levelLabel, exists := levelNames[level]; exists {
					a.Value = slog.StringValue(levelLabel)
				} else {
					levelLabel = level.String()
				}

			} else if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(this.queueForTime.Front().String())
				this.queueForTime.Pop()
			}

			return a
		},
	}

	writer := io.Writer(os.Stdout)

	switch outputLocation {
	case outPutStdout:
		writer = os.Stdout
	case outPutStderr:
		writer = os.Stderr
	case outPutFile:
		options := os.O_WRONLY | os.O_APPEND | os.O_CREATE

		if this.addDate {
			fileName = this.fileName + "_" + time.Now().Format("20060102") + "." + this.fileExtensionName
		} else {
			fileName = this.fileName + "." + this.fileExtensionName
		}

		if file, err := os.OpenFile(fileName, options, os.FileMode(0644)); err != nil {
			writer = os.Stdout
			slog.Default().Error("file open fail", "error", err)
		} else {
			writer = file
		}
	}

	this.logger = slog.New(slog.NewJSONHandler(writer, opts))
}

func (this *Log) producer(level Level, message string, arguments ...any) {
	t := time.Now()
	callerInfo, errForCallerInfo := utility.GetCallerInfo(3)
	logger := this.getLogger()

	f := func() {
		if this.lastDay != t.Day() {
			this.lastDay = t.Day()
			this.setLogger(this.level, this.outputLocation, this.fileName, this.fileExtensionName, this.addDate)
		}

		logger = this.getLogger()
		if this.withCallerInfo && errForCallerInfo == nil {
			logger = logger.With(slog.Any("CallerInfo", callerInfo))
		} else if this.withCallerInfo && errForCallerInfo != nil {
			logger.Error("utility.GetCallerInfo fail", "error", errForCallerInfo)
		}

		logger.Log(context.Background(), slog.Level(level), message, arguments...)
	}

	this.queueForTime.Push(t)
	this.queueForLogging.Push(f)

	go this.cosumer()
}

func (this *Log) cosumer() {
	if this.mutexForLogging.TryLock() == false {
		return
	}
	defer this.mutexForLogging.Unlock()

	for this.queueForLogging.Size() != 0 {
		this.queueForLogging.Front()()
		this.queueForLogging.Pop()
	}
}

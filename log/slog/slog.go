// Package slog provides slog logging.
package slog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/common-library/go/collection"
	"github.com/common-library/go/lock"
	"github.com/common-library/go/utility"
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

	queueForTime    collection.Queue[time.Time]
	queueForLogging collection.Queue[func()]
}

// Trace means recording trace level logs.
//
// ex) testLog.Trace("message-01", "key-01", "value-01", "key-02", 1)
func (l *Log) Trace(message string, arguments ...any) {
	l.producer(LevelTrace, message, arguments...)
}

// Debug means recording debug level logs.
//
// ex) testLog.Debug("message-02", "key-01", "value-02", "key-02", 2)
func (l *Log) Debug(message string, arguments ...any) {
	l.producer(LevelDebug, message, arguments...)
}

// Info means recording info level logs.
//
// ex) testLog.Info("message-03", "key-01", "value-03", "key-02", 3)
func (l *Log) Info(message string, arguments ...any) {
	l.producer(LevelInfo, message, arguments...)
}

// Warn means recording warn level logs.
//
// ex) testLog.Warn("message-04", "key-01", "value-04", "key-02", 4)
func (l *Log) Warn(message string, arguments ...any) {
	l.producer(LevelWarn, message, arguments...)
}

// Error means recording error level logs.
//
// ex) testLog.Error("message-05", "key-01", "value-05", "key-02", 5)
func (l *Log) Error(message string, arguments ...any) {
	l.producer(LevelError, message, arguments...)
}

// Fatal means recording fatal level logs.
//
// ex) testLog.Fatal("message-06", "key-01", "value-06", "key-02", 6)
func (l *Log) Fatal(message string, arguments ...any) {
	l.producer(LevelFatal, message, arguments...)
}

// Flush waits to record the logs accumulated up to the time it was called.
//
// ex) testLog.Flush()
func (l *Log) Flush() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	defer wg.Wait()

	l.queueForLogging.Push(func() { wg.Done() })
	go l.cosumer()
}

// GetLevel gets the level.
//
// ex) level := testLog.GetLevel()
func (l *Log) GetLevel() Level {
	return l.level
}

// SetLevel sets the level.
//
// ex) testLog.SetLevel(log.LevelInfo)
func (l *Log) SetLevel(level Level) {
	l.queueForLogging.Push(func() {
		l.setLogger(level, l.outputLocation, l.fileName, l.fileExtensionName, l.addDate)
	})
}

// SetOutputToStdout sets the output to standard output.
//
// ex) testLog.SetOutputToStdout()
func (l *Log) SetOutputToStdout() {
	l.queueForLogging.Push(func() { l.setLogger(l.level, outPutStdout, "", "", l.addDate) })
}

// SetOutputToStderr sets the output to standard error.
//
// ex) testLog.SetOutputToStderr()
func (l *Log) SetOutputToStderr() {
	l.queueForLogging.Push(func() { l.setLogger(l.level, outPutStderr, "", "", l.addDate) })
}

// SetOutputToFile sets the output to file.
//
// ex) testLog.SetOutputToFile(fileName, fileExtensionName, true)
func (l *Log) SetOutputToFile(fileName, fileExtensionName string, addDate bool) {
	l.queueForLogging.Push(func() {
		l.setLogger(l.level, outPutFile, fileName, fileExtensionName, addDate)
	})
}

// SetWithCallerInfo also records caller information.
//
// ex) testLog.SetWithCallerInfo(true)
func (l *Log) SetWithCallerInfo(withCallerInfo bool) {
	l.queueForLogging.Push(func() { l.withCallerInfo = withCallerInfo })
}

func (l *Log) getLogger() *slog.Logger {
	l.mutexForLogger.Lock()
	defer l.mutexForLogger.Unlock()

	if l.logger == nil {
		l.logger = slog.Default()
	}

	return l.logger
}

func (l *Log) setLogger(level Level, outputLocation outPut, fileName, fileExtensionName string, addDate bool) {
	l.mutexForLogger.Lock()
	defer l.mutexForLogger.Unlock()

	l.level = level
	l.outputLocation = outputLocation
	l.fileName = fileName
	l.fileExtensionName = fileExtensionName
	l.addDate = addDate

	opts := &slog.HandlerOptions{
		Level: slog.Level(level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				level := a.Value.Any().(slog.Level)
				if levelLabel, exists := levelNames[level]; exists {
					a.Value = slog.StringValue(levelLabel)
				}

			case slog.TimeKey:
				a.Value = slog.StringValue(l.queueForTime.Front().String())
				l.queueForTime.Pop()
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

		if l.addDate {
			fileName = l.fileName + "_" + time.Now().Format("20060102") + "." + l.fileExtensionName
		} else {
			fileName = l.fileName + "." + l.fileExtensionName
		}

		if file, err := os.OpenFile(fileName, options, os.FileMode(0644)); err != nil {
			writer = os.Stdout
			slog.Default().Error("file open fail", "error", err)
		} else {
			writer = file
		}
	}

	l.logger = slog.New(slog.NewJSONHandler(writer, opts))
}

func (l *Log) producer(level Level, message string, arguments ...any) {
	t := time.Now()
	callerInfo, errForCallerInfo := utility.GetCallerInfo(3)
	logger := l.getLogger()

	f := func() {
		if l.lastDay != t.Day() {
			l.lastDay = t.Day()
			l.setLogger(l.level, l.outputLocation, l.fileName, l.fileExtensionName, l.addDate)
		}

		logger = l.getLogger()
		if l.withCallerInfo && errForCallerInfo == nil {
			logger = logger.With(slog.Any("CallerInfo", callerInfo))
		} else if l.withCallerInfo && errForCallerInfo != nil {
			logger.Error("utility.GetCallerInfo fail", "error", errForCallerInfo)
		}

		logger.Log(context.Background(), slog.Level(level), message, arguments...)
	}

	l.queueForTime.Push(t)
	l.queueForLogging.Push(f)

	go l.cosumer()
}

func (l *Log) cosumer() {
	if !l.mutexForLogging.TryLock() {
		return
	}
	defer l.mutexForLogging.Unlock()

	for l.queueForLogging.Size() != 0 {
		l.queueForLogging.Front()()
		l.queueForLogging.Pop()
	}
}

// Package slog provides structured logging with asynchronous output and flexible configuration.
//
// This package wraps Go's standard log/slog with additional features including
// asynchronous logging, file rotation, multiple output destinations, and caller
// information tracking.
//
// Features:
//   - Structured logging with key-value pairs
//   - Multiple log levels (Trace, Debug, Info, Warn, Error, Fatal)
//   - Asynchronous logging with queue-based buffering
//   - Output to stdout, stderr, or files
//   - Daily log file rotation
//   - Caller information tracking
//   - Thread-safe operations
//
// Example:
//
//	var logger slog.Log
//	logger.SetLevel(slog.LevelInfo)
//	logger.SetOutputToFile("app", "log", true)
//	logger.Info("Server started", "port", 8080)
//	logger.Flush()
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

// Trace logs a message at trace level with optional key-value pairs.
//
// Trace is the lowest log level, typically used for very detailed debugging
// information that is rarely needed in production.
//
// Parameters:
//   - message: The log message
//   - arguments: Optional key-value pairs (must be even number of arguments)
//
// The arguments are interpreted as alternating keys and values.
// Keys should be strings, values can be any type.
//
// Example:
//
//	var logger slog.Log
//	logger.Trace("Function entered", "function", "processData", "params", 3)
//
// Example with multiple pairs:
//
//	logger.Trace("Processing request",
//	    "requestID", "req-123",
//	    "userID", 456,
//	    "timestamp", time.Now(),
//	)
func (l *Log) Trace(message string, arguments ...any) {
	l.producer(LevelTrace, message, arguments...)
}

// Debug logs a message at debug level with optional key-value pairs.
//
// Debug level is used for detailed diagnostic information useful during
// development and troubleshooting.
//
// Parameters:
//   - message: The log message
//   - arguments: Optional key-value pairs (must be even number of arguments)
//
// Example:
//
//	var logger slog.Log
//	logger.Debug("Query executed", "sql", "SELECT * FROM users", "duration", 45)
//
// Example with error context:
//
//	logger.Debug("Retrying connection",
//	    "attempt", 2,
//	    "maxAttempts", 3,
//	    "error", lastErr.Error(),
//	)
func (l *Log) Debug(message string, arguments ...any) {
	l.producer(LevelDebug, message, arguments...)
}

// Info logs a message at info level with optional key-value pairs.
//
// Info level is used for general informational messages about application
// operation, such as startup, shutdown, and significant state changes.
//
// Parameters:
//   - message: The log message
//   - arguments: Optional key-value pairs (must be even number of arguments)
//
// Example:
//
//	var logger slog.Log
//	logger.Info("Server started", "port", 8080, "environment", "production")
//
// Example for lifecycle events:
//
//	logger.Info("Database connection established",
//	    "host", "localhost",
//	    "database", "myapp",
//	    "poolSize", 10,
//	)
func (l *Log) Info(message string, arguments ...any) {
	l.producer(LevelInfo, message, arguments...)
}

// Warn logs a message at warning level with optional key-value pairs.
//
// Warn level indicates potentially harmful situations that don't prevent
// the application from functioning but should be investigated.
//
// Parameters:
//   - message: The log message
//   - arguments: Optional key-value pairs (must be even number of arguments)
//
// Example:
//
//	var logger slog.Log
//	logger.Warn("High memory usage", "used", "85%", "threshold", "80%")
//
// Example for degraded performance:
//
//	logger.Warn("Slow query detected",
//	    "query", "SELECT * FROM large_table",
//	    "duration", "5.2s",
//	    "threshold", "1s",
//	)
func (l *Log) Warn(message string, arguments ...any) {
	l.producer(LevelWarn, message, arguments...)
}

// Error logs a message at error level with optional key-value pairs.
//
// Error level is used for errors that prevent specific operations from
// completing successfully but don't crash the application.
//
// Parameters:
//   - message: The log message
//   - arguments: Optional key-value pairs (must be even number of arguments)
//
// Example:
//
//	var logger slog.Log
//	logger.Error("Failed to save user", "userID", 123, "error", err.Error())
//
// Example with context:
//
//	logger.Error("API request failed",
//	    "endpoint", "/api/users",
//	    "method", "POST",
//	    "statusCode", 500,
//	    "error", err.Error(),
//	)
func (l *Log) Error(message string, arguments ...any) {
	l.producer(LevelError, message, arguments...)
}

// Fatal logs a message at fatal level with optional key-value pairs.
//
// Fatal level indicates critical errors that require immediate attention.
// Note: Unlike some logging libraries, this does NOT terminate the application.
//
// Parameters:
//   - message: The log message
//   - arguments: Optional key-value pairs (must be even number of arguments)
//
// Example:
//
//	var logger slog.Log
//	logger.Fatal("Database connection lost", "error", err.Error())
//
// Example for critical failures:
//
//	logger.Fatal("Configuration file corrupt",
//	    "file", "config.yaml",
//	    "error", err.Error(),
//	    "action", "manual intervention required",
//	)
func (l *Log) Fatal(message string, arguments ...any) {
	l.producer(LevelFatal, message, arguments...)
}

// Flush blocks until all queued log entries are written.
//
// This method ensures all pending asynchronous log writes are completed
// before returning. It's essential to call this before application shutdown
// to ensure no logs are lost.
//
// Behavior:
//   - Blocks the calling goroutine until all queued logs are written
//   - Processes the internal logging queue completely
//   - Thread-safe and can be called concurrently
//
// Example at application shutdown:
//
//	var logger slog.Log
//	defer logger.Flush()
//
//	logger.Info("Application starting")
//	// ... application logic ...
//	logger.Info("Application shutting down")
//	// Flush ensures shutdown message is written
//
// Example in tests:
//
//	func TestLogging(t *testing.T) {
//	    var logger slog.Log
//	    logger.SetOutputToFile("test", "log", false)
//
//	    logger.Info("Test message")
//	    logger.Flush() // Ensure log is written before test ends
//
//	    // Verify log file contents
//	}
func (l *Log) Flush() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	defer wg.Wait()

	l.queueForLogging.Push(func() { wg.Done() })
	go l.cosumer()
}

// GetLevel returns the current minimum log level.
//
// Returns:
//   - Level: Current log level threshold
//
// Only log messages at or above this level will be written.
// The returned value is one of: LevelTrace, LevelDebug, LevelInfo,
// LevelWarn, LevelError, or LevelFatal.
//
// Example:
//
//	var logger slog.Log
//	currentLevel := logger.GetLevel()
//
//	if currentLevel == slog.LevelDebug {
//	    fmt.Println("Debug logging is enabled")
//	}
//
// Example for conditional logging:
//
//	if logger.GetLevel() <= slog.LevelDebug {
//	    // Perform expensive debug data collection
//	    debugData := collectDebugInfo()
//	    logger.Debug("Debug info", "data", debugData)
//	}
func (l *Log) GetLevel() Level {
	return l.level
}

// SetLevel sets the minimum log level threshold.
//
// Parameters:
//   - level: New minimum log level (LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal)
//
// Only log messages at or above this level will be written.
// This operation is asynchronous and queued for processing.
//
// Example:
//
//	var logger slog.Log
//	logger.SetLevel(slog.LevelInfo)
//
//	logger.Debug("This will not be logged")
//	logger.Info("This will be logged")
//
// Example for dynamic level adjustment:
//
//	if os.Getenv("DEBUG") == "true" {
//	    logger.SetLevel(slog.LevelDebug)
//	} else {
//	    logger.SetLevel(slog.LevelInfo)
//	}
//
// Example for production:
//
//	logger.SetLevel(slog.LevelWarn) // Only warnings and errors in production
func (l *Log) SetLevel(level Level) {
	l.queueForLogging.Push(func() {
		l.setLogger(level, l.outputLocation, l.fileName, l.fileExtensionName, l.addDate)
	})
}

// SetOutputToStdout configures logging output to standard output.
//
// All subsequent log messages will be written to stdout (os.Stdout).
// This operation is asynchronous and queued for processing.
//
// Behavior:
//   - Switches output destination to stdout
//   - Maintains current log level
//   - Logs are written in JSON format
//
// Example:
//
//	var logger slog.Log
//	logger.SetOutputToStdout()
//	logger.Info("This goes to stdout")
//
// Example for development:
//
//	if os.Getenv("ENV") == "development" {
//	    logger.SetOutputToStdout()
//	} else {
//	    logger.SetOutputToFile("app", "log", true)
//	}
func (l *Log) SetOutputToStdout() {
	l.queueForLogging.Push(func() { l.setLogger(l.level, outPutStdout, "", "", l.addDate) })
}

// SetOutputToStderr configures logging output to standard error.
//
// All subsequent log messages will be written to stderr (os.Stderr).
// This operation is asynchronous and queued for processing.
//
// Behavior:
//   - Switches output destination to stderr
//   - Maintains current log level
//   - Logs are written in JSON format
//
// Example:
//
//	var logger slog.Log
//	logger.SetOutputToStderr()
//	logger.Error("This goes to stderr")
//
// Example for error-only logging:
//
//	errorLogger := &slog.Log{}
//	errorLogger.SetOutputToStderr()
//	errorLogger.SetLevel(slog.LevelError)
func (l *Log) SetOutputToStderr() {
	l.queueForLogging.Push(func() { l.setLogger(l.level, outPutStderr, "", "", l.addDate) })
}

// SetOutputToFile configures logging output to a file.
//
// Parameters:
//   - fileName: Base name of the log file (without extension)
//   - fileExtensionName: File extension (e.g., "log", "txt")
//   - addDate: If true, appends current date (YYYYMMDD) to filename
//
// Behavior:
//   - Creates or appends to the specified file
//   - If addDate is true, filename becomes: fileName_YYYYMMDD.extension
//   - If addDate is false, filename becomes: fileName.extension
//   - Automatic daily rotation when addDate is true
//   - Falls back to stdout if file cannot be opened
//   - Logs are written in JSON format
//
// Example without date:
//
//	var logger slog.Log
//	logger.SetOutputToFile("application", "log", false)
//	// Writes to: application.log
//
// Example with daily rotation:
//
//	logger.SetOutputToFile("app", "log", true)
//	// Writes to: app_20231218.log
//	// Automatically creates new file next day: app_20231219.log
//
// Example for different environments:
//
//	env := os.Getenv("ENV")
//	logger.SetOutputToFile(env+"-app", "log", true)
//	// production-app_20231218.log or development-app_20231218.log
func (l *Log) SetOutputToFile(fileName, fileExtensionName string, addDate bool) {
	l.queueForLogging.Push(func() {
		l.setLogger(l.level, outPutFile, fileName, fileExtensionName, addDate)
	})
}

// SetWithCallerInfo enables or disables caller information in log entries.
//
// Parameters:
//   - withCallerInfo: If true, includes file name, line number, and function name in logs
//
// When enabled, each log entry will include a "CallerInfo" field with:
//   - File: Source file name
//   - Line: Line number
//   - Function: Function name
//
// This operation is asynchronous and queued for processing.
//
// Example:
//
//	var logger slog.Log
//	logger.SetWithCallerInfo(true)
//	logger.Info("User logged in", "userID", 123)
//	// Log includes: {"CallerInfo":{"File":"main.go","Line":45,"Function":"handleLogin"},...}
//
// Example for debugging:
//
//	if os.Getenv("DEBUG") == "true" {
//	    logger.SetWithCallerInfo(true) // Enable in debug mode
//	} else {
//	    logger.SetWithCallerInfo(false) // Disable in production for performance
//	}
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

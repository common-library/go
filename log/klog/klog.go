// Package klog provides wrapper functions for Kubernetes klog logging.
//
// This package wraps k8s.io/klog/v2 with optional caller information tracking,
// providing structured logging capabilities commonly used in Kubernetes applications.
//
// Features:
//   - Info, Error, and Fatal log levels
//   - Formatted logging (Infof, Errorf, Fatalf)
//   - Line-based logging (Infoln, Errorln, Fatalln)
//   - Structured logging (InfoS, ErrorS)
//   - Optional caller information (file, line, function)
//   - Thread-safe operations
//
// Example:
//
//	klog.SetWithCallerInfo(true)
//	klog.Info("Server started")
//	klog.InfoS("Request processed", "method", "GET", "path", "/api/users")
//	klog.Flush()
package klog

import (
	"fmt"
	"sync/atomic"

	"github.com/common-library/go/utility"
	"k8s.io/klog/v2"
)

var withCallerInfo atomic.Bool

// Info logs informational messages.
//
// This function logs at info level using klog. If caller information is enabled,
// it includes file name, line number, and function name.
//
// Parameters:
//   - arguments: Values to log (similar to fmt.Print)
//
// Example:
//
//	klog.Info("Server started on port 8080")
//
// Example with multiple arguments:
//
//	klog.Info("Processing request", requestID, "from user", userID)
//
// Example with caller info enabled:
//
//	klog.SetWithCallerInfo(true)
//	klog.Info("Database connected")
//	// Output: [callerInfo:{File:"main.go" Line:42 Function:"main.init"}] Database connected
func Info(arguments ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.InfoDepth(1, append([]any{fmt.Sprintf("[callerInfo:%#v] ", callerInfo)}, arguments...)...)
		}
	} else {
		klog.InfoDepth(1, arguments...)
	}
}

// InfoS logs structured informational messages with key-value pairs.
//
// This function provides structured logging at info level. Arguments are
// interpreted as alternating keys and values.
//
// Parameters:
//   - message: The log message
//   - keysAndValues: Alternating keys (strings) and values (any type)
//
// Example:
//
//	klog.InfoS("Request completed", "method", "GET", "path", "/api/users", "duration", 45)
//
// Example with caller info:
//
//	klog.SetWithCallerInfo(true)
//	klog.InfoS("User logged in",
//	    "userID", 123,
//	    "username", "alice",
//	    "timestamp", time.Now(),
//	)
func InfoS(message string, keysAndValues ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.InfoSDepth(1, message, append([]any{"callerInfo", callerInfo}, keysAndValues...)...)
		}
	} else {
		klog.InfoSDepth(1, message, keysAndValues...)
	}
}

// Infof logs formatted informational messages.
//
// This function logs at info level using printf-style formatting.
//
// Parameters:
//   - format: Printf-style format string
//   - arguments: Values for format string placeholders
//
// Example:
//
//	klog.Infof("Server listening on port %d", 8080)
//
// Example with multiple values:
//
//	klog.Infof("Processing %d requests for user %s", count, username)
//
// Example with caller info:
//
//	klog.SetWithCallerInfo(true)
//	klog.Infof("Cache hit rate: %.2f%%", hitRate*100)
func Infof(format string, arguments ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.InfofDepth(1, "[callerInfo:%#v] "+format, append([]any{callerInfo}, arguments...)...)
		}
	} else {
		klog.InfofDepth(1, format, arguments...)
	}
}

// Infoln logs informational messages with a newline.
//
// Similar to Info but always appends a newline, like fmt.Println.
//
// Parameters:
//   - arguments: Values to log
//
// Example:
//
//	klog.Infoln("Application initialized")
//
// Example with multiple arguments:
//
//	klog.Infoln("User", userID, "logged in at", time.Now())
func Infoln(arguments ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.InfolnDepth(1, append([]any{fmt.Sprintf("[callerInfo:%#v]", callerInfo)}, arguments...)...)
		}
	} else {
		klog.InfolnDepth(1, arguments...)
	}
}

// Error logs error messages.
//
// This function logs at error level using klog. If caller information is enabled,
// it includes file name, line number, and function name.
//
// Parameters:
//   - arguments: Values to log (similar to fmt.Print)
//
// Example:
//
//	klog.Error("Failed to connect to database")
//
// Example with context:
//
//	klog.Error("Request failed:", err)
//
// Example with caller info enabled:
//
//	klog.SetWithCallerInfo(true)
//	klog.Error("Authentication failed for user", userID)
func Error(arguments ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.ErrorDepth(1, append([]any{fmt.Sprintf("[callerInfo:%#v] ", callerInfo)}, arguments...)...)
		}
	} else {
		klog.ErrorDepth(1, arguments...)
	}
}

// ErrorS logs structured error messages with key-value pairs.
//
// This function provides structured logging at error level with an explicit
// error parameter and additional context.
//
// Parameters:
//   - err: The error to log (can be nil)
//   - message: The log message
//   - keysAndValues: Alternating keys (strings) and values (any type)
//
// Example:
//
//	klog.ErrorS(err, "Database query failed", "query", sql, "duration", elapsed)
//
// Example with nil error:
//
//	klog.ErrorS(nil, "Validation failed", "field", "email", "value", userEmail)
//
// Example with caller info:
//
//	klog.SetWithCallerInfo(true)
//	klog.ErrorS(err, "Failed to save user",
//	    "userID", 123,
//	    "operation", "UPDATE",
//	    "table", "users",
//	)
func ErrorS(err error, message string, keysAndValues ...any) {
	if withCallerInfo.Load() {
		if callerInfo, errTemp := utility.GetCallerInfo(2); errTemp != nil {
			klog.ErrorS(errTemp, "")
		} else {
			klog.ErrorSDepth(1, err, message, append([]any{"callerInfo", callerInfo}, keysAndValues...)...)
		}
	} else {
		klog.ErrorSDepth(1, err, message, keysAndValues...)
	}
}

// Errorf logs formatted error messages.
//
// This function logs at error level using printf-style formatting.
//
// Parameters:
//   - format: Printf-style format string
//   - arguments: Values for format string placeholders
//
// Example:
//
//	klog.Errorf("Failed to open file: %v", err)
//
// Example with multiple values:
//
//	klog.Errorf("User %d attempted invalid action: %s", userID, action)
//
// Example with caller info:
//
//	klog.SetWithCallerInfo(true)
//	klog.Errorf("Timeout after %d attempts", retries)
func Errorf(format string, arguments ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.ErrorfDepth(1, "[callerInfo:%#v] "+format, append([]any{callerInfo}, arguments...)...)
		}
	} else {
		klog.ErrorfDepth(1, format, arguments...)
	}
}

// Errorln logs error messages with a newline.
//
// Similar to Error but always appends a newline, like fmt.Println.
//
// Parameters:
//   - arguments: Values to log
//
// Example:
//
//	klog.Errorln("Connection lost")
//
// Example with error:
//
//	klog.Errorln("Operation failed:", err)
func Errorln(arguments ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.ErrorlnDepth(1, append([]any{fmt.Sprintf("[callerInfo:%#v]", callerInfo)}, arguments...)...)
		}
	} else {
		klog.ErrorlnDepth(1, arguments...)
	}
}

// Fatal logs a message and then calls os.Exit(255).
//
// This function logs at fatal level and terminates the program.
// Deferred functions will NOT run.
//
// Parameters:
//   - arguments: Values to log before exiting
//
// Warning: This terminates the application. Use Error for recoverable errors.
//
// Example:
//
//	if configFile == nil {
//	    klog.Fatal("Configuration file required")
//	}
//
// Example with context:
//
//	klog.Fatal("Critical dependency missing:", dependency)
func Fatal(arguments ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.FatalDepth(1, append([]any{fmt.Sprintf("[callerInfo:%#v] ", callerInfo)}, arguments...)...)
		}
	} else {
		klog.FatalDepth(1, arguments...)
	}
}

// Fatalf logs a formatted message and then calls os.Exit(255).
//
// This function logs at fatal level using printf-style formatting and
// terminates the program. Deferred functions will NOT run.
//
// Parameters:
//   - format: Printf-style format string
//   - arguments: Values for format string placeholders
//
// Warning: This terminates the application. Use Errorf for recoverable errors.
//
// Example:
//
//	if port < 1024 {
//	    klog.Fatalf("Invalid port: %d", port)
//	}
//
// Example with error:
//
//	klog.Fatalf("Failed to initialize: %v", err)
func Fatalf(format string, arguments ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.FatalfDepth(1, "[callerInfo:%#v] "+format, append([]any{callerInfo}, arguments...)...)
		}
	} else {
		klog.FatalfDepth(1, format, arguments...)
	}
}

// Fatalln logs a message with a newline and then calls os.Exit(255).
//
// Similar to Fatal but always appends a newline. This function terminates
// the program. Deferred functions will NOT run.
//
// Parameters:
//   - arguments: Values to log before exiting
//
// Warning: This terminates the application.
//
// Example:
//
//	klog.Fatalln("Database connection failed")
//
// Example with error:
//
//	klog.Fatalln("Startup failed:", err)
func Fatalln(arguments ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.FatallnDepth(1, append([]any{fmt.Sprintf("[callerInfo:%#v]", callerInfo)}, arguments...)...)
		}
	} else {
		klog.FatallnDepth(1, arguments...)
	}
}

// Flush flushes all pending log I/O.
//
// This function blocks until all buffered log data has been written.
// It's essential to call this before application shutdown to ensure
// no logs are lost.
//
// Example at shutdown:
//
//	defer klog.Flush()
//
//	klog.Info("Application starting")
//	// ... application logic ...
//	klog.Info("Application shutting down")
//
// Example in signal handler:
//
//	sigChan := make(chan os.Signal, 1)
//	signal.Notify(sigChan, os.Interrupt)
//	<-sigChan
//	klog.Info("Interrupt received")
//	klog.Flush()
func Flush() {
	klog.Flush()
}

// SetWithCallerInfo enables or disables caller information in log entries.
//
// Parameters:
//   - with: If true, includes file name, line number, and function name in logs
//
// When enabled, each log entry will include caller information showing
// where the log call originated. This is useful for debugging but adds
// overhead.
//
// Example:
//
//	klog.SetWithCallerInfo(true)
//	klog.Info("Server started")
//	// Output includes: [callerInfo:{File:"main.go" Line:42 Function:"main"}] Server started
//
// Example for environment-based configuration:
//
//	if os.Getenv("DEBUG") == "true" {
//	    klog.SetWithCallerInfo(true)
//	} else {
//	    klog.SetWithCallerInfo(false)
//	}
func SetWithCallerInfo(with bool) {
	withCallerInfo.Store(with)
}

// Package klog provides klog logging.
package klog

import (
	"fmt"
	"sync/atomic"

	"github.com/heaven-chp/common-library-go/utility"
	"k8s.io/klog/v2"
)

var withCallerInfo atomic.Bool

// Info means recording info logs.
//
// ex) klog.Info("message-01")
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

// InfoS means recording info logs.
//
// ex) klog.InfoS("message-01", "key-01", 1, "key-02", "value-01")
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

// Infof means recording info logs.
//
// ex) klog.Infof("%s", "message-01")
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

// Infoln means recording info logs.
//
// ex) klog.Infoln("message-01")
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

// Error means recording error logs.
//
// ex) klog.Error("message-01")
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

// ErrorS means recording error logs.
//
// ex) klog.ErrorS(err, "message-01", "key-01", 1, "key-02", "value-01")
func ErrorS(err error, message string, keysAndValues ...any) {
	if withCallerInfo.Load() {
		if callerInfo, err := utility.GetCallerInfo(2); err != nil {
			klog.ErrorS(err, "")
		} else {
			klog.ErrorSDepth(1, err, message, append([]any{"callerInfo", callerInfo}, keysAndValues...)...)
		}
	} else {
		klog.ErrorSDepth(1, err, message, keysAndValues...)
	}
}

// Errorf means recording error logs.
//
// ex) klog.Errorf("%s", "message-01")
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

// Errorln means recording error logs.
//
// ex) klog.Errorln("message-01")
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

// Fatal means recording fatal logs.
//
// ex) klog.Fatal("message-01")
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

// Fatalf means recording fatal logs.
//
// ex) klog.Fatalf("%s", "message-01")
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

// Fatalln means recording fatal logs.
//
// ex) klog.Fatalln("message-01")
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

// Flush flushes all log.
//
// ex) klog.Flush()
func Flush() {
	klog.Flush()
}

// SetWithCallerInfo also records caller information.
//
// ex) klog.SetWithCallerInfo(true)
func SetWithCallerInfo(with bool) {
	withCallerInfo.Store(with)
}

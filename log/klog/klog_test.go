package klog_test

import (
	"errors"
	"testing"

	"github.com/common-library/go/log/klog"
)

func TestInfo(t *testing.T) {
	t.Parallel()

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.Info("message-01")

	klog.SetWithCallerInfo(true)
	klog.Info("message-02")
}

func TestInfoS(t *testing.T) {
	t.Parallel()

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.InfoS("message-01", "key-01", 1, "key-02", "value-01")

	klog.SetWithCallerInfo(true)
	klog.InfoS("message-02", "key-01", 2, "key-02", "value-02")
}

func TestInfof(t *testing.T) {
	t.Parallel()

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.Infof("%s", "message-01")

	klog.SetWithCallerInfo(true)
	klog.Infof("%s", "message-02")
}

func TestInfoln(t *testing.T) {
	t.Parallel()

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.Infoln("message-01")

	klog.SetWithCallerInfo(true)
	klog.Infoln("message-02")
}

func TestError(t *testing.T) {
	t.Parallel()

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.Error("asd")

	klog.SetWithCallerInfo(true)
	klog.Error("asd")
}

func TestErrorS(t *testing.T) {
	t.Parallel()

	defer klog.Flush()

	err := errors.New("error")

	klog.SetWithCallerInfo(false)
	klog.ErrorS(err, "message-01", "key-01", 1, "key-02", "value-01")

	klog.SetWithCallerInfo(true)
	klog.ErrorS(err, "message-02", "key-01", 2, "key-02", "value-02")
}

func TestErrorf(t *testing.T) {
	t.Parallel()

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.Errorf("%s", "message-01")

	klog.SetWithCallerInfo(true)
	klog.Errorf("%s", "message-02")
}

func TestErrorln(t *testing.T) {
	t.Parallel()

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.Errorln("message-01")

	klog.SetWithCallerInfo(true)
	klog.Errorln("message-02")
}

func TestFatal(t *testing.T) {
	return

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.Fatal("message-01")

	klog.SetWithCallerInfo(true)
	klog.Fatal("message-02")
}

func TestFatalf(t *testing.T) {
	return

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.Fatalf("%s", "message-01")

	klog.SetWithCallerInfo(true)
	klog.Fatalf("%s", "message-02")
}

func TestFatalln(t *testing.T) {
	return

	defer klog.Flush()

	klog.SetWithCallerInfo(false)
	klog.Fatalln("message-01")

	klog.SetWithCallerInfo(true)
	klog.Fatalln("message-02")
}

func TestFlush(t *testing.T) {
	klog.Flush()
}

func TestSetWithCallerInfo(t *testing.T) {
	klog.SetWithCallerInfo(true)
	klog.SetWithCallerInfo(false)
}

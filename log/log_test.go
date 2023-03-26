package log_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/heaven-chp/common-library-go/file"
	"github.com/heaven-chp/common-library-go/log"
)

func check1(log_level int) error {
	err := log.Initialize(log_level, "", "")
	if err != nil {
		return errors.New(fmt.Sprintf("Initialize fail - error : (%s)", err))
	}

	log.Critical("(%d) (%s)", 1, "a")
	log.Error("(%d) (%s)", 2, "b")
	log.Warning("(%d) (%s)", 3, "c")
	log.Info("(%d) (%s)", 4, "d")
	log.Debug("(%d) (%s)", 5, "e")

	err = log.Finalize()
	if err != nil {
		return errors.New(fmt.Sprintf("Finalize fail - error : (%s)", err))
	}

	return nil
}

func check2(outputPath string, log_level int) error {
	err := log.Initialize(log_level, outputPath, "test")
	if err != nil {
		return errors.New(fmt.Sprintf("Initialize fail - error : (%s)", err))
	}

	log.Critical("(%d) (%s)", 1, "a")
	log.Error("(%d) (%s)", 2, "b")
	log.Warning("(%d) (%s)", 3, "c")
	log.Info("(%d) (%s)", 4, "d")
	log.Debug("(%d) (%s)", 5, "e")

	err = resultCheck(log_level)
	if err != nil {
		return err
	}

	err = log.Finalize()
	if err != nil {
		return errors.New(fmt.Sprintf("Finalize fail - error : (%s)", err))
	}

	return nil
}

func resultCheck(log_level int) error {
	log.Flush()

	results := []string{
		" [CRITICAL] : (1) (a)",
		" [ERROR] : (2) (b)",
		" [WARNING] : (3) (c)",
		" [INFO] : (4) (d)",
		" [DEBUG] : (5) (e)",
	}

	fileName := log.GetFileName()

	data, err := file.Read(fileName)

	if err != nil {
		return err
	}

	for index, value := range data {
		if strings.Contains(value, results[index]) {
			continue
		}

		return errors.New(fmt.Sprintf("resultCheck fail - (%s)(%s)", value, results[index]))
	}

	return nil
}

func TestInitialize(t *testing.T) {
	err := log.Initialize(log.DEBUG, "", "")
	if err != nil {
		t.Error(err)
	}

	err = log.Finalize()
	if err != nil {
		t.Error(err)
	}
}
func TestFinalize(t *testing.T) {
	err := log.Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestCritical(t *testing.T) {
	log_level := log.CRITICAL

	err := check1(log_level)
	if err != nil {
		t.Error(err)
	}

	err = check2(t.TempDir(), log_level)
	if err != nil {
		t.Error(err)
	}
}

func TestError(t *testing.T) {
	log_level := log.ERROR

	err := check1(log_level)
	if err != nil {
		t.Error(err)
	}

	err = check2(t.TempDir(), log_level)
	if err != nil {
		t.Error(err)
	}
}

func TestWarning(t *testing.T) {
	log_level := log.WARNING

	err := check1(log_level)
	if err != nil {
		t.Error(err)
	}

	err = check2(t.TempDir(), log_level)
	if err != nil {
		t.Error(err)
	}
}

func TestInfo(t *testing.T) {
	log_level := log.INFO

	err := check1(log_level)
	if err != nil {
		t.Error(err)
	}

	err = check2(t.TempDir(), log_level)
	if err != nil {
		t.Error(err)
	}
}

func TestDebug(t *testing.T) {
	log_level := log.DEBUG

	err := check1(log_level)
	if err != nil {
		t.Error(err)
	}

	err = check2(t.TempDir(), log_level)
	if err != nil {
		t.Error(err)
	}
}

func TestFlush(t *testing.T) {
	log.Flush()
}

func TestToIntLevel(t *testing.T) {
	level, err := log.ToIntLevel("DEBUG")
	if err != nil {
		t.Error(err)
	}
	if level != log.DEBUG {
		t.Error("ToIntLevel fail")
	}

	invalidLevel := "ABC"
	level, err = log.ToIntLevel(invalidLevel)
	if level != -1 || err.Error() != "invalid level - level : ("+invalidLevel+")" {
		t.Error("ToIntLevel fail")
	}
}

func TestGetSetLevel(t *testing.T) {
	for i := log.CRITICAL; i <= log.DEBUG; i++ {
		log.SetLevel(i)
		if log.GetLevel() != i {
			t.Errorf("invalid GetLevel() - (%d)", i)
		}
	}
}

func TestGetFileName(t *testing.T) {
	outputPath := t.TempDir()
	fileNamePrefix := "test"

	if log.GetFileName() != "" {
		t.Errorf("invalid file name - fileName : (%s)", log.GetFileName())
	}

	err := log.Initialize(log.DEBUG, outputPath, fileNamePrefix)
	if err != nil {
		t.Error(err)
	}

	fileName := log.GetFileName()
	compare := outputPath + "/" + fileNamePrefix + "_" + time.Now().Format("20060102") + ".log"
	if fileName != compare {
		t.Errorf("invalid file name - fileName : (%s), compare : (%s", fileName, compare)
	}

	err = log.Finalize()
	if err != nil {
		t.Error(err)
	}
}

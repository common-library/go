package log

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/heaven-chp/common-library-go/file"
)

func check1(log_level int) error {
	err := Initialize(log_level, "", "")
	if err != nil {
		return errors.New(fmt.Sprintf("Initialize fail - error : (%s)", err))
	}

	Critical("(%d) (%s)", 1, "a")
	Error("(%d) (%s)", 2, "b")
	Warning("(%d) (%s)", 3, "c")
	Info("(%d) (%s)", 4, "d")
	Debug("(%d) (%s)", 5, "e")

	err = Finalize()
	if err != nil {
		return errors.New(fmt.Sprintf("Finalize fail - error : (%s)", err))
	}

	return nil
}

func check2(outputPath string, log_level int) error {
	err := Initialize(log_level, outputPath, "test")
	if err != nil {
		return errors.New(fmt.Sprintf("Initialize fail - error : (%s)", err))
	}

	Critical("(%d) (%s)", 1, "a")
	Error("(%d) (%s)", 2, "b")
	Warning("(%d) (%s)", 3, "c")
	Info("(%d) (%s)", 4, "d")
	Debug("(%d) (%s)", 5, "e")

	err = resultCheck(log_level)
	if err != nil {
		return err
	}

	err = Finalize()
	if err != nil {
		return errors.New(fmt.Sprintf("Finalize fail - error : (%s)", err))
	}

	return nil
}

func resultCheck(log_level int) error {
	Flush()

	results := []string{
		" [CRITICAL] : (1) (a)",
		" [ERROR] : (2) (b)",
		" [WARNING] : (3) (c)",
		" [INFO] : (4) (d)",
		" [DEBUG] : (5) (e)",
	}

	fileName := GetFileName()

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

func TestSingleton(t *testing.T) {
	loggerManager1 := singleton()
	loggerManager2 := singleton()

	if loggerManager1 != loggerManager2 {
		t.Errorf("invalid Singleton()")
	}
}

func TestInitialize(t *testing.T) {
	err := Initialize(DEBUG, "", "")
	if err != nil {
		t.Error(err)
	}

	err = Finalize()
	if err != nil {
		t.Error(err)
	}
}
func TestFinalize(t *testing.T) {
	err := Finalize()
	if err != nil {
		t.Error(err)
	}
}

func TestCritical(t *testing.T) {
	log_level := CRITICAL

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
	log_level := ERROR

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
	log_level := WARNING

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
	log_level := INFO

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
	log_level := DEBUG

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
	Flush()
}

func TestToIntLevel(t *testing.T) {
	level, err := ToIntLevel("DEBUG")
	if err != nil {
		t.Error(err)
	}
	if level != DEBUG {
		t.Error("ToIntLevel fail")
	}

	invalidLevel := "ABC"
	level, err = ToIntLevel(invalidLevel)
	if level != -1 || err.Error() != "invalid level - level : ("+invalidLevel+")" {
		t.Error("ToIntLevel fail")
	}
}

func TestGetSetLevel(t *testing.T) {
	for i := CRITICAL; i <= DEBUG; i++ {
		SetLevel(i)
		if GetLevel() != i {
			t.Errorf("invalid GetLevel() - (%d)", i)
		}
	}
}

func TestGetFileName(t *testing.T) {
	outputPath := t.TempDir()
	fileNamePrefix := "test"

	if GetFileName() != "" {
		t.Errorf("invalid file name - fileName : (%s)", GetFileName())
	}

	err := Initialize(DEBUG, outputPath, fileNamePrefix)
	if err != nil {
		t.Error(err)
	}

	fileName := GetFileName()
	compare := outputPath + "/" + fileNamePrefix + "_" + time.Now().Format("20060102") + ".log"
	if fileName != compare {
		t.Errorf("invalid file name - fileName : (%s), compare : (%s", fileName, compare)
	}

	err = Finalize()
	if err != nil {
		t.Error(err)
	}
}

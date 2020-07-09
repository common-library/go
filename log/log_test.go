package log

import (
	"errors"
	"fmt"
	"github.com/heaven-chp/common-library-go/file"
	"os"
	"strings"
	"testing"
)

func TestSingleton(t *testing.T) {
	loggerManager1 := singleton()
	loggerManager2 := singleton()

	if loggerManager1 != loggerManager2 {
		t.Errorf("invalid Singleton()")
	}
}

func TestGetLevel(t *testing.T) {
	for i := CRITICAL; i <= DEBUG; i++ {
		SetLevel(i)
		if GetLevel() != i {
			t.Errorf("invalid GetLevel() - (%d)", i)
		}
	}
}

func check1(log_level int) error {
	const outputPath = "./test"

	os.RemoveAll(outputPath)

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

	os.RemoveAll(outputPath)

	return nil
}

func check2(log_level int) error {
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

	content, err := file.GetContent(fileName)

	if err != nil {
		return err
	}

	for index, value := range content {
		if strings.Contains(value, results[index]) {
			continue
		}

		return errors.New(fmt.Sprintf("resultCheck fail - (%s)(%s)", value, results[index]))
	}

	return nil
}

func TestCritical(t *testing.T) {
	log_level := CRITICAL

	err := check1(log_level)
	if err != nil {
		t.Error(err)
	}

	err = check2(log_level)
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

	err = check2(log_level)
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

	err = check2(log_level)
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

	err = check2(log_level)
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

	err = check2(log_level)
	if err != nil {
		t.Error(err)
	}
}

package log

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type fileLog struct {
	level          int
	outputPath     string
	fileNamePrefix string

	file *os.File
}

func (fileLog *fileLog) initialize(level int, outputPath string, fileNamePrefix string) error {
	fileLog.level = level
	fileLog.outputPath = outputPath
	fileLog.fileNamePrefix = fileNamePrefix

	fileLog.finalize()

	if len(fileLog.outputPath) == 0 {
		return nil
	}

	err := os.MkdirAll(outputPath, 0777)
	if err != nil {
		os.RemoveAll(outputPath)
		return err
	}

	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	fileLog.file, err = os.OpenFile(fileLog.makeFileName(), options, os.FileMode(0644))
	return err
}

func (fileLog *fileLog) finalize() error {
	if fileLog.file != nil {
		err := fileLog.file.Close()
		fileLog.file = nil

		return err
	}

	return nil
}

func (fileLog *fileLog) isShow(level int) bool {
	if level > fileLog.level {
		return false
	}

	return true
}

func (fileLog *fileLog) makeContents(level int, format string, value ...interface{}) string {
	t := time.Now()
	contents := fmt.Sprintf(format, value...)

	contentsFinal := fmt.Sprintf("[%02d:%02d:%02d] [%s] : %s", t.Hour(), t.Minute(), t.Second(), logLevelInfo[level], contents)

	return contentsFinal
}

func (fileLog *fileLog) logging(level int, format string, value ...interface{}) {
	if fileLog.isShow(level) == false {
		return
	}

	if fileLog.file != nil && strings.Contains(fileLog.file.Name(), time.Now().Format("20060102")) == false {
		err := fileLog.initialize(fileLog.level, fileLog.outputPath, fileLog.fileNamePrefix)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	contents := fileLog.makeContents(level, format, value...)

	if len(singleton().outputPath) == 0 {
		fmt.Println(contents)
		return
	}

	_, err := fmt.Fprintln(fileLog.file, contents)
	if err != nil {
		log.Fatal(err)
	}
}

func (fileLog *fileLog) makeFileName() string {
	fileName := ""

	if len(fileLog.outputPath) != 0 {
		fileName = fileLog.outputPath + "/"
	}

	if len(fileLog.fileNamePrefix) != 0 {
		fileName += fileLog.fileNamePrefix + "_"
	}

	fileName += time.Now().Format("20060102") + ".log"

	return fileName
}

func (fileLog *fileLog) getLevel() int {
	return fileLog.level
}

func (fileLog *fileLog) setLevel(level int) {
	fileLog.level = level
}

func (fileLog *fileLog) getFileName() string {
	if fileLog.file != nil {
		return fileLog.file.Name()
	}

	return ""
}

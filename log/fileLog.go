package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

type fileLog struct {
	logLevel       int
	outputPath     string
	fileNamePrefix string

	logFile *os.File
}

func (fileLog *fileLog) initialize(logLevel int, outputPath string, fileNamePrefix string) error {
	fileLog.logLevel = logLevel
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
	fileLog.logFile, err = os.OpenFile(fileLog.getFileName(), options, os.FileMode(0644))
	if err != nil {
		return err
	}

	return nil
}

func (fileLog *fileLog) finalize() error {
	if fileLog.logFile != nil {
		err := fileLog.logFile.Close()
		fileLog.logFile = nil

		return err
	}

	return nil
}

func (fileLog *fileLog) isShow(logLevel int) bool {
	if logLevel > fileLog.logLevel {
		return false
	}

	return true
}

func (fileLog *fileLog) makeContent(logLevel int, format string, value ...interface{}) string {
	t := time.Now()
	content := fmt.Sprintf(format, value...)

	contentFinal := fmt.Sprintf("[%02d:%02d:%02d] [%s] : %s", t.Hour(), t.Minute(), t.Second(), log_level_string[logLevel], content)

	return contentFinal
}

func (fileLog *fileLog) logging(logLevel int, format string, value ...interface{}) {
	if fileLog.isShow(logLevel) == false {
		return
	}

	content := fileLog.makeContent(logLevel, format, value...)

	if len(singleton().outputPath) == 0 {
		fmt.Printf("~~~ (%s) ~~~\n", content)
		fmt.Println(content)
		return
	}

	_, err := fmt.Fprintln(fileLog.logFile, content)
	if err != nil {
		log.Fatal(err)
	}
}

func (fileLog *fileLog) getFileName() string {
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

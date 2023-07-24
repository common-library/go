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

func (this *fileLog) initialize(level int, outputPath string, fileNamePrefix string) error {
	this.finalize()

	this.level = level
	this.outputPath = outputPath
	this.fileNamePrefix = fileNamePrefix

	if len(this.outputPath) == 0 {
		return nil
	}

	err := os.MkdirAll(outputPath, 0777)
	if err != nil {
		os.RemoveAll(outputPath)
		return err
	}

	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile(this.makeFileName(), options, os.FileMode(0644))
	if err != nil {
		return err
	}
	this.file = file

	return nil
}

func (this *fileLog) finalize() error {
	this.level = 0
	this.outputPath = ""
	this.fileNamePrefix = ""

	if this.file != nil {
		err := this.file.Close()
		this.file = nil

		return err
	}

	return nil
}

func (this *fileLog) isShow(level int) bool {
	if level > this.level {
		return false
	}

	return true
}

func (this *fileLog) makeContents(level int, format string, value ...interface{}) string {
	t := time.Now()
	contents := fmt.Sprintf(format, value...)

	contentsFinal := fmt.Sprintf("[%02d:%02d:%02d] [%s] : %s", t.Hour(), t.Minute(), t.Second(), logLevelInfo[level], contents)

	return contentsFinal
}

func (this *fileLog) logging(level int, format string, value ...interface{}) {
	if this.isShow(level) == false {
		return
	}

	if this.file != nil && strings.Contains(this.file.Name(), time.Now().Format("20060102")) == false {
		err := this.initialize(this.level, this.outputPath, this.fileNamePrefix)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	contents := this.makeContents(level, format, value...)

	if this.file == nil {
		fmt.Println(contents)
		return
	}

	_, err := fmt.Fprintln(this.file, contents)
	if err != nil {
		log.Fatal(err)
	}
}

func (this *fileLog) makeFileName() string {
	fileName := ""

	if len(this.outputPath) != 0 {
		fileName = this.outputPath + "/"
	}

	if len(this.fileNamePrefix) != 0 {
		fileName += this.fileNamePrefix + "_"
	}

	fileName += time.Now().Format("20060102") + ".log"

	return fileName
}

func (this *fileLog) getLevel() int {
	return this.level
}

func (this *fileLog) setLevel(level int) {
	this.level = level
}

func (this *fileLog) getFileName() string {
	if this.file != nil {
		return this.file.Name()
	}

	return ""
}

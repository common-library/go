// Package log provides a file logging by level.
package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/heaven-chp/common-library-go/file"
	"github.com/heaven-chp/common-library-go/utility"
)

var levelInfo = map[string]int{
	"CRITICAL": 1,
	"ERROR":    2,
	"WARNING":  3,
	"INFO":     4,
	"DEBUG":    5}

type contentsInfo struct {
	t       time.Time
	setting Setting

	level      string
	isFormat   bool
	format     string
	value      []any
	callerInfo utility.CallerInfo

	isFlush bool
}

// Setting is struct that provides FileLog setting.
type Setting struct {
	Level           string
	OutputPath      string
	FileNamePrefix  string
	PrintCallerInfo bool
	ChannelSize     int
}

// FileLog is struct that provides file log related methods.
type FileLog struct {
	condition    atomic.Bool
	conditionJob atomic.Bool

	mutexForSetting sync.Mutex
	mutexForChannel sync.Mutex

	channelForLogging chan contentsInfo
	channelForFlush   chan contentsInfo

	setting Setting
}

// Initialize is initialize.
// log level priority : CRITICAL, ERROR, WARNING, INFO, DEBUG
//
// ex) err := fileLog.Initialize(Setting{...})
func (this *FileLog) Initialize(setting Setting) error {
	if err := this.Finalize(); err != nil {
		return nil
	}

	_, exists := levelInfo[setting.Level]
	if exists == false {
		return errors.New("please select one of CRITICAL, ERROR, WARNING, INFO, DEBUG")
	}

	this.SetSetting(setting)

	this.makeChannel()

	this.condition.Store(true)
	go this.job()

	return nil
}

// Finalize is finalize.
//
// ex) err := fileLog.Finalize()
func (this *FileLog) Finalize() error {
	this.condition.Store(false)

	this.Flush()

	this.closeChannel()

	for this.conditionJob.Load() {
	}

	this.SetSetting(Setting{})

	return nil
}

// Critical is logging a log of the critical level.
//
// ex) fileLog.Critical(1, 1.1, "a")
func (this *FileLog) Critical(value ...any) {
	this.logging("CRITICAL", false, "", value...)
}

// Critical is logging the critical level log by specifying the format.
//
// ex) fileLog.Criticalf("(%d) (%s)", 1, "a")
func (this *FileLog) Criticalf(format string, value ...any) {
	this.logging("CRITICAL", true, format, value...)
}

// Error is logging a log of the error level.
//
// ex) fileLog.Error(1, 1.1, "a")
func (this *FileLog) Error(value ...any) {
	this.logging("ERROR", false, "", value...)
}

// Errorf is logging the error level log by specifying the format.
//
// ex) fileLog.Errorf("(%d) (%s)", 1, "a")
func (this *FileLog) Errorf(format string, value ...any) {
	this.logging("ERROR", true, format, value...)
}

// Warning is logging a log of the warning level.
//
// ex) fileLog.Warning(1, 1.1, "a")
func (this *FileLog) Warning(value ...any) {
	this.logging("WARNING", false, "", value...)
}

// Warningf is logging the warning level log by specifying the format.
//
// ex) fileLog.Warningf("(%d) (%s)", 1, "a")
func (this *FileLog) Warningf(format string, value ...any) {
	this.logging("WARNING", true, format, value...)
}

// Info is logging a log of the info level.
//
// ex) fileLog.Info(1, 1.1, "a")
func (this *FileLog) Info(value ...any) {
	this.logging("INFO", false, "", value...)
}

// Infof is logging the info level log by specifying the format.
//
// ex) fileLog.Infof("(%d) (%s)", 1, "a")
func (this *FileLog) Infof(format string, value ...any) {
	this.logging("INFO", true, format, value...)
}

// Debug is logging a log of the debug level.
//
// ex) fileLog.Debug(1, 1.1, "a")
func (this *FileLog) Debug(value ...any) {
	this.logging("DEBUG", false, "", value...)
}

// Debugf is logging the debug level log by specifying the format.
//
// ex) fileLog.Debugf("(%d) (%s)", 1, "a")
func (this *FileLog) Debugf(format string, value ...any) {
	this.logging("DEBUG", true, format, value...)
}

// Flush is will wait for the requested log to be logging before Flush is called.
//
// ex) fileLog.Flush()
func (this *FileLog) Flush() {
	if this.conditionJob.Load() == false {
		return
	}

	this.sendChannel(contentsInfo{isFlush: true})

	<-this.channelForFlush
}

// GetSetting is get the Setting
//
// ex) setting := fileLog.GetSetting()
func (this *FileLog) GetSetting() Setting {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	return this.setting
}

// SetSetting is set the Setting
//
// ex) fileLog.SetSetting(Setting{...})
func (this *FileLog) SetSetting(setting Setting) {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	this.setting = setting
}

// GetLevel is get the level of Setting
//
// ex) level := fileLog.GetLevel()
func (this *FileLog) GetLevel() string {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	return this.setting.Level
}

// SetLevel is set the level of Setting
//
// ex) fileLog.SetLevel("DEBUG")
func (this *FileLog) SetLevel(level string) {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	this.setting.Level = level
}

// GetOutputPath is get the output path of Setting
//
// ex) outputPath := fileLog.GetOutputPath()
func (this *FileLog) GetOutputPath() string {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	return this.setting.OutputPath
}

// SetOutputPath is set the output path of Setting
//
// ex) fileLog.SetOutputPath("./")
func (this *FileLog) SetOutputPath(outputPath string) {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	this.setting.OutputPath = outputPath
}

// GetFileNamePrefix is get the file name prefix of Setting
//
// ex) fileNamePrefix := fileLog.GetFileNamePrefix()
func (this *FileLog) GetFileNamePrefix() string {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	return this.setting.FileNamePrefix
}

// SetFileNamePrefix is set the the file name prefix of Setting
//
// ex) fileLog.SetFileNamePrefix("xxx")
func (this *FileLog) SetFileNamePrefix(fileNamePrefix string) {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	this.setting.FileNamePrefix = fileNamePrefix
}

// GetPrintCallerInfo is get the print caller info of Setting
//
// ex) printCallerInfo := fileLog.GetPrintCallerInfo()
func (this *FileLog) GetPrintCallerInfo() bool {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	return this.setting.PrintCallerInfo
}

// SetPrintCallerInfo is set the print caller info of Setting
//
// ex) fileLog.SetPrintCallerInfo(true)
func (this *FileLog) SetPrintCallerInfo(printCallerInfo bool) {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	this.setting.PrintCallerInfo = printCallerInfo
}

// GetChannelSize is get the channel size of Setting
//
// ex) channelSize := fileLog.GetChannelSize()
func (this *FileLog) GetChannelSize() int {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	return this.setting.ChannelSize
}

// SetChannelSize is set the channel size of Setting
//
// ex) fileLog.SetChannelSize(1024)
func (this *FileLog) SetChannelSize(channelSize int) {
	this.mutexForSetting.Lock()
	defer this.mutexForSetting.Unlock()

	this.setting.ChannelSize = channelSize

	this.makeChannel()
}

func (this *FileLog) logging(level string, isFormat bool, format string, value ...any) {
	setting := this.GetSetting()

	if this.condition.Load() == false ||
		levelInfo[level] > levelInfo[setting.Level] {
		return
	}

	callerInfo, err := utility.GetCallerInfo(3)
	if err != nil && setting.PrintCallerInfo {
		log.Println(err)
	}

	this.sendChannel(contentsInfo{
		t:          time.Now(),
		setting:    setting,
		level:      level,
		isFormat:   isFormat,
		format:     format,
		value:      value,
		callerInfo: callerInfo,
		isFlush:    false})
}

func (this *FileLog) makeChannel() {
	this.mutexForChannel.Lock()
	defer this.mutexForChannel.Unlock()

	if this.channelForFlush == nil {
		this.channelForFlush = make(chan contentsInfo)
	}

	for len(this.channelForLogging) != 0 {
	}

	if this.channelForLogging != nil {
		close(this.channelForLogging)
	}

	this.channelForLogging = make(chan contentsInfo, this.setting.ChannelSize)
}

func (this *FileLog) sendChannel(data contentsInfo) {
	this.mutexForChannel.Lock()
	defer this.mutexForChannel.Unlock()

	if this.channelForLogging == nil {
		return
	}

	this.channelForLogging <- data
}

func (this *FileLog) closeChannel() {
	this.mutexForChannel.Lock()
	defer this.mutexForChannel.Unlock()

	if this.channelForLogging != nil {
		close(this.channelForLogging)
	}
}

func (this *FileLog) job() {
	this.conditionJob.Store(true)
	defer this.conditionJob.Store(false)

	makeFileName := func(t time.Time, setting Setting) string {
		fileName := ""

		if len(setting.OutputPath) != 0 {
			fileName = setting.OutputPath + "/"
		}

		if len(setting.FileNamePrefix) != 0 {
			fileName += setting.FileNamePrefix + "_"
		}

		fileName += t.Format("20060102") + ".log"

		return fileName
	}

	getFile := func(contentsInfo contentsInfo) (*os.File, error) {
		if err := file.CreateDirectoryAll(contentsInfo.setting.OutputPath, 0777); err != nil {
			return nil, err
		}

		fileName := makeFileName(contentsInfo.t, contentsInfo.setting)
		options := os.O_WRONLY | os.O_APPEND | os.O_CREATE

		return os.OpenFile(fileName, options, os.FileMode(0644))
	}

	getContents := func(contentsInfo contentsInfo) string {
		prefix := fmt.Sprintf("[%02d:%02d:%02d]%s[%s]",
			contentsInfo.t.Hour(),
			contentsInfo.t.Minute(),
			contentsInfo.t.Second(),
			"%s",
			contentsInfo.level)

		if contentsInfo.setting.PrintCallerInfo {
			prefix = fmt.Sprintf(prefix,
				fmt.Sprintf("[%d:%s:%s:%d]",
					contentsInfo.callerInfo.GoroutineID,
					contentsInfo.callerInfo.FileName,
					contentsInfo.callerInfo.FunctionName,
					contentsInfo.callerInfo.Line))
		} else {
			prefix = fmt.Sprintf(prefix, "")
		}

		contents := ""
		if contentsInfo.isFormat {
			contents = fmt.Sprintf("%s : %s\n", prefix, fmt.Sprintf(contentsInfo.format, contentsInfo.value...))
		} else {
			contents = fmt.Sprintf("%s : %s", prefix, fmt.Sprintln(contentsInfo.value...))
		}

		return contents
	}

	closeFile := func(file *os.File) {
		if file == nil {
			return
		}

		if err := file.Close(); err != nil {
			log.Println(err)
		}

		file = nil
	}

	var file *os.File = nil
	defer closeFile(file)

	existingSetting := Setting{}
	for this.condition.Load() {
		for contentsInfo := range this.channelForLogging {
			if contentsInfo.isFlush {
				this.channelForFlush <- contentsInfo
				continue
			}

			if existingSetting != contentsInfo.setting {
				existingSetting = contentsInfo.setting

				closeFile(file)

				newFile, err := getFile(contentsInfo)
				if err != nil {
					log.Println(err)
					continue
				}
				file = newFile
			}

			if _, err := fmt.Fprint(file, getContents(contentsInfo)); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

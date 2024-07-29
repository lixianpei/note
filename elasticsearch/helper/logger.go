package helper

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

var Logger = &logger{}

type logger struct {
	//cxt *gin.Context
	*logrus.Logger
}

func InitLogger() {
	log := &logger{}
	log.Logger = logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&MyFormatter{})
	log.AddHook(&CustomLogFile{})
	Logger = log
}

// MyFormatter ================= 自定义日志内容格式 =================
type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Local().Format("2006-01-02 15:04:05.000")
	var newLog string

	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("[%s][%s][%s][%s:%d %s][%s]\n",
			timestamp, "traceId", entry.Level, fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s][%s][%s][%s]\n", timestamp, "traceId", entry.Level, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

// CustomLogFile ================= 自定义文件 =================
type CustomLogFile struct{}

func (hook *CustomLogFile) Fire(entry *logrus.Entry) error {
	entry.Logger.Out = logFileOut()
	return nil
}
func (hook *CustomLogFile) Levels() []logrus.Level {
	return logrus.AllLevels
}

func logFileOut() (file *os.File) {
	runPath, _ := os.Getwd()
	dateDay := time.Now().Local().Format(time.DateOnly)
	logFilename := filepath.Join(runPath, "logs", dateDay+".log")
	if NewFile(logFilename) {
		var err error
		file, err = os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("初始化日志文件失败：", err.Error())
		}
	}
	return
}

// NewFile 新建一个文件
func NewFile(filename string) bool {
	//检查文件是否存在
	if _, err := os.Stat(filename); errors.Is(err, fs.ErrNotExist) {
		//如果文件不存在，则创建文件所在的目录
		if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			fmt.Println("创建目录失败：", err)
			return false
		}

		//创建文件
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("创建文件失败：", err)
			return false
		}
		defer func() {
			_ = file.Close()
		}()
		return true
	} else if err != nil {
		fmt.Println("检查文件失败：", err)
		return false
	} else {
		return true
	}
}

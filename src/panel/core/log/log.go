package core

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"goPanel/src/panel/common"
	"goPanel/src/panel/config"
	"os"
	"time"
)

func init() {
	if config.Conf.App.LogOutputType == 0 {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 显示调用位置，方法
	log.SetReportCaller(config.Conf.App.Debug)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.Level(config.Conf.App.LogLevel))
}

// 设置日志输出方式
func LogSetOutput(path string) {
	if config.Conf.App.LogOutputFlag == 1 && path != "" {
		if !common.DirOrFileByIsExists(path) && !common.CreatePath(path) {
			return
		}

		initDate := time.Now().Format("2006-01-02")
		// 创建文件
		logFilePath := path + initDate + ".log"
		file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			log.Error("日志文件打开失败或创建失败！")
		}

		fileWriter := logFileWriter{file, initDate, path}
		log.SetOutput(&fileWriter)
	}
}

type logFileWriter struct {
	file *os.File
	date string
	path string
}

func (p *logFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}
	n, e := p.file.Write(data)

	// 日志文件按照时间划分
	initDate := time.Now().Format("2006-01-02")
	if p.date != initDate {
		_ = p.file.Close()
		logFilePath := p.path + initDate + ".log"
		p.file, _ = os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	}

	return n, e
}

package common

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"runtime"
	"strings"
)

type QYLog struct {
	log *logrus.Logger
	fd  *os.File
}

func (l *QYLog) Debug(format string, args ...interface{}) {
	file, funcName := fileInfo(2)
	fields := map[string]interface{}{
		"file": file,
	}
	if _, ok := fields["method"]; !ok {
		fields["method"] = funcName
	}
	l.DebugWithFields(fields, format, args...)
}

func (l *QYLog) DebugWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	file, funcName := fileInfo(2)
	if _, ok := fields["file"]; !ok {
		fields["file"] = file
	}
	if _, ok := fields["method"]; !ok {
		fields["method"] = funcName
	}
	l.log.WithFields(fields).Debugf(format, args...)
}

//格式化信息
func (l *QYLog) Info(format string, args ...interface{}) {
	file, funcName := fileInfo(2)
	fields := map[string]interface{}{
		"file": file,
	}
	if _, ok := fields["method"]; !ok {
		fields["method"] = funcName
	}
	l.InfoWithFields(fields, format, args...)
}

//格式化信息字段
func (l *QYLog) InfoWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	file, funcName := fileInfo(2)
	if _, ok := fields["file"]; !ok {
		fields["file"] = file
	}
	if _, ok := fields["method"]; !ok {
		fields["method"] = funcName
	}
	l.log.WithFields(fields).Infof(format, args...)
}

//格式化警告的格式
func (l *QYLog) Warn(format string, args ...interface{}) {
	file, funcName := fileInfo(2)
	fields := map[string]interface{}{
		"file": file,
	}
	if _, ok := fields["method"]; !ok {
		fields["method"] = funcName
	}
	l.WarnWithFields(fields, format, args...)
}

//格式化警告字段的格式
func (l *QYLog) WarnWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	file, funcName := fileInfo(2)
	if _, ok := fields["file"]; !ok {
		fields["file"] = file
	}
	if _, ok := fields["method"]; !ok {
		fields["method"] = funcName
	}
	l.log.WithFields(fields).Warnf(format, args...)
}

//格式化错误日志格式
func (l *QYLog) Error(format string, args ...interface{}) {
	file, funcName := fileInfo(2)
	fields := map[string]interface{}{
		"file": file,
	}
	if _, ok := fields["method"]; !ok {
		fields["method"] = funcName
	}
	l.ErrorWithFields(fields, format, args...)
}

//格式化错误字段的格式
func (l *QYLog) ErrorWithFields(fields map[string]interface{}, format string, args ...interface{}) {
	file, funcName := fileInfo(2)
	if _, ok := fields["file"]; !ok {
		fields["file"] = file
	}
	if _, ok := fields["method"]; !ok {
		fields["method"] = funcName
	}
	l.log.WithFields(fields).Errorf(format, args...)
}

//关闭日志文件句柄
func (l *QYLog) Close() {
	if l.fd.Close() != nil {
		fmt.Printf("日志文件句柄 %v 无法正常关闭", l.fd)
	}
}

var Log *QYLog

//pc为调用函数的标识符,file为文件名,line为行号,ok为是否返回成功
func fileInfo(skip int) (string, string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	}

	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s:%d", file, line), funcName
}

//初始化日志
func InitLog(path, level, format string) {
	logFd, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		fmt.Printf("无法打开日志文件 %s 错误信息为: %v \n", path, err)
		//os.Exit(1)
	}

	Log = &QYLog{log: logrus.New(), fd: logFd}

	switch strings.ToLower(level) {
	case DEBUG:
		Log.log.Level = logrus.DebugLevel
	case ERROR:
		Log.log.Level = logrus.ErrorLevel
	case INFO:
		Log.log.Level = logrus.InfoLevel
	case WARN:
		Log.log.Level = logrus.WarnLevel
	default:
		Log.log.Level = logrus.InfoLevel
	}

	switch strings.ToLower(format) {
	case JSON:
		Log.log.Formatter = &logrus.JSONFormatter{}
	case TEXT:
		Log.log.Formatter = &logrus.TextFormatter{FullTimestamp: true, DisableColors: false, DisableSorting: false}
	default:
		Log.log.Formatter = &logrus.TextFormatter{FullTimestamp: true, DisableColors: false, DisableSorting: false}
	}

	Log.log.Out = logFd
}

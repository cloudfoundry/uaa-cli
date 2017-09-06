package utils

import (
	"io"
	"log"
	"fmt"
)

type Logger struct {
	infoLog *log.Logger
	// robots is like info, but should only receive machine-parsable output
	robotLog *log.Logger
	warnLog *log.Logger
	errorLog *log.Logger
}

func NewLogger(infoHandle, robotsHandle, warningHandle, errorHandle io.Writer) Logger {
	infoLog := log.New(infoHandle,"",0)
	robotLog := log.New(robotsHandle,"",0)
	warnLog := log.New(warningHandle,"",0)
	errorLog := log.New(errorHandle,"",0)
	return Logger{infoLog, robotLog, warnLog, errorLog}
}

func (l *Logger) Info(msg string) {
	l.infoLog.Println(msg)
}
func (l *Logger) Infof(format string, a ...interface{}) {
	l.infoLog.Println(fmt.Sprintf(format, a...))
}
func (l *Logger) Warn(msg string) {
	l.warnLog.Println(msg)
}
func (l *Logger) Error(msg string) {
	l.errorLog.Println(msg)
}
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.errorLog.Println(fmt.Sprintf(format, a...))
}
func (l *Logger) Robots(msg string) {
	l.robotLog.Println(msg)
}
func (l *Logger) Robotsf(format string, a ...interface{}) {
	l.robotLog.Println(fmt.Sprintf(format, a...))
}
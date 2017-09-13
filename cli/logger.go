package cli

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
)

type Logger struct {
	infoLog *log.Logger
	// robots is like info, but should only receive machine-parsable output
	robotLog *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	muted    bool
}

func NewLogger(infoHandle, robotsHandle, warningHandle, errorHandle io.Writer) Logger {
	infoLog := log.New(infoHandle, "", 0)
	robotLog := log.New(robotsHandle, "", 0)
	warnLog := log.New(warningHandle, "", 0)
	errorLog := log.New(errorHandle, "", 0)
	return Logger{infoLog, robotLog, warnLog, errorLog, false}
}

func (l *Logger) Info(msg string) {
	if l.muted {
		return
	}
	l.infoLog.Println(msg)
}
func (l *Logger) Infof(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *Logger) Warn(msg string) {
	if l.muted {
		return
	}
	yellow := color.New(color.FgYellow).SprintFunc()
	l.warnLog.Println(yellow(msg))
}

func (l *Logger) Error(msg string) {
	if l.muted {
		return
	}
	red := color.New(color.FgRed).SprintFunc()
	l.errorLog.Println(red(msg))
}
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}

func (l *Logger) Robots(msg string) {
	if l.muted {
		return
	}
	l.robotLog.Println(msg)
}
func (l *Logger) Robotsf(format string, a ...interface{}) {
	l.Robots(fmt.Sprintf(format, a...))
}
func (l *Logger) Mute() {
	l.muted = true
}
func (l *Logger) Unmute() {
	l.muted = false
}

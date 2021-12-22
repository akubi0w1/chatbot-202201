package log

import (
	"fmt"
	glog "log"
	"os"
)

var ()

type Logger struct {
	logger *glog.Logger
}

func New() *Logger {
	return &Logger{
		logger: glog.New(os.Stdout, "", glog.Ldate|glog.Ltime),
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	_format := fmt.Sprintf("[INFO] %s", format)
	l.logger.Printf(_format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	_format := fmt.Sprintf("[WARN] %s", format)
	l.logger.Printf(_format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	_format := fmt.Sprintf("[ERROR] %s", format)
	l.logger.Printf(_format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	_format := fmt.Sprintf("[DEBUG] %s", format)
	l.logger.Printf(_format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	_format := fmt.Sprintf("[Fatal] %s", format)
	l.logger.Fatalf(_format, args...)
}

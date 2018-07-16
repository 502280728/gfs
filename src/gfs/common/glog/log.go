// log
package glog

import (
	"os"
	"strings"

	"github.com/op/go-logging"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})
}

var LogOn bool = false
var File string
var Level string
var Format string

func Init() {
	console := logging.NewLogBackend(os.Stdout, "", 0)
	consoleLogger := logging.AddModuleLevel(console)

	fs, _ := os.OpenFile(File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend)
	fsLogger := logging.AddModuleLevel(logging.NewLogBackend(fs, "", 0))

	logging.SetBackend(fsLogger, consoleLogger)
	logging.SetFormatter(logging.MustStringFormatter(Format))
	switch strings.ToLower(Level) {
	case "debug":
		logging.SetLevel(logging.DEBUG, "")
	case "info":
		logging.SetLevel(logging.INFO, "")
	case "notice":
		logging.SetLevel(logging.NOTICE, "")
	case "warning":
		logging.SetLevel(logging.WARNING, "")
	case "error":
		logging.SetLevel(logging.ERROR, "")
	case "critical", "fatal":
		logging.SetLevel(logging.CRITICAL, "")
	default:
		logging.SetLevel(logging.INFO, "")
	}
}

type LoggingLogger struct {
	logger *logging.Logger
}

func GetLogger(module string) Logger {
	return &LoggingLogger{logging.MustGetLogger(module)}
}

func (log *LoggingLogger) Debug(args ...interface{}) {
	if LogOn {
		log.logger.Debug(args...)
	}
}
func (log *LoggingLogger) Debugf(format string, args ...interface{}) {
	if LogOn {
		log.logger.Debugf(format, args...)
	}
}
func (log *LoggingLogger) Info(args ...interface{}) {
	if LogOn {
		log.logger.Info(args...)
	}
}
func (log *LoggingLogger) Infof(format string, args ...interface{}) {
	if LogOn {
		log.logger.Infof(format, args...)
	}
}
func (log *LoggingLogger) Notice(args ...interface{}) {
	if LogOn {
		log.logger.Info(args...)
	}
}
func (log *LoggingLogger) Noticef(format string, args ...interface{}) {
	if LogOn {
		log.logger.Infof(format, args...)
	}
}
func (log *LoggingLogger) Warning(args ...interface{}) {
	if LogOn {
		log.logger.Info(args...)
	}
}
func (log *LoggingLogger) Warningf(format string, args ...interface{}) {
	if LogOn {
		log.logger.Infof(format, args...)
	}
}
func (log *LoggingLogger) Error(args ...interface{}) {
	if LogOn {
		log.logger.Info(args...)
	}
}
func (log *LoggingLogger) Errorf(format string, args ...interface{}) {
	if LogOn {
		log.logger.Infof(format, args...)
	}
}
func (log *LoggingLogger) Critical(args ...interface{}) {
	if LogOn {
		log.logger.Info(args...)
	}
}
func (log *LoggingLogger) Criticalf(format string, args ...interface{}) {
	if LogOn {
		log.logger.Infof(format, args...)
	}
}

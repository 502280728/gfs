// logging
package logging

import (
	"os"

	"github.com/op/go-logging"
)

func Init(file, level string) {
	format := logging.MustStringFormatter(`%{time:2006-01-02 15:04:05} %{longfile} %{level} %{message}`)
	console := logging.NewLogBackend(os.Stdout, "", 0)
	fs, _ := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModeAppend)
	backend1 := logging.NewLogBackend(fs, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)

	logging.SetBackend(logging.AddModuleLevel(console), backend1Leveled)
	logging.SetFormatter(format)
	logging.SetLevel(logging.DEBUG, "")
}
func GetLogger(module string) *logging.Logger {
	return logging.MustGetLogger(module)
}

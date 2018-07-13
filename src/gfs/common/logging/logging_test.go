// logging_test
package logging

import (
	"testing"
)

func Test1(t *testing.T) {
	Init("d:/temp/log.log", "")
	logger1 := GetLogger("com")
	logger2 := GetLogger("com2")
	logger1.Info("com")
	logger2.Info("com2")
}

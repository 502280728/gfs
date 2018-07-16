// utils_test.go
package gutils

import (
	"fmt"
	"testing"
)

func TestCeil(t *testing.T) {
	var a float64 = 10
	var b float64 = 3
	fmt.Print(Ceil64(a / b))
}

func TestCreateDirIfNotExists(t *testing.T) {
	CreateDirIfNotExists("d:\\temp\\gfsnode\\a\\b")

}

func TestGetIP(t *testing.T) {
	fmt.Println(GetLocalIP())
}

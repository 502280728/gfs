package utils

import (
	"fmt"
	"math"
	"os"

	"github.com/satori/uuid"
)

func UUID() (string, error) {
	if tmp, err := uuid.NewV4(); err == nil {
		return fmt.Sprintf("%s", tmp), nil
	} else {
		return "", err
	}
}

func Ceil(x float64) int {
	return int(math.Ceil(x))
}

func Floor(x float64) int {
	return int(math.Floor(x))
}

func Ceil64(x float64) int64 {
	return int64(math.Ceil(x))
}

func Floor64(x float64) int64 {
	return int64(math.Floor(x))
}

//创建一个文件夹，如果该文件夹不存在的话;如果该文件夹父文件夹也不存在的话，创建之
//如果存在，报错
func CreateDirIfNotExists(file string) error {
	if _, err1 := os.Stat(file); err1 != nil {
		return os.MkdirAll(file, os.ModePerm)
	} else {
		return fmt.Errorf("文件夹%s已经存在", file)
	}
}

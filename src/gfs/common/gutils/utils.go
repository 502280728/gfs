package gutils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"reflect"
	"time"

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

//传对象
func EncodeToBytes(obj interface{}) []byte {
	return EncodeToByteBuffer(obj).Bytes()
}

//传对象,使用GOB编码
func EncodeToByteBuffer(obj interface{}) *bytes.Buffer {
	var res bytes.Buffer
	EncodeToWriter(obj, &res)
	return &res
}

func EncodeToWriter(obj interface{}, writer io.Writer) {
	enc := gob.NewEncoder(writer)
	enc.Encode(obj)
}

//obj 必须是地址
func DecodeFromBytes(obj interface{}, bb []byte) {
	var buf bytes.Buffer
	buf.Write(bb)
	DecodeFromByteBuffer(obj, &buf)
}

// obj必须是地址
func DecodeFromByteBuffer(obj interface{}, bb *bytes.Buffer) {
	DecodeFromReader(obj, bb)
}

// obj必须是地址
func DecodeFromReader(obj interface{}, reader io.Reader) {
	dec := gob.NewDecoder(reader)
	if value, ok := obj.(reflect.Value); ok {
		dec.DecodeValue(value)
	} else {
		dec.Decode(obj)
	}
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "localhost"
	}
	res := ""
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				res = ipnet.IP.String()
			}
		}
	}
	return res

}

const (
	TIME_BEGIN string = "2006-01-02 15:04:05"
)

func GetNowString(pattern string) string {
	return time.Now().Format(pattern)
}
func GetNowStringSimple() string {
	return GetNowString("20060102150405")
}

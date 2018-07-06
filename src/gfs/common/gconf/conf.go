// conf
//这里的conf仅负责读取配置或者添加配置
//由于这里的conf是个map，无法具体话配置，那么
//在masternode、datanode、sdk中的具体的配置由各个使用端自身负责
package gconf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//系统有会有三个需要配置的地方：sdk、masternode、datanode

//由于sdk、masternode和datanode会处于不同的程序中，那么可以使用一个类型代表
//这三种配置。这里使用map作为总的配置
type Configuration map[string]string

const (
	GFS_BLOCK_SIZE = "gfs.block.size"
	GFS_MASTER     = "gfs.master"
)

var (
	NO_CONF_ERROR_TEMPLATE string = "NO CONF %s FOUND!!!"
	NOT_SUPPORT_ERROR      error  = errors.New("NOT SUPPORT ERROR")
	WRONG_CONF_FILE_ERROR  error  = errors.New("WRONG CONF FILE")
	NOT_VALIDATED          error  = errors.New("CANT NOT PASS VALIDATOR")
	WILL_SUPPORT_IN_FUTURE error  = errors.New("NOT SUPPORT NOW.WILL SUPPORT IN THE FUTURE")
)

var Conf Configuration = make(map[string]string) //该变量表示配置的实例

//从yaml中读取配置
func (conf Configuration) LoadYaml(path string) {
	panic(WILL_SUPPORT_IN_FUTURE)
}

//从properties文件中读取配置
// a=b #a
// #注释
// c=1
func (conf Configuration) LoadProperties(confFile string) {
	if f, err := os.Open(confFile); err == nil {
		buf := bufio.NewReader(f)
		for {
			b, _, err := buf.ReadLine()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			tmp := strings.TrimSpace(string(b))
			if strings.HasPrefix(tmp, "#") {
				continue
			}
			index := strings.Index(tmp, "=")
			if index == 0 {
				panic(WRONG_CONF_FILE_ERROR)
			}
			first := strings.TrimSpace(tmp[:index])
			second := strings.TrimSpace(tmp[index+1:])

			sIndex := strings.Index(second, "\t#")
			if sIndex > -1 {
				second = second[0:sIndex]
			}
			sIndex = strings.Index(second, " #")
			if sIndex > -1 {
				second = second[0:sIndex]
			}
			conf[first] = strings.TrimSpace(second)
		}
	} else {
		panic(err.Error())
	}
}

//添加额外的配置，会覆盖已有的配置
func (conf Configuration) AddExtra(key, value string) {
	conf[key] = value
}

//添加额外的配置，会覆盖已有的配置
func (conf Configuration) AddExtras(extra map[string]string) {
	for k, v := range extra {
		conf[k] = v
	}
}

//获取一个值，将该值存入res表示的指针,如果不存在那就放入默认值, 目前仅支持int和string类型
func (conf Configuration) GetOrDefault(key string, res interface{}, def interface{}) error {
	if v, found := conf[key]; found {
		return put(res, v)
	} else {
		value := reflect.ValueOf(res)
		ele := value.Elem()
		if !ele.IsValid() || !ele.CanSet() || value.Kind() != reflect.Ptr {
			return NOT_SUPPORT_ERROR
		}
		ele.Set(reflect.ValueOf(def))
		return nil
	}
}

//获取一个值，将该值存入res表示的指针,目前仅支持int和string类型
func (conf Configuration) Get(key string, res interface{}) error {
	if v, found := conf[key]; found {
		return put(res, v)
	} else {
		return fmt.Errorf(NO_CONF_ERROR_TEMPLATE, key)
	}
}

func put(res interface{}, v string) error {
	value := reflect.ValueOf(res)
	ele := value.Elem()
	if !ele.IsValid() || !ele.CanSet() || value.Kind() != reflect.Ptr {
		return NOT_SUPPORT_ERROR
	}
	switch ele.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if tmp, err := strconv.ParseInt(v, 10, 64); err == nil {
			ele.SetInt(tmp)
			return nil
		} else {
			return err
		}
	case reflect.String:
		ele.SetString(v)
		return nil
	default:
		return NOT_SUPPORT_ERROR
	}
}

//获取一个值
func (conf Configuration) GetString(key string) (string, error) {
	if v, found := conf[key]; found {
		return v, nil
	} else {
		return "", fmt.Errorf(NO_CONF_ERROR_TEMPLATE, key)
	}
}

//将前缀为prefixKey的配置项存入一个struct
func (conf Configuration) GetStruct(prefixKey string, res interface{}) error {
	panic(WILL_SUPPORT_IN_FUTURE)
}

func (conf Configuration) GetStringValidated(key string, cv ConfValidator) (string, error) {
	if v, err := conf.GetString(key); err == nil {
		if err1 := cv(v); err1 == nil {
			return v, nil
		} else {
			return "", err1
		}
	} else {
		return "", err
	}
}

func (conf Configuration) GetURL(key string) (string, error) {
	return conf.GetStringValidated(key, URL_VALIDATOR)
}

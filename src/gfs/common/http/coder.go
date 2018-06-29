// Coder
package http

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
)

type Coder interface {
	Decode(obj interface{}, reader io.Reader) //仅接收指针或者指针的reflect.Value
	Encode(obj interface{}) *bytes.Buffer
	IsReflectValueSupported() bool //该方法没有任何作用，仅是提醒上面的Decode和Encode，特别是Decode必须要能够支持reflect.Value类型的参数
	DecodeRequest(obj interface{}, req *http.Request)
}

// 目前仅支持gob

var (
	GOB_CODER  *GobCoder  = &GobCoder{}
	JSON_CODER *JsonCoder = &JsonCoder{}
)

type GobCoder struct{}

func (dc *GobCoder) Decode(obj interface{}, reader io.Reader) {
	dec := gob.NewDecoder(reader)
	if value, ok := obj.(reflect.Value); ok {
		dec.DecodeValue(value)
	} else {
		dec.Decode(obj)
	}
}

func (dc *GobCoder) DecodeRequest(obj interface{}, req *http.Request) {
	dc.Decode(obj, req.Body)
	req.ParseForm()
	var val reflect.Value
	if value, ok := obj.(reflect.Value); ok {
		val = value
	} else {
		val = reflect.ValueOf(obj).Elem()
	}

	var field reflect.Value
	for k, v := range req.Form {
		field = val.FieldByName(k)
		switch field.Kind() {
		case reflect.Int:
			tmp, _ := strconv.Atoi(v[0])
			field.SetInt(int64(tmp))
		case reflect.String:
			field.SetString(v[0])
		case reflect.Slice:
		default:
		}
	}
}

func (dc *GobCoder) Encode(obj interface{}) *bytes.Buffer {
	var writer bytes.Buffer
	enc := gob.NewEncoder(&writer)
	enc.Encode(obj)
	return &writer
}

func (dc *GobCoder) IsReflectValueSupported() bool {
	return true
}

type JsonCoder struct{}

func (jc *JsonCoder) Decode(obj interface{}, reader io.Reader) {
	var bb bytes.Buffer
	bb.ReadFrom(reader)
	if value, ok := obj.(reflect.Value); ok {
		json.Unmarshal(bb.Bytes(), value.Interface())
	} else {
		json.Unmarshal(bb.Bytes(), obj)
	}
}
func (jc *JsonCoder) Encode(obj interface{}) *bytes.Buffer {
	bb, _ := json.Marshal(obj)
	var res bytes.Buffer
	res.Write(bb)
	return &res
}
func (jc *JsonCoder) IsReflectValueSupported() bool {
	return true
}

func (jc *JsonCoder) DecodeRequest(obj interface{}, req *http.Request) {
	jc.Decode(obj, req.Body)
	req.ParseForm()
	if len(req.Form) <= 0 {
		return
	}
	var val reflect.Value
	var ok bool
	if val, ok = obj.(reflect.Value); !ok {
		val = reflect.ValueOf(obj).Elem()
	}
	for k, v := range req.Form {
		if field := val.FieldByName(k); !checkIfZeroValue(field) {
			switch field.Kind() {
			case reflect.Int:
				tmp, _ := strconv.Atoi(v[0])
				field.SetInt(int64(tmp))
			case reflect.String:
				field.SetString(v[0])
			case reflect.Slice:
			default:
			}
		}
	}

}

//判断一个Value类型对象是否是一个zero对象
func checkIfZeroValue(val reflect.Value) bool {
	return val == reflect.Value{}
}

// Coder
package http

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
)

type Coder interface {
	Decode(obj interface{}, reader io.Reader)
	Encode(obj interface{}) *bytes.Buffer
	IsReflectValueSupported() bool //该方法没有任何作用，仅是提醒上面的Decode和Encode，特别是Decode必须要能够支持reflect.Value类型的参数
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
	val := reflect.ValueOf(obj)
	for k, v := range req.Form {
		val.FieldByName(k).Set(reflect.ValueOf(v))
	}
}

package http

import (
	"fmt"
	"gfs/common"
	"gfs/common/http/session"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

//第一代router，将post的body反序列化一个对象，它的handler仅能接受两个参数，第一个是该对象，第二个是session
//目前仅支持 /a/b  func(a A,s session)的写法
type GFSRouter struct {
	config map[string]interface{}
	sm     *session.Manager
}

func (r *GFSRouter) Config(path string, handler interface{}) {
	if r.config == nil {
		r.config = make(map[string]interface{})
	}
	if _, found := r.config[path]; found {
		panic("配置router是path重复了")
	} else {
		r.config[path] = handler
	}

}

func (r *GFSRouter) Handle(w http.ResponseWriter, req *http.Request) {
	//	sess, _ := r.sm.SessionStart(w, req)
	uri, _ := url.Parse(req.RequestURI)
	path := uri.Path
	handlderfunc := r.getHandler(path)
	if handlderfunc != nil {
		result := handle(handlderfunc, req.Body)
		buf := common.EncodeToByteBuffer(result)
		w.Write(buf.Bytes())
	}
}

//最最核心的方法，该方法通过反射运行 handlderfunc
func handle(handlderfunc interface{}, body io.Reader) interface{} {
	tt := reflect.TypeOf(handlderfunc)
	rt := tt.In(0)
	aaa := reflect.New(rt)
	common.DecodeFromReader(aaa, body)
	fu := reflect.ValueOf(handlderfunc)
	params := []reflect.Value{aaa.Elem()}
	result := fu.Call(params)
	switch result[0].Kind() {
	case reflect.Interface, reflect.Ptr:
		return result[0].Elem().Interface()
	default:
		return result[0].Interface()
	}
}

func (r *GFSRouter) getHandler(path string) interface{} {
	for k, v := range r.config {
		if strings.HasPrefix(path, k) {
			return v
		}
	}
	return nil
}

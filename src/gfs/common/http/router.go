package http

import (
	"gfs/common/http/session"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

//第一代router，将post的body反序列化一个对象，它的handler仅能接受两个参数，第一个是该对象，第二个是session
//目前仅支持 /a/b  func(a A,s session)的写法
//
type GFSRouter struct {
	config       map[string]interface{}
	sm           *session.Manager
	DefaultCoder Coder // 默认的请求和结果的编码器,如果在Congif时没有指定编码器，那么就使用该
	coderCofig   map[string]Coder
}

func (r *GFSRouter) ConfigWithCoder(path string, handler interface{}, coder Coder) {
	if r.config == nil {
		r.config = make(map[string]interface{})
	}
	if _, found := r.config[path]; found {
		panic("配置router是path重复了")
	} else {
		r.config[path] = handler
	}
	if r.coderCofig == nil {
		r.coderCofig = make(map[string]Coder)
	}
	r.coderCofig[path] = coder
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
	if r.coderCofig == nil {
		r.coderCofig = make(map[string]Coder)
	}
	if r.DefaultCoder == nil {
		panic("一定要配置默认编码器")
	}
	r.coderCofig[path] = r.DefaultCoder
}

func (r *GFSRouter) Handle(w http.ResponseWriter, req *http.Request) {
	//	sess, _ := r.sm.SessionStart(w, req)
	uri, _ := url.Parse(req.RequestURI)
	path := uri.Path
	handlderfunc := r.getHandler(path)
	coder := r.getCoder(path)
	if handlderfunc != nil {
		result := handle(handlderfunc, req.Body, coder)
		buf := coder.Encode(result)
		w.Write(buf.Bytes())
	}
}

//最最核心的方法，该方法通过反射运行 handlderfunc
func handle(handlderfunc interface{}, body io.Reader, coder Coder) interface{} {
	tt := reflect.TypeOf(handlderfunc)
	rt := tt.In(0)
	aaa := reflect.New(rt)
	coder.Decode(aaa, body)
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
		if pathAdaptor(path, k) {
			return v
		}
	}
	return nil
}

func (r *GFSRouter) getCoder(path string) Coder {
	for k, c := range r.coderCofig {
		if pathAdaptor(path, k) {
			return c
		}
	}
	return nil
}

//请求路径与实际配置路径的适配器。目前仅支持前缀匹配
func pathAdaptor(path string, pattern string) bool {
	return strings.HasPrefix(path, pattern)
}

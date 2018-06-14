// main
package main

import (
	"fmt"
	"gfs/common"
	http1 "gfs/common/http"
	"net/http"
	"reflect"
	"strconv"
)

type A struct {
	Name string
	Age  int
	MM   []string
}

type B struct {
	CName string
	CAge  int
}

func main() {
	var router = &http1.GFSRouter{DefaultCoder: http1.JSON_CODER}
	router.Config("/work/1/2", work)

	var hh = common.Handler(func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		val := reflect.New(reflect.TypeOf(A{})).Elem()
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
				kk := reflect.New(field.Type()).Elem()
				value := make([]reflect.Value, len(v))
				for index, tmpv := range v {
					value[index] = reflect.ValueOf(tmpv)
				}
				kk = reflect.Append(kk, value...)
				field.Set(kk)
			default:
			}
		}
		fmt.Println(val.Interface())
		bb := http1.JSON_CODER.Encode(val.Interface())
		w.Write(bb.Bytes())
	})

	http.ListenAndServe(":8080", hh)
}

func work(a A) B {
	return B{a.Name, a.Age}
}

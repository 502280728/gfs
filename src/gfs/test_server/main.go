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

		for k, v := range req.Form {
			switch val.FieldByName(k).Kind() {
			case reflect.Int:
				tmp, _ := strconv.Atoi(v[0])
				val.FieldByName(k).SetInt(int64(tmp))
			case reflect.String:
				val.FieldByName(k).SetString(v[0])

			}
		}
		fmt.Println(val.Interface())
	})

	http.ListenAndServe(":8080", hh)
}

func work(a A) B {
	return B{a.Name, a.Age}
}

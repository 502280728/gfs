// main
package main

import (
	"gfs/common"
	http1 "gfs/common/http"
	"net/http"
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
	var router = &http1.GFSRouter{}
	router.Config("/work/1/2", work)

	var hh = common.Handler(func(w http.ResponseWriter, req *http.Request) {
		router.Handle(w, req)
	})

	http.ListenAndServe(":8080", hh)
}

func work(a A) B {
	return B{a.Name, a.Age}
}

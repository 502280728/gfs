package common

import (
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request)

func (handler Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler(w, req)
}

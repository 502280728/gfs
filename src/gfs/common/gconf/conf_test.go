// conf_test
package gconf

import (
	"testing"
)

func TestGet(t *testing.T) {
	var a int64
	var b string
	var c string
	Conf.AddExtra("a", "123")
	Conf.AddExtra("abc", "abc")
	Conf.AddExtra("master", "http://localhost:8080")
	Conf.AddExtra("master1", "localhost:8080")
	Conf.Get("a", &a)
	Conf.Get("abc", &b)

	if a != int64(123) {
		t.FailNow()
	}
	if b != "abc" {
		t.FailNow()
	}
	if _, err1 := Conf.GetURL("master"); err1 != nil {
		t.FailNow()
	}
	if _, err2 := Conf.GetURL("master1"); err2 == nil {
		t.FailNow()
	}
	if err3 := Conf.GetOrDefault("c", &c, "this"); err3 != nil || c != "this" {
		t.FailNow()
	}

}

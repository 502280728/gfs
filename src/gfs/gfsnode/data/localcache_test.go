// localcache_test.go
package data

import (
	"testing"
	"time"
)

func TestNoticeMaster(t *testing.T) {
	LocalStore.Config(map[string]string{"master": "http://localhost"})
	LocalStore.Set("/abc", &Data{File: "asdf", Block: 1, Path: "asdf", reported: false, modifyTime: time.Now()})
	LocalStore.noticeMaster()
}

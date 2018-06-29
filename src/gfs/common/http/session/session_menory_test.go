// session_menory_test.go
package session

import (
	"container/list"
	"fmt"
	"testing"
	"time"
)

var sp SessionProvider = &MapSessionProvider{
	content:     make(map[string]*list.Element),
	list:        list.New(),
	maxLifeTime: 10000,
}

func init() {
	fmt.Println("aaaaa")
}

func TestSession(t *testing.T) {
	var sess Session = sp.GetSession("a")
	sess.Set("bb", "cc")
	fmt.Println(sess.Get("bb"))
}

// session_menory_test.go
package session

import (
	"fmt"
	"testing"
)

var sp SessionProvider = MapSessionProvider{content: make(map[string]MapSession)}

func init() {
	fmt.Println("aaaaa")
}

func TestSession(t *testing.T) {
	var sess Session = sp.GetSession("a")
	sess.Set("bb", "cc")
	fmt.Println(sess.Get("bb"))
}

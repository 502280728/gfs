// session_init.go
package session

import (
	"container/list"
)

const (
	MAP_SESSION_PROVIDER = "mm"
)

func init() {
	var sp SessionProvider = &MapSessionProvider{
		content:     make(map[string]*list.Element),
		maxLifeTime: 1800,
		list:        list.New(),
	}
	RegisterProvider(MAP_SESSION_PROVIDER, sp)
}

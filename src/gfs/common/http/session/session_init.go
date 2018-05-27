// session_init.go
package session

const (
	MAP_SESSION_PROVIDER = "mm"
)

func init() {
	var sp SessionProvider = MapSessionProvider{content: make(map[string]MapSession)}
	RegisterProvider(MAP_SESSION_PROVIDER, sp)
}

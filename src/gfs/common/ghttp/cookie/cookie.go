// cookie
package cookie

import (
	"net/http"
)

type CookieStore interface {
	StoreCookie(w *http.Response)
	WriteCookie(req *http.Request)
}

type GFSCookieStore struct {
	Cookies []*http.Cookie
}

func (cs *GFSCookieStore) StoreCookie(w *http.Response) {
	cs.Cookies = append(cs.Cookies, w.Cookies()...)
}

func (cs *GFSCookieStore) WriteCookie(req *http.Request) {
	for _, cookie := range cs.Cookies {
		req.AddCookie(cookie)
	}
}

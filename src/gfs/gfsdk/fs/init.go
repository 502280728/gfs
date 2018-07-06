// init
package fs

import (
	http1 "gfs/common/ghttp"
	"gfs/common/ghttp/cookie"
	"net/http"
)

var cs = &cookie.GFSCookieStore{Cookies: []*http.Cookie{}}

var GFSREQ = http1.GFSRequest{
	CS: cs,
}

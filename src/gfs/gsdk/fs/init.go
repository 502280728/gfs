// init
package fs

import (
	http1 "gfs/common/http"
	"gfs/common/http/cookie"
	"net/http"
)

var cs = &cookie.GFSCookieStore{Cookies: []*http.Cookie{}}

var GFSREQ = http1.GFSRequest{
	CS: cs,
}

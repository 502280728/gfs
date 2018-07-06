// init
package gfsdk

import (
	http1 "gfs/common/ghttp"
	"gfs/common/ghttp/cookie"
	"net/http"

	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("gfs/gfsdk")

var cs = &cookie.GFSCookieStore{Cookies: []*http.Cookie{}}

var GFSREQ = http1.GFSRequest{
	CS: cs,
}

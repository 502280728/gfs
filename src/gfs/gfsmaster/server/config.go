package server

import (
	"gfs/common"
	"gfs/common/http/session"

	logging "github.com/op/go-logging"
)

var sm *session.Manager
var logger = logging.MustGetLogger("gfs/gfsmaster/server")

var (
	InvalidSession = common.MessageInFS{
		Success: false,
		Msg:     "InvalidSession",
	}

	SuccessLogin = common.MessageInFS{
		Success: true,
	}
)

func failMessage(err error) common.MessageInFS {
	return common.MessageInFS{
		Success: false,
		Msg:     err.Error(),
	}
}

const SessionUser = "session_user"

func init() {
	sm = session.NewProvider(session.MAP_SESSION_PROVIDER, "")
	//sm.SessionGC()
}

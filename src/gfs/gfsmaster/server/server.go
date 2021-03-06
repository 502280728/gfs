package server

//该package主要负责接收URL请求，并分发到相应的处理器，可以认为是MVC中的C
import (
	"gfs/common"
	"gfs/gfsmaster/fs"
	"gfs/gfsmaster/fs/user"
	"gfs/gfsmaster/node"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var conf string
	var cmd = &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			initFileSystem()
			svr := Server{}
			svr.start()
		},
	}
	cmd.Flags().StringVarP(&conf, "conf", "c", "", "配置文件位置")
	return cmd
}

type Server common.Conf

func (svr *Server) start() {
	//http.ListenAndServe(":8080", createFSHandler())
	http.ListenAndServe(":8081", createListener())
	http.ListenAndServe(":8082", createFSHandler())
}

func createTestListener() common.Handler {
	handler := common.Handler(func(w http.ResponseWriter, req *http.Request) {
		sess, _ := sm.SessionStart(w, req)
		logger.Infof("wowo %s", sess.Get("cc"))
		sess.Set("cc", "bb")
		logger.Info("a request")
		w.Write([]byte("a test"))
	})
	return handler
}

//初始化整个Filesystem
func initFileSystem() {
	fs.RootStorePath = "D:/temp/fs.binary"
	fs.RecoverFromStore()
}

//该handler用来处理node节点定时发送的ping信息
func createListener() common.Handler {
	node.DoSomework()
	handler := common.Handler(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			uri, _ := url.Parse(req.RequestURI)
			if strings.HasPrefix(uri.Path, "/node") {
				aa := req.Header.Get("AdviseAddress")
				bb, _ := ioutil.ReadAll(req.Body)
				res := node.HandleNodeRequest(aa, bb)
				w.Write(res)
			} else if strings.HasPrefix(uri.Path, "/cli/load") {
				userName := uri.Query().Get("user")
				blocksize, _ := strconv.Atoi(req.FormValue("blocks"))
				res := node.HandleClientRequest(req.FormValue("filename"), blocksize, &user.User{Name: userName})
				w.Write(res)
			}
		}

	})
	return handler
}

//开启服务器，监听来自client端的关于filesystem的请求
//这些请求包含了 rm,touch,mkdir,chmod,chown,mv,ll,ls等等修改文件系统文件树的请求，
//也包含了load这个将本地文件加载到gfs中的请求
//该handler专门用来处理关于fs文件树的操作
func createFSHandler() common.Handler {
	handler := common.Handler(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			sess, _ := sm.SessionStart(w, req)
			uri, _ := url.Parse(req.RequestURI)
			logger.Infof("receive request with uri: %s and session id :%s", uri.Path, sess.SessionId())
			if strings.HasPrefix(uri.Path, "/user/login") {
				uu := map[string]string{}
				common.DecodeFromReader(&uu, req.Body)
				sess.Set(SessionUser, &user.User{Name: uu["name"]})
				logger.Infof("receive request for login,with username '%s' ", uu["name"])
				w.Write(common.EncodeToBytes(SuccessLogin))
			} else if strings.HasPrefix(uri.Path, "/fs") {
				if sess.Get(SessionUser) == nil {
					w.Write(common.EncodeToBytes(InvalidSession))
				} else {
					bb := fs.Handle(uri.Path, sess.Get(SessionUser).(*user.User), parseBody(req.Body))
					w.Write(bb)
				}
			} else if strings.HasPrefix(uri.Path, "/load") {
				if uu := sess.Get(SessionUser); uu == nil {
					w.Write(common.EncodeToBytes(InvalidSession))
				} else {
					req.ParseForm()
					file := req.FormValue("file")
					size, _ := strconv.ParseInt(req.FormValue("size"), 10, 64)

					gfsw := fs.Load(file, size, uu.(*user.User))
					w.Write(common.EncodeToBytes(*gfsw))
				}
			} else if strings.HasPrefix(uri.Path, "/get") {

			} else if strings.HasPrefix(uri.Path, "/more") {
				if uu := sess.Get(SessionUser); uu == nil {
					w.Write(common.EncodeToBytes(InvalidSession))
				} else {
					req.ParseForm()
					file := req.FormValue("file")
					if gfsr, err := fs.Get(file, uu.(*user.User)); err == nil {
						w.Write(common.EncodeToBytes(*gfsr))
					} else {
						w.Write(common.EncodeToBytes(failMessage(err)))
					}
				}
			}
		} else {
			w.Write([]byte("仅支持POST请求"))
		}
	})
	return handler
}

func parseBody(body io.ReadCloser) string {
	var path string
	common.DecodeFromReader(&path, body)
	return path
}

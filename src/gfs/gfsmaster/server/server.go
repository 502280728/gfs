package server

//该package主要负责接收URL请求，并分发到相应的处理器，可以认为是MVC中的C
import (
	"gfs/gfsmaster/fs"
	"gfs/gfsmaster/fs/user"
	"gfs/gfsmaster/node"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Cmd() *cobra.Command {
	var conf string
	var cmd = &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			initFileSystem()
			createFSServer()
		},
	}
	cmd.Flags().StringVarP(&conf, "conf", "c", "", "配置文件位置")
	return cmd
}

//初始化整个Filesystem
func initFileSystem() {
	fs.RootStorePath = "D:/temp/fs.binary"
	fs.RecoverFromStore()
}

type Handler func(http.ResponseWriter, *http.Request)

func (handler Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler(w, req)
}

//开启服务器，监听来自client端的关于filesystem的请求
//这些请求包含了 rm,touch,mkdir,chmod,chown,mv,ll,ls等等修改文件系统文件树的请求，
//也包含了load这个将本地文件加载到gfs中的请求
func createFSServer() {
	http.ListenAndServe(":8080", createFSHandler())
	http.ListenAndServe(":8081", createListener())
}

//该handler用来处理node节点定时发送的ping信息

func createListener() Handler {
	node.DoSomework()
	handler := Handler(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			uri, _ := url.Parse(req.RequestURI)
			if strings.HasPrefix(uri.Path, "/node") {
				aa := req.Header.Get("AdviseAddress")
				bb, _ = ioutil.ReadAll(req.Body)
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

//该handler专门用来处理关于fs文件树的操作
func createFSHandler() Handler {
	handler := Handler(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			uri, _ := url.Parse(req.RequestURI)
			if strings.HasPrefix(uri.Path, "/fs") {
				userName := uri.Query().Get("user")
				bb := fs.Handle(uri.Path, &user.User{Name: userName}, parseBody(req.Body))
				w.Write(bb)
			} else if strings.HasPrefix(uri.Path, "/file") {

			}
		} else {
			w.Write([]byte("仅支持POST请求"))
		}

	})
	return handler
}

func parseBody(body io.ReadCloser) string {
	b, _ := ioutil.ReadAll(body)
	return string(b)
}

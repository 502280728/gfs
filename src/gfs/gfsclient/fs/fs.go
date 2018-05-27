package fs

import (
	"bufio"
	"bytes"
	"fmt"
	"gfs/common"
	http1 "gfs/common/http"
	"gfs/common/http/cookie"
	"gfs/gfsclient/cmd"
	"net/http"
	"net/url"
	"os"
	"strings"

	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"
)

var logger = logging.MustGetLogger("gfs/gfsclient/fs")
var CS = &cookie.GFSCookieStore{[]*http.Cookie{}}
var REQ = http1.GFSRequest{}

func init() {
	REQ.SetCookieStore(CS)
}

func Cmd() *cobra.Command {
	var master string
	var user string
	var pass string
	var cmd = &cobra.Command{
		Use: "login",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("connecting to master : %s", master)
			startFsCli(master, user)
		},
	}
	cmd.Flags().StringVarP(&master, "master", "m", "", "master的位置")
	cmd.Flags().StringVarP(&user, "user", "u", "", "用户名")
	cmd.Flags().StringVarP(&pass, "pass", "p", "", "密码")
	return cmd
}

func login(master string, user string, pass string) bool {
	uu := map[string]string{"name": user, "pass": pass}
	var msg common.MessageInFS
	err := REQ.PostObj(master+"/fs/login", uu, &msg)
	if err == nil {
		if msg.Success {
			logger.Info(msg.Data)
			return true
		} else {
			panic(msg.Msg)
		}
	} else {
		panic(err.Error())
	}
}

//  cd /d
//  ls
//  ls /
//  mkdir -p /a/b/c
//  rm -rf /a
//  touch /a.txt
//  mv /a/b /c/d
//  chmod u+a /ab/c/d
//  chown -r root /ab/c/d
func startFsCli(master string, user string) {
	uri, _ := url.Parse(master)
	fullPrint := fmt.Sprintf("[%s@%s] > ", user, uri.Host)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(fullPrint)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "exit" {
			os.Exit(0)
		}
		ss := strings.Split(line, " ")
		if ss[0] != "load" {
			cmd.MainCmd.ExecuteString(ss)
			//sentCommand(ss[0], ss[1])
		} else {

		}
		fmt.Print(fullPrint)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("there is an error")
	}
}

func sentCommand(cmdLine string, path string) {
	var bb bytes.Buffer
	bb.WriteString(path)
	resp, err := http.Post("http://localhost:8080/fs/"+cmdLine, "application/octet-stream", &bb)
	if err != nil {
		logger.Info("error occurs")
	}
	defer resp.Body.Close()
	var msg common.MessageInFS
	common.DecodeFromReader(&msg, resp.Body)
	if msg.Success {
		fmt.Println(msg.Data)
	} else {
		fmt.Println(msg.Msg)
	}
}
func sentCommand1(cmdLine string, targetfile string, sourcefile string) {
	//	var length int64
	//	if fs, err := os.Stat(sourcefile); err == nil {
	//		length = fs.Size()
	//	}
	//	if length%common.BlockSize == 0 {
	//		blocks = int(length / common.BlockSize)
	//	} else {
	//		blocks = int(length/common.BlockSize) + 1
	//	}
}

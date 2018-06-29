package fs

import (
	"bufio"
	"bytes"
	"fmt"
	"gfs/common"
	http1 "gfs/common/http"
	"gfs/common/http/cookie"
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
			logger.Infof("connecting to master : %s,%s,%s", master, user, pass)
			if login(master, user, pass) {
				startFsCli(master, user)
			}
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
	err := REQ.PostObj(master+"/user/login", uu, &msg)
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
//  load sourceFile targetFile
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
		if ss[0] != "load" && ss[0] != "more" {
			//cmd.MainCmd.ExecuteString(ss)
			sentCommand(master, ss[0], ss[1])
		} else if ss[0] == "load" {
			sentLoadCommand(master, ss[1], ss[2])
		} else if ss[0] == "more" {
			sendMoreCommand(master, ss[1])
		}
		fmt.Print(fullPrint)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("there is an error")
	}
}

func sentCommand(master, cmdLine, path string) {
	var bb bytes.Buffer
	bb.WriteString(path)
	var msg common.MessageInFS
	err := REQ.PostObj(master+"/fs/"+cmdLine, path, &msg)
	if err != nil {
		logger.Errorf("error occurs when sending command:%s", err.Error())
	} else {
		if msg.Success {
			fmt.Println(msg.Data)
		} else {
			fmt.Println(msg.Msg)
		}
	}
}
func sentLoadCommand(master, sourcefile, targetfile string) {
	sentCommand(master, "touch", targetfile)
	var length int64
	if fs, err := os.Stat(sourcefile); err == nil {
		length = fs.Size()
	} else {
		logger.Error(err.Error())
	}
	url := fmt.Sprintf("%s/load?file=%s&size=%d", master, targetfile, length)
	logger.Info(url)
	var res common.GFSWriter
	REQ.Post(url, nil, nil, &res)
	file, _ := os.Open(sourcefile)
	var bb = make([]byte, length, length)
	file.Read(bb)
	res.Write(bb)
}

func sendMoreCommand(master, file string) {
	url := fmt.Sprintf("%s/more?file=%s", master, file)
	var res common.GFSReader
	REQ.Post(url, nil, nil, &res)
	logger.Info(res)
	var bb = make([]byte, 1500, 1500)

	length, _ := res.Read(bb)
	fmt.Print(string(bb[0:length]))
	length, _ = res.Read(bb)
	fmt.Print(string(bb[0:length]))
	length, _ = res.Read(bb)
	fmt.Print(string(bb[0:length]))
	length, _ = res.Read(bb)
	fmt.Print(string(bb[0:length]))
}

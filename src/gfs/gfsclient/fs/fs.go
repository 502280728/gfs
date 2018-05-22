package fs

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"gfs/common"
	logging "github.com/op/go-logging"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var logger = logging.MustGetLogger("gfs/gfsclient/fs")

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
			sentCommand(ss[0], ss[1])
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
	var length int64
	if fs, err := os.Stat(sourcefile); err == nil {
		length = fs.Size()
	}
	var blocks int
	if length%common.BlockSize == 0 {
		blocks = int(length / common.BlockSize)
	} else {
		blocks = int(length/common.BlockSize) + 1
	}

	var bb bytes.Buffer
	bb.WriteString(path)
	var data url.Values
	data.Set("filename", targetfile)
	data.Set("blocks", strconv.Itoa(blocks))
	resp, err := http.PostForm("http://localhost:8081/cli/load", data)
	if err != nil {
		logger.Info("error occurs")
	}
	defer resp.Body.Close()
	var msg common.MasterToClientMessage
	common.DecodeFromReader(&msg, resp.Body)
	var bb bytes.Buffer
	file, _ := os.Open(sourcefile)
	for i := 0; i < blocks; i++ {
		var kk = make([]byte, 0, common.BlockSize)
		file.Read(kk)
		bb.Write(bb)
		http.Post(msg.Nodes[0]+"/data", "application/octet-stream", &bb)
		bb.Reset()
	}

}

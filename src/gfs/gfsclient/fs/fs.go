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
	"strings"
)

var logger = logging.MustGetLogger("gfs/gfsclient/fs")

func Cmd() *cobra.Command {
	var master string
	var user string
	var cmd = &cobra.Command{
		Use: "login",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Infof("connecting to master : %s", master)
			startFsCli(master, user)
		},
	}
	cmd.Flags().StringVarP(&master, "master", "m", "", "master的位置")
	cmd.Flags().StringVarP(&user, "user", "u", "", "用户名")
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
		sentCommand(ss[0], ss[1])
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
	dec := gob.NewDecoder(resp.Body)
	dec.Decode(&msg)
	if msg.Success {
		fmt.Println(msg.Data)
	} else {
		fmt.Println(msg.Msg)
	}
}

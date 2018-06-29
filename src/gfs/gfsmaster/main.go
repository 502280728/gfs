package main

import (
	"gfs/gfsmaster/server"

	"github.com/spf13/cobra"
)

// gfsmaster.exe start 启动主节点
var gfsmasterCmd = &cobra.Command{
	Use: "gfsmaster",
}

func main() {
	gfsmasterCmd.AddCommand(server.Cmd())
	gfsmasterCmd.Execute()
}

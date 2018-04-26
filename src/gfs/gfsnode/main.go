package main

import (
	"gfs/gfsnode/server"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use: "gfsnode",
}

func main() {
	mainCmd.AddCommand(server.Cmd())
	mainCmd.Execute()
}

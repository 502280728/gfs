package main

import (
	"gfs/gfsclient/fs"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use: "gfsclient",
}

// gfsclient fs ls -d /
// gfsclient fs mkdir -p /a/b/c
// gfsclient fs rm -rf /a
// gfsclient fs touch /a.txt
// gfsclient fs mv /a/b /c/d
// gfsclient fs chmod u+a /ab/c/d
// gfsclient fs chown -r root /ab/c/d

func main() {
	mainCmd.AddCommand(fs.Cmd())
	mainCmd.Execute()
}

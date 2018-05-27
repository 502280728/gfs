package main

import (
	"gfs/gfsclient/fs"
	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use: "gfsclient",
}

//  ls
//  ls /
//  mkdir -p /a/b/c
//  rm -rf /a
//  touch /a.txt
//  mv /a/b /c/d
//  chmod u+a /ab/c/d
//  chown -r root /ab/c/d

func main() {
	mainCmd.AddCommand(fs.Cmd())
	mainCmd.Execute()
}

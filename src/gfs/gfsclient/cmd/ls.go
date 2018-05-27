package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var LsCmd = &cobra.Command{
	Use: "ls",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("this is ls")
	},
}

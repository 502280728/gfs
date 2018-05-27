package cmd

import (
	"github.com/spf13/cobra"
)

var MainCmd = &cobra.Command{
	Use: "main",
}

func init() {
	MainCmd.AddCommand(LsCmd)
}

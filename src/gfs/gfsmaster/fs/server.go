package fs

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	var conf string
	var cmd = &cobra.Command{
		Use: "start",
	}
	cmd.Flags().StringVarP(&conf, "conf", "c", "", "配置文件位置")
	return cmd
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hidemaruowo/ytac/lib"
)

func versionCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "version",
		Short: "show version",
		Long:  "show YTAC's version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("v" + lib.Version())
		},
	}
	return cmd
}

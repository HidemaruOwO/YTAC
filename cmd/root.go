package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/hidemaruowo/ytac/lib"
)

// RootCmd is root command
var RootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		lib.CheckDiffVersion()
		fmt.Println("YTAC v" + lib.Version())
		fmt.Println("✨ Thank you for installing YTAC!!")
		fmt.Println("If you would like to confirm any commands, Run:")
		color.New(color.Bold).Println("\t" + color.BlueString("$ ") + "ytac --help")
	},
}

func init() {
	lib.CheckDiffVersion()
	cobra.OnInitialize()
	RootCmd.AddCommand(
		versionCmd(),
		getCmd(),
	)
}

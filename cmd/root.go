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
		fmt.Println("YTAC v" + lib.Version())
		fmt.Println("✨ Thanks for installing YTAC!!")
		fmt.Println("To verify the command, run:")
		color.New(color.Bold).Println("\t" + color.BlueString("$ ") + "ytac --help")
	},
}

/*
func init() {
	cobra.OnInitialize()
	RootCmd.AddCommand(
		sumCmd(),
	)
}
*/

package lib

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func CheckCmdFFMPEG() {
	var printBold = color.New(color.Bold)
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		fmt.Println("🔥 Failed to run ffmpeg command\nPlease install ffmpeg and set env path")
		printBold.Println("🔎 Download Page: " + color.HiBlueString("https://ffmpeg.org/download.html"))
		os.Exit(1)
	}
}

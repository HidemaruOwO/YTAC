package cmd

import (
	"fmt"
	"io"
	"os"

	"path/filepath"

	"github.com/fatih/color"
	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"

	"github.com/hidemaruowo/ytac/lib"
)

var videoPath = ""

func getCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Download a video from YouTube",
		Long:  `Download a video from YouTube`,
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] != "" {
				download(args[0])
			} else {
				var errorMessage string = color.HiRedString("🔥 Please type a video ID")
				fmt.Println(errorMessage + "\nRun:")
				color.New(color.Bold).Println("\t" + color.BlueString("$ ") + "ytac get <video ID or video URL>")
			}
		},
	}

	return cmd
}

func download(videoID string) {
	var client = youtube.Client{}
	var printBold = color.New(color.Bold)

	var thumbnail string = "https://img.youtube.com/vi/" + videoID + "/hqdefault.jpg"

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	var formats = video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	videoPath = filepath.Join(lib.GetYtacPath(), "temp", video.Title+"-"+videoID+".mp4")

	file, err := os.Create(videoPath)
	if err != nil {
		fmt.Println("🔥 Failed to create video file")
		lib.GenTempDirectory()
		printBold.Println("♻️  Restarting donwload function")
		download(videoID)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
	if lib.UseSixel() == true {
		lib.ShowImage(thumbnail)
	}
}

package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"path/filepath"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"

	"github.com/hidemaruowo/ytac/lib"
)

var videoPath = ""
var p [1]byte

var printBold = color.New(color.Bold)

func getCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Download a video from YouTube",
		Long:  `Download a video from YouTube`,
		Run: func(cmd *cobra.Command, args []string) {
			var isYouTubeURL = regexp.MustCompile(`^(?:https|http):\/\/(www\.youtube\.com)(?:\/(?:.*)|\?(?:.*)|$)$`)
			var getVideoURL = regexp.MustCompile(`\?v=([^&]+)`)
			var index int = 0
			var value string = ""

			if len(args) != 0 {
				for index, value = range args {
					if isYouTubeURL.MatchString(value) {
						var videoID = getVideoURL.FindStringSubmatch(value)[1]
						ytac(videoID, index)
					} else {
						ytac(value, index)
					}
				}
				var tempPath string = filepath.Join(lib.GetYtacPath(), "temp")
				os.RemoveAll(tempPath)
			} else {
				var errorMessage string = color.HiRedString("🔥 Please type a video ID")
				fmt.Println(errorMessage + "\nRun:")
				color.New(color.Bold).Println("\t" + color.BlueString("$ ") + "ytac get <video ID or video URL>")
			}
		},
	}

	return cmd
}

func ytac(videoID string, index int) {
	checkCmdFFMPEG()
	printBold.Println("✨ " + strconv.Itoa(index+1) + ", Running YTAC...")
	var videoPath, videoTitle = download(videoID)
	var audioPath string = audioConv(videoPath, videoTitle)
	printBold.Println(color.HiYellowString("==>") + " Saved path: " + color.HiBlueString(audioPath))
}

func download(videoID string) (string, string) {
	var client = youtube.Client{}

	var thumbnail string = "https://img.youtube.com/vi/" + videoID + "/hqdefault.jpg"

	video, err := client.GetVideo(videoID)
	if err != nil {
		printBold.Println("🔥 " + color.HiRedString("No YouTube videos were found with that VideoID") + "\nThe video may not exist or may be a private video")
		panic(err)

	}

	printBold.Println("🔎 Found a " + color.HiBlueString(video.Title) + " video")

	if lib.UseSixel() == true {
		lib.ShowImage(thumbnail)
	}

	var formats = video.Formats.WithAudioChannels() // only get videos with audio
	stream, size, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	videoPath = filepath.Join(lib.GetYtacPath(), "temp", videoID+".mp4")

	file, err := os.Create(videoPath)
	if err != nil {
		fmt.Println("🔥 Failed to create video file")
		lib.GenTempDirectory()
		printBold.Println("♻️  Restarting donwload function")
		download(videoID)
	}
	defer file.Close()

	var tmpl = `{{ red "Downloading:" }} {{ bar . "[" (blue "=") (rndcolor "~>") "." "]"}} {{speed . | green }} {{percent .}}`

	var bar = pb.ProgressBarTemplate(tmpl).Start64(int64(size))

	//var reader = io.LimitReader(rand.Reader, int64(n))
	var barReader = bar.NewProxyReader(stream)

	_, err = io.Copy(file, barReader)
	if err != nil {
		panic(err)
	}
	bar.Finish()
	return videoPath, video.Title
}

func audioConv(videoPath string, videoTitle string) string {
	var today string = time.Now().Format("2006-01-02")
	var distPath string = path.Join(lib.GetYtacPath(), "dist")
	videoTitle = strings.Replace(videoTitle, "/", "", -1)
	videoTitle = strings.Replace(videoTitle, "\\", "", -1)
	// progress bar setting
	var tmpl = `{{ red "Converting:" }} {{ bar . "[" (blue "=") (rndcolor "~>") "." "]"}}`
	var max int64 = 100
	var bar = pb.ProgressBarTemplate(tmpl).Start64(max)

	var audioPath string = path.Join(distPath, today, videoTitle+".mp3")
	for i := 0; i < 70; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond * 30)
	}
	var err = ffmpeg_go.Input(videoPath).Output(audioPath).OverWriteOutput().Run()
	if err != nil {
		bar.Finish()
		fmt.Println("🔥 Failed to convert video to audio")
		if f, err := os.Stat(distPath); os.IsNotExist(err) || !f.IsDir() {
			lib.GenDistDirectory()
		}
		lib.GenDistTodayDirectory()
		printBold.Println("♻️  Restarting audioConv function")
		audioConv(videoPath, videoTitle)
	}
	for i := 0; i < 30; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond * 30)
	}
	bar.Finish()
	return audioPath
}

func checkCmdFFMPEG() {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		fmt.Println("🔥 Failed to run ffmpeg command\nPlease install ffmpeg and set env path")
		printBold.Println("🔎 Download Page: " + color.HiBlueString("https://ffmpeg.org/download.html"))
		os.Exit(1)
	}
}

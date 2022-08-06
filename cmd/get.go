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

	// "github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"

	"github.com/hidemaruowo/ytac/lib"
)

var videoPath string
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
			var value string

			if len(args) != 0 {
				// read url loop
				status := make(chan string)
				defer close(status)
				for index, value = range args {
					if isYouTubeURL.MatchString(value) {
						var videoID = getVideoURL.FindStringSubmatch(value)[1]
						go ytac(videoID, index, status)
					} else {
						go ytac(value, index, status)
					}
				}
				_ = <-status
				var tempPath string = filepath.Join(lib.GetYtacPath(), "temp")
				err := removeContents(tempPath)
				if err != nil {
					panic(err)
				}
			} else {
				var errorMessage string = color.HiRedString("🔥 Please type a video ID")
				fmt.Println(errorMessage + "\nRun:")
				color.New(color.Bold).Println("\t" + color.BlueString("$ ") + "ytac get <video ID or video URL>")
			}
		},
	}

	return cmd
}

func ytac(videoID string, index int, status chan<- string) {
	chAudioConv := make(chan string)
  var pathTitle [2]string
  var convertList [][2]string
	printBold.Println("✨ " + strconv.Itoa(index+1) + ", Running YTAC...")
	pathTitle[0], pathTitle[1] = download(videoID)
  convertList = append(convertList, pathTitle)

  // TODO for文で回せるようにしたい
  for _, pt := range convertList {
    go audioConv(videoPath, pt[1], chAudioConv)
  }
	audioPath := <-chAudioConv
	defer close(chAudioConv)
	printBold.Println(color.HiYellowString("==>") + " Saved path: " + color.HiBlueString(audioPath))
}

func download(videoID string) (string, string) {
  // TODO gorutineで表示がバラバラになってしまったので、表示内容をまとめて管理するようにする。
	client := youtube.Client{}

	// thumbnail := "https://img.youtube.com/vi/" + videoID + "/hqdefault.jpg"

	video, err := client.GetVideo(videoID)
	if err != nil {
		printBold.Println("🔥 " + color.HiRedString("No YouTube videos were found with that VideoID") + "\nThe video may not exist or may be a private video")
		panic(err)
	}

	printBold.Println("🔎 Found a " + color.HiBlueString(video.Title) + " video")

  // Sixel表示を無効化
	// if lib.UseSixel() == true {
	// 	lib.ShowImage(thumbnail)
	// }

	var formats = video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	videoPath = filepath.Join(lib.GetYtacPath(), "temp", videoID+".mp4")

	file, err := os.Create(videoPath)
	defer file.Close()
	if err != nil {
		fmt.Println("🔥 Failed to create video file")
		lib.GenTempDirectory()
		printBold.Println("♻️  Restarting donwload function")
		download(videoID)
		return videoPath, video.Title
	}

  _, err = io.Copy(file, stream)
  if err != nil {
    panic(err)
  }
	// var tmpl = `{{ red "Downloading:" }} {{ bar . "[" (blue "=") (rndcolor "~>") "." "]"}} {{speed . | green }} {{percent .}}`
	//
  // FIXME 並列で動くようにした結果、Convertのプログレスバーと重なってしまったので修正する
	// var bar = pb.ProgressBarTemplate(tmpl).Start64(int64(size))

	//var reader = io.LimitReader(rand.Reader, int64(n))
	// var barReader = bar.NewProxyReader(stream)

	// _, err = io.Copy(file, barReader)
	// if err != nil {
	// 	panic(err)
	// }
	// bar.Finish()
  fmt.Println("Done!")
	return videoPath, video.Title
}

func audioConv(videoPath string, videoTitle string, chAudiPath chan string) {
	var today string = time.Now().Format("2006-01-02")
	var distPath string = path.Join(lib.GetYtacPath(), "dist")
	videoTitle = strings.Replace(videoTitle, "/", "", -1)
	videoTitle = strings.Replace(videoTitle, "\\", "", -1)
	// progress bar setting
	// var tmpl = `{{ red "Converting:" }} {{ bar . "[" (blue "=") (rndcolor "~>") "." "]"}} {{percent .}}`
	// var max int64 = 100
	// var bar = pb.ProgressBarTemplate(tmpl).Start64(max)

	var audioPath string = path.Join(distPath, today, videoTitle+".mp3")
	// for i := 0; i < 70; i++ {
	// 	bar.Increment()
	// 	time.Sleep(time.Millisecond * 30)
	// }
	var err = ffmpeg_go.Input(videoPath).Output(audioPath).OverWriteOutput().Run()
	if err != nil {
		// bar.Finish()
		fmt.Println("🔥 Failed to convert video to audio")
		if f, err := os.Stat(distPath); os.IsNotExist(err) || !f.IsDir() {
			lib.GenDistDirectory()
		}
		lib.GenDistTodayDirectory()
		printBold.Println("♻️  Restarting audioConv function")
		audioConv(videoPath, videoTitle, chAudiPath)
	}
	// for i := 0; i < 30; i++ {
	// 	bar.Increment()
	// 	time.Sleep(time.Millisecond * 30)
	// }
	// bar.Finish()
	chAudiPath <- audioPath
}

func CheckCmdFFMPEG() {
  // FIXME mainで呼び出すように変更したので、libに移動するのが望ましい
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		fmt.Println("🔥 Failed to run ffmpeg command\nPlease install ffmpeg and set env path")
		printBold.Println("🔎 Download Page: " + color.HiBlueString("https://ffmpeg.org/download.html"))
		os.Exit(1)
	}
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

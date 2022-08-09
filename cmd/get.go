package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"path/filepath"

	// "github.com/cheggaaa/pb/v3"
	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
	"github.com/kkdai/youtube/v2"
	"github.com/spf13/cobra"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"

	"github.com/hidemaruowo/ytac/lib"
)

var converted int = 0
var videoPath string
var p [1]byte
var convertList [][2]string
var savedPathes []string

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

			var videoIDs []string

			for _, value = range args {
				if isYouTubeURL.MatchString(value) {
					videoIDs = append(videoIDs, getVideoURL.FindStringSubmatch(value)[1])
				} else {
					videoIDs = append(videoIDs, value)
				}
			}

			args = removeArrayDuplicate(videoIDs)

			if len(args) != 0 {
				for _, value = range args {
					fmt.Println(value)
				}
				// read url loop
				status := make(chan string)
				defer close(status)
				for index, value = range args {
					runDownload(value)
					go ytac(value, index, status)
				}
				// progress bar setting
				var tmpl = `{{ red "Converting:" }} {{ bar . "[" (blue "=") (rndcolor "~>") "." "]"}} {{percent .}}`
				var max int64 = int64(len(args))
				var bar = pb.ProgressBarTemplate(tmpl).Start64(max)
				var applyConverted int = converted
				for i := 0; i < applyConverted; i++ {
					bar.Increment()
					time.Sleep(time.Millisecond * 30)
				}
				for {
					if applyConverted != converted {
						for i := 0; converted-applyConverted > i; i++ {
							bar.Increment()
							time.Sleep(time.Millisecond * 30)
						}
						fmt.Println(applyConverted)
						fmt.Println(converted)
					} else {
						break
					}
					time.Sleep(time.Millisecond * 500)
				}

				bar.Finish()

				for _, value = range savedPathes {
					printBold.Println(color.HiYellowString("==>") + " Saved path: " + color.HiBlueString(value))
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

// TODOEND Downloadを同期処理で実行して、コンバート処理を非同期で実行する
// TODOEND その後は変換処理が終わったあと、出力先をまとめて出力する
func ytac(videoID string, index int, status chan<- string) {
	// TODO for文で回せるようにしたい
	for _, pt := range convertList {
		chAudioConv := make(chan string)
		go audioConv(videoPath, pt[1], chAudioConv)
		audioPath := <-chAudioConv
		savedPathes = append(savedPathes, audioPath)
		defer close(chAudioConv)
	}
}

func download(videoID string) (string, string) {
	var client = youtube.Client{}

	var thumbnail string = "https://img.youtube.com/vi/" + videoID + "/hqdefault.jpg"

	video, err := client.GetVideo(videoID)
	if err != nil {
		printBold.Println("🔥 " + color.HiRedString("No YouTube videos were found with that VideoID") + "\nThe video may not exist or may be a private video")
		os.Exit(1)
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
		return videoPath, video.Title
	}
	defer file.Close()

	printBold.Println("🔎 Found a " + color.HiBlueString(video.Title))
	if lib.UseSixel() == true {
		lib.ShowImage(thumbnail)
	}

	var tmpl = `{{ red "Downloading:" }} {{ bar . "[" (blue "=") (rndcolor "~>") "." "]"}} {{speed . | green }} {{percent .}}`
	var bar = pb.ProgressBarTemplate(tmpl).Start64(int64(size))
	var barReader = bar.NewProxyReader(stream)
	_, err = io.Copy(file, barReader)
	if err != nil {
		panic(err)
	}
	bar.Finish()

	return videoPath, video.Title
}

func audioConv(videoPath string, videoTitle string, chAudioPath chan string) {
	var today string = time.Now().Format("2006-01-02")
	var distPath string = path.Join(lib.GetYtacPath(), "dist")
	videoTitle = strings.Replace(videoTitle, "/", "", -1)
	videoTitle = strings.Replace(videoTitle, "\\", "", -1)

	var audioPath string = path.Join(distPath, today, videoTitle+".mp3")
	log.SetOutput(ioutil.Discard)
	var err = ffmpeg_go.Input(videoPath).Output(audioPath).OverWriteOutput().Run()
	if err != nil {
		// bar.Finish()

		fmt.Println("🔥 Failed to convert video to audio")
		if f, err := os.Stat(distPath); os.IsNotExist(err) || !f.IsDir() {
			lib.GenDistDirectory()
		}
		lib.GenDistTodayDirectory()
		printBold.Println("♻️  Restarting audioConv function")
		audioConv(videoPath, videoTitle, chAudioPath)
	}
	chAudioPath <- audioPath
	converted += 1
}

func runDownload(videoID string) {
	var pathTitle [2]string
	pathTitle[0], pathTitle[1] = download(videoID)
	convertList = append(convertList, pathTitle)

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

func removeArrayDuplicate(args []string) []string {
	results := make([]string, 0, len(args))
	encountered := map[string]bool{}
	for i := 0; i < len(args); i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
			results = append(results, args[i])
		}
	}
	return results
}

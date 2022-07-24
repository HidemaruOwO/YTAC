package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"

	"github.com/hidemaruowo/ytac/config"
)

func GenConfig() {
	var printBold = color.New(color.Bold)
	var ytacPath string = GetYtacPath()
	var OutYtacPath string = color.HiBlueString(ytacPath)

	fmt.Println("🔎 $YTACPATH: " + OutYtacPath)
	fmt.Println("🔨 Creating config..")
	var configData string = config.DefaultConfig()
	var err = ioutil.WriteFile(filepath.Join(ytacPath, "config.json"), []byte(configData), 0644)
	if err != nil {
		fmt.Println("🔥 Failed to create config")
		var isGenYtacDir = GenYtacDirectory()
		if isGenYtacDir {
			return
		} else {
			printBold.Println("♻️  Restarting GenerateConfig function")
			fmt.Println("")
			GenConfig()
		}
	} else {
		printBold.Println("✨ Please restart YTAC")
		os.Exit(0)
	}
}

func genDirectory(path string, outPath string) bool {
	fmt.Println("🔨 Creating " + outPath + "..")
	var err = os.Mkdir(path, 0755)
	if err != nil {
		fmt.Println("🔥 Failed to create " + outPath)
		fmt.Errorf(err.Error())
		return false
	}
	return true
}

func GenYtacDirectory() bool {
	var ytacPath string = GetYtacPath()
	var outYtacPath string = color.HiBlueString(ytacPath)

	return genDirectory(ytacPath, outYtacPath)
}

func GenTempDirectory() bool {
	var ytacPath string = GetYtacPath()
	var tempPath string = filepath.Join(ytacPath, "temp")
	var outTempPath string = color.HiBlueString(tempPath)

	return genDirectory(tempPath, outTempPath)
}

func GenDistDirectory() bool {
	var ytacPath string = GetYtacPath()
	var distPath string = filepath.Join(ytacPath, "dist")
	var OutDistPath string = color.HiBlueString(distPath)

	return genDirectory(distPath, OutDistPath)
}

func GenDistTodayDirectory() bool {
	var ytacPath string = GetYtacPath()
	var today string = time.Now().Format("2006-01-02")
	var todayPath string = filepath.Join(ytacPath, "dist", today)
	var outTodayPath string = color.HiBlueString(todayPath)

	return genDirectory(todayPath, outTodayPath)
}

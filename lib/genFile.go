package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

func GenYtacDirectory() bool {
	var ytacPath string = GetYtacPath()
	var OutYtacPath string = color.HiBlueString(ytacPath)

	fmt.Println("🔨 Creating " + OutYtacPath + "..")
	var err = os.Mkdir(ytacPath, 0755)
	if err != nil {
		fmt.Println("🔥 Failed to create " + OutYtacPath)
		fmt.Errorf(err.Error())
		return false
	}
	return true
}

func GenTempDirectory() bool {
	var ytacPath string = GetYtacPath()
	var tempPath string = filepath.Join(ytacPath, "temp")
	var OutTempPath string = color.HiBlueString(tempPath)

	fmt.Println("🔨 Creating " + OutTempPath + "..")
	var err = os.Mkdir(ytacPath, 0755)
	if err != nil {
		fmt.Println("🔥 Failed to create " + OutTempPath)
		fmt.Errorf(err.Error())
		return false
	}
	return true

}

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
	fmt.Println("🔨 creating config..")
	var configData string = config.DefaultConfig()
	var err = ioutil.WriteFile(filepath.Join(ytacPath, "/config.json"), []byte(configData), 0644)
	if err != nil {
		fmt.Println("🔥 failed to create config")
		fmt.Println("🔨 creating " + OutYtacPath + "..")
		err = os.Mkdir(ytacPath, 0755)
		if err != nil {
			fmt.Println("🔥 failed to create " + OutYtacPath)
		} else {
			printBold.Println("♻️ Restarting GenerateConfig function")
			fmt.Println("")
			GenConfig()
		}
	} else {
		printBold.Println("✨ Please restart YTAC")
	}
}

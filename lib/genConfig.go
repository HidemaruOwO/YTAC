package lib

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
)

func GenConfig() {
	var printBold = color.New(color.Bold)
	var ytacPath string = GetYtacPath()
	var OutYtacPath string = color.HiBlueString(ytacPath)

	fmt.Println("🔎 $YTACPATH: " + OutYtacPath)
	fmt.Println("🔨 creating config..")
	var configData string = `{
	"name": "YouTube Video to Audio Converter",
	"version": "1.0.0"
}`
	var err = ioutil.WriteFile(ytacPath+"/config.json", []byte(configData), 0644)
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

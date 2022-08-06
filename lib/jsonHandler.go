package lib

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hidemaruowo/ytac/config"
	"github.com/tidwall/gjson"
)

var _config = readConfig()

//readJson
func readConfig() string {
	var ytacPath string = GetYtacPath()
	var configPath string = filepath.Join(ytacPath, "/config.json")

	f, err := os.Open(configPath)
	if err != nil {
		GenConfig()
		//log.Fatal(err)
		readConfig()
	}
	defer f.Close()
	json, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return string(json)
}

//Config
func Version() string {
	var value = gjson.Get(_config, "version")
	return value.String()
}

func UseSixel() bool {
	var value = gjson.Get(_config, "useSixel")
	return value.Bool()
}

//Functions
func CheckDiffVersion() {
	var json string = config.DefaultConfig()
	var version string = gjson.Get(json, "version").String()
	var _version string = Version()
	if _version != version {
		GenConfig()
	}
}

package lib

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/hidemaruowo/ytac/config"
	"github.com/tidwall/gjson"
)

var _config = readConfig()

//readJson
func readConfig() string {
	var ytacPath string = GetYtacPath()
	var configPath string = ytacPath + "/config.json"

	var f, err = os.Open(configPath)
	if err != nil {
		GenConfig()
		//log.Fatal(err)
		readConfig()
	}
	defer f.Close()
	var json, err2 = ioutil.ReadAll(f)
	if err2 != nil {
		log.Fatal(err2)
	}
	return string(json)
}

//Config
func Version() string {
	var version = gjson.Get(_config, "version")
	return version.String()
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

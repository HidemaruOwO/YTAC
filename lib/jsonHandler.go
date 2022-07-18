package lib

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/tidwall/gjson"
)

var config = readConfig()

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
	var version = gjson.Get(config, "version")
	return version.String()
}

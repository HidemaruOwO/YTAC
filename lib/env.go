package lib

import (
	"os"
)

func GetYtacPath() string {
	var ytacPath = os.Getenv("YTACPATH")

	if ytacPath == "" {
		var homeDir, _ = os.UserHomeDir()
		ytacPath = homeDir + "/ytac"
	}

	ytacPath = ytacPath + ""

	return ytacPath
}

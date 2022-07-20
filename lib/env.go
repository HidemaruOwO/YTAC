package lib

import (
	"os"
  "path/filepath"
)

func GetYtacPath() string {
	var ytacPath = os.Getenv("YTACPATH")

	if ytacPath == "" {
		var homeDir, _ = os.UserHomeDir()
    ytacPath = filepath.Join(homeDir, "/ytac")
	}

  // ?
	// ytacPath = ytacPath + ""

	return ytacPath
}

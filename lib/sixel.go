package lib

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mattn/go-sixel"
)

// ShowImage is display image
func ShowImage(url string) error {
	res, err := http.Get(url)
	if err != nil {
		fmt.Errorf("🔥 Cannot get image: %v: %v", err, res.Status)
	}
	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("🔥 Cannot read body %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(buf))

	if err != nil {
		return fmt.Errorf("🔥 Cannot decode image %v", err)
	}

	enc := sixel.NewEncoder(os.Stdout)
	return enc.Encode(img)
}

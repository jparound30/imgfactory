package main // imgfactory

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"html"
	"imgfactory/logger"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// load font file.
	_, err := getFont()
	if os.IsNotExist(err) {
		logger.Printf("Downloading IPA Gothic font file.")
		err = fontDownload()
		if err != nil {
			logger.Panicf("cannot download ipa gothic font file. %v", err)
			// never return
			return
		}
		logger.Printf("Finish\n")
		// retry loading font file.
		_, err = getFont()
	}
	if err != nil {
		logger.Panicf("cannot load font file. %v", err)
		// never return
		return
	}
	http.HandleFunc("/", imageGenerateHandler)

	logger.Fatal(http.ListenAndServe(":8080", nil))
}

func imageGenerateHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("[API] %q\n", html.EscapeString(r.URL.Path))
	paths := strings.Split(r.URL.Path, "/")
	resourceName := paths[len(paths)-1]

	var imageType ImageType = ImageJpeg
	if strings.Contains(resourceName, ".") {
		names := strings.Split(resourceName, ".")
		resourceName = names[0]
		if len(names) >= 2 {
			if names[len(names)-1] == "png" {
				imageType = ImagePng
			}
		}
	}

	var width, height int
	width = 640
	height = 480

	logger.Printf("W: %d, H: %d, file : %s\n", width, height, resourceName)

	dateString := time.Now().Format("2006/01/02 15:04:05 JST")
	var imageInfo = ImageInfo{
		width:     width,
		height:    height,
		text:      fmt.Sprintf("%s\n\n%s\n\n%s", r.URL.Path, resourceName, dateString),
		imageType: imageType,
	}

	ft, _ := getFont()
	opt := truetype.Options{
		Size:              36,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	face := truetype.NewFace(ft, &opt)

	buf, err := imageInfo.generateImage(&face)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if imageType == ImageJpeg {
		w.Header().Set("Content-Type", "image/jpeg")
	} else {
		w.Header().Set("Content-Type", "image/png")
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

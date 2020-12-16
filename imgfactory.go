package main // imgfactory

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"html"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type ImageInfo struct {
	width  int
	height int
	text   string
}

func (imageInfo ImageInfo) generateImage(fontFace *font.Face, imageType int) (*bytes.Buffer, error) {
	const textTopMargin = 90

	img := image.NewRGBA(image.Rect(0, 0, imageInfo.width, imageInfo.height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{
		R: 200,
		G: 50,
		B: 200,
		A: 255,
	}}, image.Point{}, draw.Over)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: *fontFace,
		Dot:  fixed.Point26_6{},
	}

	texts := strings.Split(imageInfo.text, "\n")

	for i, text := range texts {
		dr.Dot.X = (fixed.I(imageInfo.width) - dr.MeasureString(text)) / 2
		if i != 0 {
			bounds, _ := dr.BoundString(text)
			log.Printf("    %+v", bounds)
			dr.Dot.Y = fixed.I(textTopMargin) + fixed.I(i).Mul(bounds.Max.Y-bounds.Min.Y+fixed.I(20))
		} else {
			dr.Dot.Y = fixed.I(textTopMargin)
		}
		log.Printf("    %+v", dr.Dot)
		dr.DrawString(text)
	}

	buf := &bytes.Buffer{}
	var err error
	if imageType == 0 {
		err = jpeg.Encode(buf, img, nil)
	} else {
		err = png.Encode(buf, img)
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return buf, nil
}

func main() {
	// load font file.
	_, err := getFont()
	if os.IsNotExist(err) {
		log.Printf("IPAフォントをダウンロードします...")
		err = fontDownload()
		if err != nil {
			log.Panicf("can not download ipa gothic font file. %v", err)
			// never return
			return
		}
		log.Printf("完了\n")
		// retry loading font file.
		_, err = getFont()
	}
	if err != nil {
		log.Panicf("can not load font file. %v", err)
		// never return
		return
	}
	http.HandleFunc("/", imageGenerateHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func imageGenerateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[API] %q\n", html.EscapeString(r.URL.Path))
	paths := strings.Split(r.URL.Path, "/")
	idStr := paths[len(paths)-1]
	var imageType int
	if strings.Contains(idStr, ".") {
		ids := strings.Split(idStr, ".")
		idStr = ids[0]
		if len(ids) >= 2 {
			if ids[len(ids)-1] == "jpg" {
				imageType = 0
			} else if ids[len(ids)-1] == "png" {
				imageType = 1
			} else {
				imageType = 0
			}
		}
	}

	var width, height int
	width = 640
	height = 480
	log.Printf("W: %d, H: %d, file : %s\n", width, height, idStr)
	dateString := time.Now().Format("2006/01/02 15:04:05 JST")
	var imageInfo = ImageInfo{
		width:  width,
		height: height,
		text:   fmt.Sprintf("%s\n%s\n%s", r.URL.Path, idStr, dateString),
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

	buf, err := imageInfo.generateImage(&face, imageType)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if imageType == 0 {
		w.Header().Set("Content-Type", "image/jpeg")
	} else {
		w.Header().Set("Content-Type", "image/png")
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

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
	"strconv"
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
	log.Println(texts)

	for i, text := range texts {
		dr.Dot.X = (fixed.I(imageInfo.width) - dr.MeasureString(text)) / 2
		if i != 0 {
			bounds, _ := dr.BoundString(text)
			log.Printf("%+v", bounds)
			dr.Dot.Y = fixed.I(textTopMargin) + fixed.I(i).Mul(bounds.Max.Y-bounds.Min.Y+fixed.I(20))
		} else {
			dr.Dot.Y = fixed.I(textTopMargin)
		}
		log.Printf("%+v", dr.Dot)
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
	_, err := getFont()
	if err != nil {
		log.Panicf("can not load font file. %v", err)
		// never return
		return
	}
	http.HandleFunc("/cam/", imageGenerateHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func imageGenerateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello, %q\n", html.EscapeString(r.URL.Path))
	paths := strings.Split(r.URL.Path, "/")
	typ := paths[2]
	idStr := paths[3]
	var imageType int
	if strings.Contains(idStr, ".") {
		ids := strings.Split(idStr, ".")
		idStr = ids[0]
		if ids[1] == "jpg" {
			imageType = 0
		} else if ids[1] == "png" {
			imageType = 1
		} else {
			imageType = 0
		}
	}
	//wh := paths[2]
	//whs := strings.Split(wh, "x")
	//width, err := strconv.Atoi(whs[0])
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//height, err := strconv.Atoi(whs[1])
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//id := paths[3]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		idStr = "id unknown"
		id = 0
	}
	var width, height int
	if id%2 == 0 {
		width = 640
		height = 480
	} else {
		width = 1920
		height = 1080
	}
	log.Printf("W: %d\n", width)
	log.Printf("H: %d\n", height)
	log.Printf("id : %s\n", idStr)
	dateString := time.Now().Format("2006/01/02 15:04:05 JST")
	var imageInfo = ImageInfo{
		width:  width,
		height: height,
		text:   fmt.Sprintf("%s\n%s\n%s", idStr, typ, dateString),
	}

	ft, _ := getFont()
	opt := truetype.Options{
		Size:              48,
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

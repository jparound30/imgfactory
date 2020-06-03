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
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// load ttf
func loadFont() (*truetype.Font, error) {
	// フォントの読み込み
	ftBinary, err := ioutil.ReadFile("ttf/ipag.ttf")
	if err != nil {
		log.Printf("can not read font file from './ttf/ipag.ttf'. err:[%s]\n", err)
		return nil, err
	}
	ft, err := truetype.Parse(ftBinary)
	if err != nil {
		log.Printf("can not parse font file from './ttf/ipag.ttf'. err:[%s]\n", err)
		return nil, err
	}
	return ft, nil
}

type ImageInfo struct {
	width  int
	height int
	text   string
}

func (imageInfo ImageInfo) generateImage(fontFace *font.Face) (*bytes.Buffer, error) {
	const textTopMargin = 90

	img := image.NewRGBA(image.Rect(0, 0, imageInfo.width, imageInfo.height))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
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
	err := jpeg.Encode(buf, img, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return buf, nil
}

// render

var face font.Face

func main() {
	ft, err := loadFont()
	if err != nil {
		log.Println(err)
		return
	}
	opt := truetype.Options{
		Size:              48,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}
	face = truetype.NewFace(ft, &opt)

	http.HandleFunc("/image/", imageGenerateHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func imageGenerateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Hello, %q\n", html.EscapeString(r.URL.Path))
	paths := strings.Split(r.URL.Path, "/")
	wh := paths[2]
	whs := strings.Split(wh, "x")
	width, err := strconv.Atoi(whs[0])
	if err != nil {
		log.Println(err)
		return
	}
	height, err := strconv.Atoi(whs[1])
	if err != nil {
		log.Println(err)
		return
	}
	id := paths[3]
	log.Printf("W: %d\n", width)
	log.Printf("H: %d\n", height)
	log.Printf("id : %s\n", id)
	dateString := time.Now().Format("2006/01/02 15:04:05 JST")
	var imageInfo = ImageInfo{
		width:  width,
		height: height,
		text:   fmt.Sprintf("%s\n%s", id, dateString),
	}
	buf, err := imageInfo.generateImage(&face)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

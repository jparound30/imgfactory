package main

import (
	"bytes"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"imgfactory/logger"
	"strings"
)

type ImageType int

const (
	_ = iota
	ImageJpeg
	ImagePng
)

type ImageInfo struct {
	width     int
	height    int
	text      string
	imageType ImageType
}

const textTopMargin = 90
const leftMargin = 30

func (imageInfo ImageInfo) generateImage(fontFace *font.Face) (*bytes.Buffer, error) {
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

	dr.Dot.X = fixed.I(leftMargin)
	dr.Dot.Y = fixed.I(textTopMargin)

	for _, text := range texts {
		for _, c := range text {
			bounds, advance := dr.BoundString(string(c))
			logger.Printf("    %+v / %+v", bounds, advance)
			if bounds.Max.X > fixed.I(imageInfo.width-leftMargin) {
				// advance y and reset x initial position
				dr.Dot.X = fixed.I(leftMargin)
				dr.Dot.Y = dr.Dot.Y + fixed.I(36)
			}
			dr.DrawString(string(c))
		}
		// advance y and reset x initial position
		dr.Dot.X = fixed.I(leftMargin)
		dr.Dot.Y = dr.Dot.Y + fixed.I(36)
		logger.Printf("    %+v", dr.Dot)
	}

	buf := &bytes.Buffer{}
	var err error
	if imageInfo.imageType == ImagePng {
		err = png.Encode(buf, img)
	} else {
		err = jpeg.Encode(buf, img, nil)
	}

	if err != nil {
		logger.Println(err)
		return nil, err
	}
	return buf, nil
}

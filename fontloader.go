package main

import (
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"log"
)

var tf *truetype.Font

// load ttf
func loadFont() (*truetype.Font, error) {
	// フォントの読み込み
	const fontPath = "ttf/ipag.ttf"
	ftBinary, err := ioutil.ReadFile(fontPath)
	if err != nil {
		log.Printf("cannot read font file from '%s'. err:[%s]\n", fontPath, err)
		return nil, err
	}
	ft, err := truetype.Parse(ftBinary)
	if err != nil {
		log.Printf("cannot parse font file from './ttf/ipag.ttf'. err:[%s]\n", err)
		return nil, err
	}
	return ft, nil
}

func getFont() (*truetype.Font, error) {
	var err error
	if tf == nil {
		tf, err = loadFont()
	}
	return tf, err
}

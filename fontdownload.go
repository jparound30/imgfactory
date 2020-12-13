package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const IPAGothicDLUrl = "http://moji.or.jp/wp-content/ipafont/IPAfont/ipag00303.zip"
const FontFileName = "ipag.zip"
const FontDownloadDir = "ttf"

func fontDownload() error {
	response, err := http.Get(IPAGothicDLUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if _, err = os.Stat(FontDownloadDir); os.IsNotExist(err) {
		if err = os.Mkdir(FontDownloadDir, 0777); err != nil {
			return err
		}
	}

	zipData, err := ioutil.ReadAll(response.Body)
	err = ioutil.WriteFile(FontDownloadDir+string(os.PathSeparator)+FontFileName, zipData, 0777)
	defer os.Remove(FontDownloadDir + string(os.PathSeparator) + FontFileName)

	zipReader, err := zip.OpenReader(FontDownloadDir + string(os.PathSeparator) + FontFileName)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		if f.FileInfo().IsDir() {
			rc.Close()
			continue
		}
		buf := make([]byte, f.UncompressedSize64)
		_, err = io.ReadFull(rc, buf)
		if err != nil {
			rc.Close()
			return err
		}

		flatPath := filepath.Join(FontDownloadDir, path.Base(f.Name))
		if err = ioutil.WriteFile(flatPath, buf, f.Mode()); err != nil {
			rc.Close()
			return err
		}
	}
	return nil
}

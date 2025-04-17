package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"strings"
)

type ImageHandler interface {
	getImage() *Image
	setDecodedImage(image.Image)
}

type Image struct {
	Type   string          `json:"type"`
	Bitmap map[string]uint `json:"bitmap"`
	Width  uint            `json:"width"`
	Height uint            `json:"height"`
}

func (userMsg *UserMessage) getImage() *Image {
	return &userMsg.Image
}

func (userMsg *UserMessage) setDecodedImage(decodedImg image.Image) {
	// userMsg.DecodedImg = decodedImg
}

func (img *Image) getBytes() []byte {
	var bytes []byte
	for _, imgByte := range img.Bitmap {
		bytes = append(bytes, byte(imgByte))
	}
	return bytes
}

func decodeImage(img *Image) image.Image {
	var imgType string
	var imgBytes []byte
	var rd io.Reader

	imgType = strings.Split(img.Type, "/")[1]
	imgBytes = img.getBytes()
	rd = bytes.NewReader(imgBytes)

	var err error
	var decodedImg image.Image

	switch imgType {
	case "jpeg":
		decodedImg, err = jpeg.Decode(rd)
	case "png":
		decodedImg, err = png.Decode(rd)
	default:
		fmt.Println("Unknown image format")
	}

	if err != nil {
		log.Fatal(err)
	}

	return decodedImg
}

func (img *Image) redraw() {

}

func (img *Image) save() {
	bytes := img.getBytes()

	file, err := os.Create("/images/new image.png")
	if err != nil {
		log.Fatal(err)
	}

	file.Write(bytes)
	file.Close()

	// file.
}

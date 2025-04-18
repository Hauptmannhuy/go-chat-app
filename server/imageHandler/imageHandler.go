package imagehandler

import (
	"bytes"
	"cmp"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Image struct {
	Type   string          `json:"type"`
	Bitmap map[string]uint `json:"bitmap"`
	Width  uint            `json:"width"`
	Height uint            `json:"height"`
}

func LocalSave(img *Image) {
	imgBytes, err := img.getBytes()
	if err != nil {
		fmt.Println("data image is corrupted")
	}
	imgType := strings.Split(img.Type, "/")[1]

	var decodedImage image.Image = decodeImage(imgBytes, imgType)
	file, err := os.Create("./images/new_image.png")
	if err != nil {
		log.Fatal(err)
	}

	encodeImage(decodedImage, imgType, file)

	file.Close()
	fmt.Println("image saved")
}

func (img *Image) getBytes() ([]byte, error) {
	var bytes []byte
	var keys []int
	var sortedBitmap map[int]uint = map[int]uint{}

	for key, val := range img.Bitmap {
		i, err := strconv.Atoi(key)
		sortedBitmap[i] = val
		if err != nil {
			return nil, err
		}
		keys = append(keys, i)
	}

	slices.SortFunc(keys, func(a, b int) int {
		return cmp.Compare(a, b)
	})

	for _, key := range keys {
		bytes = append(bytes, byte(sortedBitmap[key]))
	}

	return bytes, nil
}

func decodeImage(imgBytes []byte, imgType string) image.Image {
	var rd io.Reader
	var decodedImg image.Image
	var err error

	if err != nil {
		fmt.Println("image data corrupted")
		return nil
	}

	rd = bytes.NewReader(imgBytes)

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

func encodeImage(img image.Image, imgType string, file *os.File) {
	var err error
	switch imgType {
	case "jpeg":
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
	case "png":
		err = png.Encode(file, img)
	}

	if err != nil {
		fmt.Printf("error encoding image: %s\n", err)
	}
}

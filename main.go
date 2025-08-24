package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func toRGBA(img image.Image) *image.RGBA {
	newImage := image.NewRGBA(img.Bounds())

	draw.Draw(newImage, img.Bounds(), img, img.Bounds().Min, draw.Src)
	
	return newImage
}

func main() {
	var filePath string

	flag.StringVar(&filePath, "i", "", "image file")

	flag.Parse()

	if filePath == "" {
		log.Fatal("choose an image")
	}

	reader, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	img, _, err := image.Decode(reader)

	if err != nil {
		log.Fatal("failed to parse image")
	}

	rgbaImage := toRGBA(img)

	res := BilinearInterpolation(rgbaImage, 400, 400)

	OutputImageForDebugResult(res, "./img/ScaleBilinear.jpg")

	fmt.Print("done")
}

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
)

func BilinearInterpolation(img image.Image, targetWidth, targetHeight int) image.Image {
	fmt.Printf("%T", img)
	switch imageType := img.(type) {
	case *image.Gray:
		break
	case *image.RGBA:
		return BilinearInterpolationRGBA(imageType, targetWidth, targetHeight)
	default:
		log.Fatal("no not implemented yet")
	}

	return nil
}

func BilinearInterpolationRGBA(img *image.RGBA, targetWidth, targetHeight int) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	scaleFactorX := float64(img.Bounds().Dx()) / float64(targetWidth)
	scaleFactorY := float64(img.Bounds().Dy()) / float64(targetHeight)

	maxX := float64(img.Bounds().Dx())
	maxY := float64(img.Bounds().Dy())

	for y := range newImg.Bounds().Dy() {
		for x := range newImg.Bounds().Dx() {
			originalX := float64(x) * scaleFactorX
			originalY := float64(y) * scaleFactorY

			// Fid The Closest Neighbor Between
			x1 := math.Floor(originalX)
			y1 := math.Floor(originalY)

			x2 := min(maxX-1, x1+1)
			y2 := min(maxY-1, y1+1)

			// Calculate Fractional Distance
			horizontalDistance := (originalX - x1) / (x2 - x1)
			verticalDistance := (originalY - y1) / (y2 - y1)

			topHorizontal := LinearInterpolation(horizontalDistance, img.RGBAAt(int(x1), int(y1)), img.RGBAAt(int(x2), int(y1)))
			bottomHorizontal := LinearInterpolation(horizontalDistance, img.RGBAAt(int(x1), int(y2)), img.RGBAAt(int(x2), int(y2)))

			finalValue := LinearInterpolation(verticalDistance, topHorizontal, bottomHorizontal)

			newImg.Set(x, y, finalValue)
		}
	}

	return newImg
}

func LinearInterpolation(dx float64, x1 color.RGBA, x2 color.RGBA) color.RGBA {

	r := calculateLinear(dx, float64(x1.R), float64(x2.R))
	g := calculateLinear(dx, float64(x1.G), float64(x2.G))
	b := calculateLinear(dx, float64(x1.B), float64(x2.B))
	a := calculateLinear(dx, float64(x1.A), float64(x2.A))

	return color.RGBA{
		R: uint8(clamp(r, 0, 255)),
		G: uint8(clamp(g, 0, 255)),
		B: uint8(clamp(b, 0, 255)),
		A: uint8(clamp(a, 0, 255)),
	}
}

func calculateLinear(dx float64, x1, x2 float64) float64 {
	return x1 + (dx * (x2 - x1))
}

func clamp(value, min, max float64) float64 {
	if value < min { return min }
	if value > max { return max }
	return value
}

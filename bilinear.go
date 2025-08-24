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

			// topHorizontal := linearInterpolation(horizontalDistance, img.RGBAAt(int(x1), int(y1)), img.RGBAAt(int(x2), int(y1)))
			// bottomHorizontal := linearInterpolation(horizontalDistance, img.RGBAAt(int(x1), int(y2)), img.RGBAAt(int(x2), int(y2)))
			//
			// finalValue := linearInterpolation(verticalDistance, topHorizontal, bottomHorizontal)

			finalValue := interpolationMathVersion(
				horizontalDistance, 
				verticalDistance, 
				img.RGBAAt(int(x1), int(y1)),
				img.RGBAAt(int(x2), int(y1)),
				img.RGBAAt(int(x1), int(y2)),
				img.RGBAAt(int(x2), int(y2)),
			)
			newImg.Set(x, y, finalValue)
		}
	}

	return newImg
}

func linearInterpolation(dx float64, x1 color.RGBA, x2 color.RGBA) color.RGBA {

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

// node1 Mean Top Left
// node2 Mean Top Right
// node3 Mean Bottom Left
// node4 Mean Bottom Right

func interpolationMathVersion(dx, dy float64, node1,node2,node3,node4 color.RGBA) color.RGBA {
	r := calculateWeightDistribution(dx, dy, node1.R, node2.R, node3.R, node4.R)
	g := calculateWeightDistribution(dx, dy, node1.G, node2.G, node3.G, node4.G)
	b := calculateWeightDistribution(dx, dy, node1.B, node2.B, node3.B, node4.B)
	a := calculateWeightDistribution(dx, dy, node1.A, node2.A, node3.A, node4.A)

return color.RGBA{
		R: uint8(clamp(r, 0, 255)),
		G: uint8(clamp(g, 0, 255)),
		B: uint8(clamp(b, 0, 255)),
		A: uint8(clamp(a, 0, 255)),
	}
}

func calculateWeightDistribution(dx, dy float64, color1,color2,color3,color4 uint8) float64 {
	topLeft := float64(color1) * ((1 - dx) * (1 - dy))
	topRight := float64(color2) * (dx - (1 - dy))

	bottomLeft := float64(color3) * ((1-dx) - dy)
	bottomRight := float64(color4) * ((1-dx) - (1-dy))

	return topLeft + topRight + bottomLeft + bottomRight
}

func clamp(value, min, max float64) float64 {
	if value < min { return min }
	if value > max { return max }
	return value
}

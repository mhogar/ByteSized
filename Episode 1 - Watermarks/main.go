package main

import (
	"image"
	"image/color"
	png "image/png"
	"log"
	"os"

	// include to initialize image formats
	_ "image/jpeg"
)

const numBits = 8

func main() {
	//createWatermark()
	extractWatermark()
}

func createWatermark() {
	// decode the images
	baseImg := decodeImage("imgs/landscape.jpeg")
	watermarkImg := decodeImage("imgs/github_logo.png")

	// get the bounds
	baseBounds := baseImg.Bounds()
	watermarkBounds := watermarkImg.Bounds()

	// create the output image with the same size as the base image
	outputImg := image.NewRGBA64(baseBounds)

	// calculate the offsets
	xOffset := calcOffset(baseBounds.Max.X, baseBounds.Min.X, watermarkBounds.Max.X, watermarkBounds.Min.X)
	yOffset := calcOffset(baseBounds.Max.Y, baseBounds.Min.Y, watermarkBounds.Max.Y, watermarkBounds.Min.Y)

	// loop through each pixel in the base image
	for y := baseBounds.Min.Y; y < baseBounds.Max.Y; y++ {
		for x := baseBounds.Min.X; x < baseBounds.Max.X; x++ {

			// calc the watermark x and y index
			watermarkX := x - xOffset
			watermarkY := y - yOffset

			// make sure the indices are in range
			if watermarkX < watermarkBounds.Min.X || watermarkX >= watermarkBounds.Max.X || watermarkY < watermarkBounds.Min.Y || watermarkY >= watermarkBounds.Max.Y {
				outputImg.Set(x, y, baseImg.At(x, y))
				continue
			}

			// extract the color from both the base and watermark image
			r1, g1, b1, a1 := baseImg.At(x, y).RGBA()
			r2, g2, b2, _ := watermarkImg.At(watermarkX, watermarkY).RGBA()

			// create a new color by combining the two values
			c := color.RGBA64{
				R: uint16(r1) | uint16(r2)>>(16-numBits),
				G: uint16(g1) | uint16(g2)>>(16-numBits),
				B: uint16(b1) | uint16(b2)>>(16-numBits),
				A: uint16(a1),
			}

			// the the color in the output image
			outputImg.SetRGBA64(x, y, c)
		}
	}

	// save the image
	saveImage("test.png", outputImg)
}

func extractWatermark() {
	// decode the image and get the bounds
	baseImg := decodeImage("test.png")
	baseBounds := baseImg.Bounds()

	// create the output image with the same size as the base image
	outputImg := image.NewRGBA64(baseBounds)

	// loop through each pixel in the base image
	for y := baseBounds.Min.Y; y < baseBounds.Max.Y; y++ {
		for x := baseBounds.Min.X; x < baseBounds.Max.X; x++ {
			// extract the color from both the base image
			r, g, b, a := baseImg.At(x, y).RGBA()

			// create the color using the base image
			c := color.RGBA64{
				R: (uint16(r) & ((1 << numBits) - 1)) << numBits,
				G: (uint16(g) & ((1 << numBits) - 1)) << numBits,
				B: (uint16(b) & ((1 << numBits) - 1)) << numBits,
				A: uint16(a),
			}

			// the the color in the output image
			outputImg.SetRGBA64(x, y, c)
		}
	}

	// save the image
	saveImage("extracted.png", outputImg)
}

func decodeImage(filename string) image.Image {
	// open the file and defer closing
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	// decode the image
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func calcOffset(max1, min1, max2, min2 int) int {
	return ((max1 - min1) - (max2 - min2)) / 2
}

func saveImage(filename string, img image.Image) {
	// create the output file and defer closing
	outputImgFile, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer outputImgFile.Close()

	// save the image to the file in the png format
	png.Encode(outputImgFile, img)
}

package main

import (
	"flag"
	"image"
	"image/color"
	png "image/png"
	"log"
	"os"

	// include to initialize image formats
	_ "image/jpeg"
)

func main() {
	// define flags
	extractMode := flag.Bool("e", false, "foo")
	baseImgFilename := flag.String("base", "", "foo")
	watermarkImgFilename := flag.String("watermark", "", "foo")
	outputImgFilename := flag.String("output", "", "foo")
	numBits := flag.Int("bits", 8, "foo")

	flag.Parse()

	//-- validate parameters --
	if *baseImgFilename == "" {
		println("Base Image Filename is required.")
		return
	}

	if *outputImgFilename == "" {
		println("Output Image Filename is required.")
		return
	}

	if *numBits < 1 || *numBits > 16 {
		println("Bits must be between 1 and 16.")
		return
	}

	// choose mode based on flag
	if *extractMode {

		// print configuration
		println("Base Image:", *baseImgFilename)
		println("Output Image:", *outputImgFilename)
		println("Number of bits:", *numBits)
		println("Extracting watermark...")

		extractWatermark(*baseImgFilename, *outputImgFilename, *numBits)
	} else {
		// validate mode specific parameters
		if *watermarkImgFilename == "" {
			println("Watermark Image Filename is required.")
			return
		}

		// print configuration
		println("Base Image:", *baseImgFilename)
		println("Watermark Image:", *watermarkImgFilename)
		println("Output Image:", *outputImgFilename)
		println("Number of bits:", *numBits)
		println("Creating watermark...")

		createWatermark(*baseImgFilename, *watermarkImgFilename, *outputImgFilename, *numBits)
	}
}

func createWatermark(baseImgFilename, watermarkImgFilename, outputImgFilename string, numBits int) {
	// decode the images
	baseImg := decodeImage(baseImgFilename)
	watermarkImg := decodeImage(watermarkImgFilename)

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
				R: combineColors(r1, r2, numBits),
				G: combineColors(g1, g2, numBits),
				B: combineColors(b1, b2, numBits),
				A: uint16(a1),
			}

			// the the color in the output image
			outputImg.SetRGBA64(x, y, c)
		}
	}

	// save the image
	saveImage(outputImgFilename, outputImg)
}

func combineColors(c1, c2 uint32, numBits int) uint16 {
	b := uint16(1<<16 - 1)
	offset := 16 - numBits

	return (uint16(c1) | b>>offset) & ((uint16(c2) >> offset) | b<<numBits)
}

func extractWatermark(baseImgFilename, outputImgFilename string, numBits int) {
	// decode the image and get the bounds
	baseImg := decodeImage(baseImgFilename)
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
				R: extractColor(r, numBits),
				G: extractColor(g, numBits),
				B: extractColor(b, numBits),
				A: uint16(a),
			}

			// the the color in the output image
			outputImg.SetRGBA64(x, y, c)
		}
	}

	// save the image
	saveImage(outputImgFilename, outputImg)
}

func extractColor(c uint32, numBits int) uint16 {
	return uint16(c) << (16 - numBits)
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

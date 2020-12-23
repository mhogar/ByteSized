package main

import (
	"image"
	png "image/png"
	"log"
	"os"

	// include to initialize image formats

	"image/color"
	_ "image/jpeg"
)

func main() {
	const numBits = 8

	dest := decodeImage("imgs/landscape.jpeg")
	src := decodeImage("imgs/github_logo.png")

	destBounds := dest.Bounds()
	srcBounds := src.Bounds()

	output := image.NewRGBA64(destBounds)

	xOffset := calcOffset(destBounds.Max.X, destBounds.Min.X, srcBounds.Max.X, srcBounds.Min.X)
	yOffset := calcOffset(destBounds.Max.Y, destBounds.Min.Y, srcBounds.Max.Y, srcBounds.Min.Y)

	for y := destBounds.Min.Y; y < destBounds.Max.Y; y++ {
		for x := destBounds.Min.X; x < destBounds.Max.X; x++ {
			srcX := x - xOffset
			if srcX < srcBounds.Min.X || srcX >= srcBounds.Max.X {
				output.Set(x, y, dest.At(x, y))
				continue
			}

			srcY := y - yOffset
			if srcY < srcBounds.Min.Y || srcY >= srcBounds.Max.Y {
				output.Set(x, y, dest.At(x, y))
				continue
			}

			// apply watermark
			r1, g1, b1, a1 := dest.At(x, y).RGBA()
			r2, g2, b2, _ := src.At(srcX, srcY).RGBA()

			c := color.RGBA64{
				R: uint16(r1) | uint16(r2)>>(16-numBits),
				G: uint16(g1) | uint16(g2)>>(16-numBits),
				B: uint16(b1) | uint16(b2)>>(16-numBits),
				A: uint16(a1),
			}

			output.SetRGBA64(x, y, c)
		}
	}

	outputFile, err := os.Create("test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	png.Encode(outputFile, output)
}

func decodeImage(filename string) image.Image {
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func calcOffset(destMax, destMin, srcMax, srcMin int) int {
	// calc x offset
	destWidth := destMax - destMin
	srcWidth := srcMax - srcMin

	return (destWidth - srcWidth) / 2
}

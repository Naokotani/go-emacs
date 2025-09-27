package images

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/Naokotani/resize"
)

func ResizePng(inputFile, outputFile string, height, width uint) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	img, err := png.Decode(in)
	if err != nil {
		return err
	}
	in.Close()

	m := resize.Resize(height, width, img, resize.Bilinear)

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	png.Encode(out, m)
	return nil
}

func IsThumbTooWide(filename string, target uint) bool {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		panic(err)
	}

	if uint(cfg.Width) > target*2 {
		return true
	}
	return false
}

func ResizeJpegToPng(inputFile, outputFile string, height, width uint) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(in)
	if err != nil {
		return err
	}
	in.Close()

	m := resize.Resize(height, width, img, resize.NearestNeighbor)

	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	png.Encode(out, m)
	return nil
}

package utils

import (
	"crypto/sha256"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func createIdenticon(data string, size int) image.Image {
	hash := sha256.Sum256([]byte(data))
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Background color
	bgColor := color.RGBA{255, 255, 255, 0}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Foreground colors
	fgColors := []color.RGBA{
		{hash[3], hash[4], hash[5], 255},
		{hash[6], hash[7], hash[8], 255},
		{hash[9], hash[10], hash[11], 255},
	}

	blockSize := size / 8
	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			if hash[(i*8+j)/8]&(1<<uint((i*8+j)%8)) != 0 {
				colorIndex := (i + j) % len(fgColors)
				drawShape(img, i, j, blockSize, fgColors[colorIndex])
				// Mirror vertically
				drawShape(img, 7-i, j, blockSize, fgColors[colorIndex])
			}
		}
	}

	return img
}

func drawShape(img *image.RGBA, x, y, size int, c color.Color) {
	for dx := 0; dx < size; dx++ {
		for dy := 0; dy < size; dy++ {
			img.Set(x*size+dx, y*size+dy, c)
		}
	}
}

func saveIdenticon(img image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func generatePath(alias string) string {
	dir := "static"
	return dir + "/" + alias + ".png"
}

func GenerateIdenticon(alias string, size int) (string, error) {
	img := createIdenticon(alias, size)
	path := generatePath(alias)
	err := saveIdenticon(img, path)
	if err != nil {
		return "", err
	}

	return path, nil
}

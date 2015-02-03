package identicon

import (
	"image"
	"image/color"
	"image/png"

	"bytes"
)

type Designed interface {
	Render(resolution int16, bgColour color.NRGBA) []byte
}

type designed struct {
	design [][]bool
	colour color.NRGBA
}

func (d *designed) Render(resolution, border int, bgColour color.NRGBA) []byte {
	return Render(resolution, border, d.design, d.colour, bgColour)
}

func Render(resolution, border int, design [][]bool, colour, bgColour color.NRGBA) []byte {
	canvas := image.Rect(0, 0, resolution, resolution)
	palette := color.Palette{
		bgColour,
		colour,
	}
	img := image.NewPaletted(canvas, palette)
	stretch := resolution / (len(design)+2*border)
	// This allows us to use copy to efficiently set many pixels in the image at once
	ones := make([]byte, stretch)
	for i := range ones {
		ones[i] = 1
	}
	for y, row := range design {
		for x, blockIsFilled := range row {
			if blockIsFilled {
				fill(img, x, y, border, stretch, ones)
			}
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)

	return buf.Bytes()
}

func fill(img *image.Paletted,
          x, y, border, stretch int,
          ones []byte) {
	trueX := (x+border) * stretch
	minY := (y+border) * stretch
	for trueY := minY; trueY < minY + stretch; trueY += 1 {
		offset := img.PixOffset(trueX, trueY)
		copy(img.Pix[offset: offset + stretch], ones)
	}
}


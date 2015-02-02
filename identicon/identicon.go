package identicon

import (
	"hash"
	"image"
	"image/color"
	"image/png"
	"fmt"
	"bytes"

	"github.com/dchest/siphash"
)

const (
	sumSize = 64 // size of the `sum` member of `Identicon`

	// Border around the identicon, in terms of block width
	xBorder = 1
	yBorder = 1
)

type Identicon struct {
	info    []byte
	hash    hash.Hash64
	sum     *uint64
	size    *int
	design  []bool
	colour  *color.NRGBA
}

func FromInfo(info []byte) *Identicon {
	return &Identicon{
		info: info,
		hash: nil,
		sum: nil,
		size: nil,
		design: nil,
		colour: nil,
	}
}

func FromHash(hash hash.Hash64) *Identicon {
	sum := hash.Sum64()
	return &Identicon{
		info: nil,
		hash: hash,
		sum: &sum,
		size: nil,
		design: nil,
		colour: nil,
	}
}

type Hasher interface {
	Hash() hash.Hash64
}

func (id *Identicon) Hash(hashKey []byte) hash.Hash64 {
	if id.hash == nil {
		id.hash = siphash.New(hashKey)
		id.hash.Write(id.info)
		tmp := id.hash.Sum64()
		id.sum = &tmp
	}
	return id.hash
}

type Designer interface {
	Design(size int) [][]bool
}

func (id *Identicon) Design(size int, hashKey []byte) ([][]bool, error) {
	if size < 0 {
		return nil, fmt.Errorf("Negative size: %d", size)
	}
	// if design is not nil, nBlocks must not be nil
	if id.design != nil && *id.size == size {
		return parseDesign(size, id.design), nil
	}
	id.size = &size
	id.Hash(hashKey)

	if id.colour == nil {
		const full = 0xff
		id.colour = &color.NRGBA{
			R: uint8(*id.sum),
			G: uint8(*id.sum >> 8),
			B: uint8(*id.sum >> 16),
			A: full,
		}
	}

	nBlocks := size * size
	id.design = make([]bool, nBlocks)

	// Cycle through the bits of the hash, using them to create the design
	const sumOffset uint = 24 // already used first 24 bits for colour
	for i := 0; i < nBlocks / 2; i += 1 {
		id.design[i] = (*id.sum >> ((sumOffset + uint(i)) % sumSize) & 1) == 1
	}

	return parseDesign(size, id.design), nil
}

func parseDesign(size int, design []bool) [][]bool {
	blocks := make([][]bool, size)

	for iRow := 0; iRow < size / 2; iRow += 1 {
		blocks[iRow] = make([]bool, size)
		copy(blocks[iRow], design[iRow*size : (iRow+1)*size - 1])
		blocks[size - 1 - iRow] = make([]bool, size)

		for i := 0; i < size; i += 1 {
			blocks[size - 1 - iRow][i] = blocks[iRow][size - 1 - i]
		}
	}

	if size % 2 == 1 {
		blocks[size/2] = make([]bool, size)
		copy(blocks[size/2], design[size/2*size : (size/2+1)*size - 1])
		for i := 0; i < size / 2; i += 1 {
			blocks[size/2][size - 1 - i] = blocks[size/2][i]
		}
	}

	return blocks
}

type Renderer interface {
	Render(size int16) []byte
}

func (id *Identicon) Render(size, nBlocks int,
                            hashKey []byte,
                            bgColour color.NRGBA) ([]byte, error) {
	design, err := id.Design(nBlocks, hashKey)
	if err != nil {
		return nil, err
	}
	return Render(size, design, *id.colour, bgColour), nil
}

func Render(size int, blocks [][]bool, colour, bgColour color.NRGBA) []byte {
	canvas := image.Rect(0, 0, size, size)
	palette := color.Palette{
		bgColour,
		colour,
	}
	img := image.NewPaletted(canvas, palette)


	border := image.Point{xBorder, yBorder}
	stretch := image.Point{
		size / (len(blocks[0]) + 2*xBorder),
		size / (len(blocks) + 2*yBorder),
	}
	// This allows us to use copy to efficiently set many pixels in the image at once
	ones := make([]byte, stretch.X)
	for i := range ones {
		ones[i] = 1
	}
	for y, row := range blocks {
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
          x, y int,
          border, stretch image.Point,
          ones []byte) {
	println(x, y)
	trueX := (x+border.X) * stretch.X
	minY := (y+border.Y) * stretch.Y
	for trueY := minY; trueY < minY + stretch.Y; trueY += 1 {
		offset := img.PixOffset(trueX, trueY)
		copy(img.Pix[offset: offset + stretch.X], ones)
	}
}


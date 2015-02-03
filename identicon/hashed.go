package identicon

import (
	"hash"

	"image/color"
)

type Hashed interface {
	Design(designSize int) Designed
}

func FromHash(hash hash.Hash64) *hashed {
	sum := hash.Sum64()
	return &hashed{hash, sum}
}

type hashed struct {
	hash hash.Hash64
	sum uint64
}

type designedHashed struct {
	*designed
	*hashed
}

func (h *hashed) Design(designSize int) *designedHashed {
	return h.DesignSymmetric(designSize)
}

func (h *designedHashed) Design(designSize int) *designedHashed {
	if len(h.design) == designSize {
		return h
	}
	return h.DesignSymmetric(designSize)
}

// Generates a symmetric identicon design of the specified size based 
// on the info and the identicon's preset information.
// designSize -- the number of blocks in the width/height of the identicon.
//               Must be positive
// 
// returns a designSize*designSize 2d slice containing the design.
func (h *hashed) DesignSymmetric(designSize int) *designedHashed {
	if designSize < 0 {
		panic("Negative design size")
	}

	const full = 0xff
	colour := color.NRGBA{
		R: uint8(h.sum),
		G: uint8(h.sum >> 8),
		B: uint8(h.sum >> 16),
		A: full,
	}
	nBlocks := designSize * (designSize / 2)

	// Middle row for odd number of rows should have half as many blocks generated
	if designSize % 2 == 1 {
		nBlocks += designSize/2 + 1 // Round up
	}
	blocks := make([]bool, nBlocks)

	// Cycle through the bits of the checksum, using them to create the design.
	const sumOffset uint = 24 // already used first 24 bits for colour
	for i := 0; i < nBlocks; i += 1 {
		blocks[i] = (h.sum >> ((sumOffset + uint(i)) % sumSize) & 1) == 1
	}
	design := parseSymmetric(designSize, blocks)
	return &designedHashed{&designed{design, colour}, h}
}

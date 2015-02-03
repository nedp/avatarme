package identicon

import (
	"github.com/dchest/siphash"
)

const (
	sumSize = 64 // size of the `sum` member of `Identicon`
)

const (
	designNone = -1 + iota
	designSymmetric
)

var hashKey = []byte{
	0x11, 0xBB, 0x22, 0xAA,
	0x33, 0x00, 0xEE, 0x66,
	0x99, 0x44, 0x77, 0x88,
	0xCC, 0xFF, 0x55, 0xDD,
}

type Identicon interface {
	Hash(info []byte) Hashed
}

type identicon struct {
	info       []byte
}

type hashedIdenticon struct {
	*hashed
	*identicon
}

type designedHashedIdenticon struct {
	*designed
	*hashed
	*identicon
}

func FromInfo(info []byte) *identicon {
	return &identicon{info}
}

// Generates a hash based on the identicon's info.
func (id *identicon) Hash(info []byte) *hashedIdenticon {
	hash := siphash.New(hashKey)
	hash.Write(info)
	hashed := FromHash(hash)
	return &hashedIdenticon{hashed, id}
}

func (id *hashedIdenticon) Hash(info []byte) *hashedIdenticon {
	return id
}

func (id *designedHashedIdenticon) Hash(info []byte) *designedHashedIdenticon {
	return id
}

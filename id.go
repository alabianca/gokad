package gokad

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type NodeID interface {
	GetBitAt(index uint) int
}

// SIZE describes how many bytes in an id
const SIZE = 20

// BITS describes how many bits in an id
const BITS = 160

// ID represents a NodeID of a node in the DHT
type ID struct {
	buf []byte
}

func (id ID) String() string {
	return fmt.Sprintf("%x", id.buf)
}

func (id *ID) Bytes() []byte {
	cpy := make([]byte, len(id.buf))
	copy(cpy, id.buf)
	return cpy
}

// Equal returns true if other is equal to the id
func (id *ID) Equal(other *ID) bool {
	for i := 0; i < SIZE; i++ {
		if id.buf[i] != other.buf[i] {
			return false
		}
	}

	return true
}

// DistanceTo calculates the 'distance' to other based on XOR metric
func (id *ID) DistanceTo(other *ID) *Distance {
	res := make([]byte, 20)

	for i := 0; i < SIZE; i++ {
		res[i] = id.buf[i] ^ other.buf[i]
	}

	return &Distance{res}
}

// CompareDistanceTo returns 0 if first and second are equally far away
// returns 1 if first is closer and return -1 if second is closer
func (id *ID) CompareDistanceTo(first *ID, second *ID) int {
	for i := 0; i < SIZE; i++ {
		b1 := id.buf[i] ^ first.buf[i]
		b2 := id.buf[i] ^ second.buf[i]

		if b1 < b2 {
			return 1
		}

		if b2 < b1 {
			return -1
		}
	}

	return 0
}

// GetBitAt returns the bit at the specified index
// If index >= BITS (160) the last bit is returned
func (id *ID) GetBitAt(index uint) int {
	if index >= BITS {
		index = 159
	}

	// find in what byte index the index falls
	bufferIndex := (index / 8) | 0
	mask := (1 << (7 - (index % 8)))
	bit := uint(id.buf[bufferIndex]) & uint(mask)

	if bit > 0 {
		return 1
	}

	return 0

}

// From returns an ID based on the provided hex id string
func From(hexID string) (*ID, error) {

	src := []byte(hexID)
	dst := make([]byte, hex.DecodedLen(len(src)))

	n, err := hex.Decode(dst, src)

	if err != nil {
		return nil, err
	}

	id := &ID{dst[:n]}

	return id, nil
}

// GenerateRandomID generates a random ID of length SIZE (20)
func GenerateRandomID() *ID {
	id := make([]byte, 20)
	rand.Read(id)

	return &ID{id}

}

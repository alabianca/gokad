package gokad

import "encoding/hex"

type Distance []byte

func (d Distance) String() string {
	return hex.EncodeToString(d)
}

func (d Distance) GetBitAt(index uint) int {
	if index >= BITS {
		index = 159
	}

	// find in what byte index the index falls
	byteIndex := (index / 8) | 0
	mask := (1 << (7 - (index % 8)))
	bit := uint(d[byteIndex]) & uint(mask)

	if bit > 0 {
		return 1
	}

	return 0
}

func (d Distance) Equal(other Distance) bool {
	for i := 0; i < SIZE; i++ {
		if d[i] != other[i] {
			return false
		}
	}

	return true
}

// type Distance struct {
// 	buf []byte
// }

// // GetBitAt returns the bit at the specified index
// // If index >= BITS (160) the last bit is returned
// func (d *Distance) GetBitAt(index uint) int {
// 	if index >= BITS {
// 		index = 159
// 	}

// 	// find in what byte index the index falls
// 	bufferIndex := (index / 8) | 0
// 	mask := (1 << (7 - (index % 8)))
// 	bit := uint(d.buf[bufferIndex]) & uint(mask)

// 	if bit > 0 {
// 		return 1
// 	}

// 	return 0

// }

// func (d Distance) String() string {
// 	return fmt.Sprintf("%x", d.buf)
// }

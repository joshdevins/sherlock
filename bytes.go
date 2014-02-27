package main

import (
	"log"
	"math"
)

const (
	BitsPerByte = 8
)

// Determines the bit-wise Hamming distance between two bytes.
func hammingDistance(left byte, right byte) int {
	if left == right {
		return 0
	}

	// find bits that differ by doing an XOR
	xor := left ^ right

	// count bitwise differences by taking the XOR result and AND with a bitmask
	diff := 0
	for i := 0; i < 8; i++ {
		bitmask := byte(math.Pow(2, float64(i)))
		if xor&bitmask > 0 {
			diff++
		}
	}

	return diff
}

// Flips a single bit within a byte.
func flipBit(b byte, i int) byte {
	if i < 0 || i >= BitsPerByte {
		log.Fatalf("Can not flip a bit in a position that does not exist: %d")
	}

	b ^= (1 << uint(BitsPerByte-1-i))

	return b
}

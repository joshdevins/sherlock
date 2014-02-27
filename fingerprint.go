package main

import (
	"fmt"
	"log"
	"math"
)

const (
	BitsPerByte             = 8
	SubFingerprintSizeBytes = 4                                     // 32-bits
	SubFingerprintSizeBits  = SubFingerprintSizeBytes * BitsPerByte // 32-bits
	FingerprintBlockSize    = 256                                   // except during testing, so this is a hint
)

type sub_fingerprint [SubFingerprintSizeBytes]byte

// size is determined at runtime (flexible to testing)
type fingerprint_block []sub_fingerprint

type fingerprint struct {
	id   string
	sfps []sub_fingerprint
}

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

// Determines the bit-wise Hamming distance from the sub-fingerprint to any
// other sub-fingerprint.
func (left *sub_fingerprint) hammingDistanceTo(right sub_fingerprint) int {
	distance := 0

	for i, leftByte := range left {
		rightByte := right[i]
		if leftByte != rightByte {
			distance += hammingDistance(leftByte, rightByte)
		}
	}

	return distance
}

// Creates a copy of the sub-fingerprint and flips a single bit.
func (sfp *sub_fingerprint) flipBit(i int) sub_fingerprint {
	if i < 0 || i >= SubFingerprintSizeBits {
		log.Fatalf("Can not flip a bit in a position that does not exist: %d")
	}

	// copy the underlying sub-fingerprint
	flipped := sub_fingerprint{}
	for i, v := range sfp {
		flipped[i] = v
	}

	// find the new indices
	byteIndex := i / 8
	bitIndex := i % 8

	// flip the bit in the byte
	flipped[byteIndex] = flipBit(flipped[byteIndex], bitIndex)

	return flipped
}

// For every bit of the sub-fingerprint, create a copy of the original and flip
// a single bit. This produces a slice of sub-fingerprints that are all Hamming
// distance of one (1) from the original sub-fingerprint.
func (sfp *sub_fingerprint) flipAllBits() []sub_fingerprint {
	flipped := make([]sub_fingerprint, SubFingerprintSizeBits)

	for i := 0; i < SubFingerprintSizeBits; i++ {
		flipped[i] = sfp.flipBit(i)
	}

	return flipped
}

// Flips the bits of a sub-fingerprint to generate all permutations that are
// equal to or less than a Hamming distance of `n`. Note that this is sequential
// and the algorithm is O(bits^n) so this can be slow for larger values of `n`.
// This algorithm could be optimised when necessary.
func (sfp *sub_fingerprint) flipAllBitsUntil(n int) ([]sub_fingerprint, error) {
	if n < 1 {
		err := fmt.Errorf("Target Hamming distance must be greater than or equal to 1: %d", n)
		return make([]sub_fingerprint, 0), err
	}

	set := make(
		map[sub_fingerprint]bool,
		int(math.Pow(
			float64(SubFingerprintSizeBits),
			float64(n),
		)),
	)
	set[*sfp] = true

	for i := 1; i <= n; i++ {
		for os, _ := range set {
			for _, fs := range os.flipAllBits() {
				set[fs] = true
			}
		}
	}

	// set into slice
	flipped := make([]sub_fingerprint, len(set))
	i := 0
	for k, _ := range set {
		flipped[i] = k
		i++
	}

	return flipped, nil
}

// Calculates the bit error rate from the fingerprint block to any other
// fingerprint block.
func (left *fingerprint_block) bitErrorRateWith(right fingerprint_block) (float32, error) {
	if len(*left) != len(right) {
		err := fmt.Errorf(
			"Fingerprint block to compare with was of size %d, but %d was expected",
			len(right),
			len(*left),
		)
		return 0.0, err
	}

	totalHammingDistance := 0
	for i, leftSfp := range *left {
		totalHammingDistance += leftSfp.hammingDistanceTo(right[i])
	}

	numBits := len(*left) * SubFingerprintSizeBits

	return float32(totalHammingDistance) / float32(numBits), nil
}

// Extract a fingerprint block from the fingerprint given the starting
// position and the final size of the fingerprint block. If the block would
// overflow off the end of the fingerprint, an error is returned instead.
func (fp *fingerprint) extractFingerprintBlock(start int, size int) (fingerprint_block, error) {
	if start < 0 || start >= int(len(fp.sfps)) {
		err := fmt.Errorf(
			"Start %d is not within the bounds of the fingerprint of size %d",
			start,
			len(fp.sfps),
		)
		return nil, err
	}

	if size < 1 || size > len(fp.sfps) {
		err := fmt.Errorf(
			"Size %d is not within the bounds of the fingerprint of size %d",
			size,
			len(fp.sfps),
		)
		return nil, err
	}

	end := start + size - 1 // end index is inclusive

	// check for overflow -- block ends past the end of the fingerprint
	if end >= len(fp.sfps) {
		err := fmt.Errorf(
			"End position %d (start position %d plus size %d) would overflow off the end of the fingerprint of size %d",
			end, start, size, len(fp.sfps),
		)
		return nil, err
	}

	fpb := fp.sfps[start : end+1] // end index is exclusive in this syntax

	return fpb, nil
}

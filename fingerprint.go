package main

import (
	"fmt"
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
func flipBit(b byte, i uint) (byte, error) {
	if i >= BitsPerByte {
		err := fmt.Errorf("Can not flip a bit in a position that does not exist: %d")
		return b, err // return unmodified byte
	}

	b ^= (1 << uint(BitsPerByte-1-i))

	return b, nil
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
func (sfp *sub_fingerprint) flipBit(i uint) (sub_fingerprint, error) {
	if i >= SubFingerprintSizeBits {
		err := fmt.Errorf("Can not flip a bit in a position that does not exist: %d")
		return *sfp, err // return unmodified sub-fingerprint
	}

	// copy the underlying sub-fingerprint
	flipped := sub_fingerprint{}
	for i, v := range sfp {
		flipped[i] = v
	}

	// find the new indices
	byteIndex := i / 8
	bitIndex := uint(i % 8)

	// flip the bit in the byte
	flippedByte, err := flipBit(flipped[byteIndex], bitIndex)
	if err != nil {
		return *sfp, err // return unmodified sub-fingerprint
	}

	flipped[byteIndex] = flippedByte
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
	if start < 0 || start >= len(fp.sfps) {
		err := fmt.Errorf(
			"Start %d is not within the bounds of the fingerprint of size %d",
			start, len(fp.sfps),
		)
		return nil, err
	}

	if size < 0 || size > len(fp.sfps) {
		err := fmt.Errorf(
			"Size %d is not within the bounds of the fingerprint of size %d",
			size, len(fp.sfps),
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

package main

import (
	"fmt"
	"math"
)

const (
	SubFingerprintSizeBytes = 4                           // 32-bits
	SubFingerprintSizeBits  = SubFingerprintSizeBytes * 8 // 32-bits
	FingerprintBlockSize    = 256                         // except during testing, so this is a hint
)

type sub_fingerprint [SubFingerprintSizeBytes]byte

// size is determined at runtime (flexible to testing)
type fingerprint_block []sub_fingerprint

type fingerprint struct {
	id   string
	sfps []sub_fingerprint
}

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

package main

import "testing"

func TestByteHammingDistance(t *testing.T) {
	fixtures := []struct {
		left     byte
		right    byte
		expected int
	}{
		{1, 1, 0},
		{1, 2, 2},
		{2, 1, 2},     // symmetrical
		{255, 0, 8},   // 11111111 -> 00000000
		{255, 1, 7},   // 11111111 -> 00000001
		{255, 2, 7},   // 11111111 -> 00000010
		{255, 3, 6},   // 11111111 -> 00000011
		{255, 131, 5}, // 11111111 -> 10000011
		{46, 41, 3},   // 00101110 -> 00101001
		{236, 29, 5},  // 11101100 -> 00011101
	}

	for i, fixture := range fixtures {
		got := hammingDistance(fixture.left, fixture.right)
		if fixture.expected != got {
			t.Errorf("[%d] Expected %d but got %d", i, fixture.expected, got)
		}
	}
}

func TestByteFlipBit(t *testing.T) {
	fixtures := []struct {
		b        byte
		i        uint
		expected byte
	}{
		{0, 7, 1}, // index 7; 00000000 (0) -> 00000001 (1)
		{0, 6, 2}, // index 6; 00000000 (0) -> 00000010 (2)
		{0, 5, 4}, // index 6; 00000000 (0) -> 00000100 (4)
		{4, 5, 0}, // index 5; 00000100 (4) -> 00000000 (0)
	}

	for i, fixture := range fixtures {
		got, err := flipBit(fixture.b, fixture.i)
		if err != nil {
			t.Errorf("[%d] Unexpected error: %s", err)
		}
		if fixture.expected != got {
			t.Errorf("[%d] Expected %d but got %d", i, fixture.expected, got)
		}
	}
}

func TestSubFingerprintHammingDistance(t *testing.T) {
	fixtures := []struct {
		left     sub_fingerprint
		right    sub_fingerprint
		expected int
	}{
		{sub_fingerprint{0, 0, 0, 0}, sub_fingerprint{0, 0, 0, 0}, 0},
		{sub_fingerprint{0, 0, 0, 0}, sub_fingerprint{0, 0, 1, 0}, 1},
		{sub_fingerprint{0, 1, 1, 0}, sub_fingerprint{0, 0, 1, 1}, 2},
		{sub_fingerprint{1, 1, 1, 1}, sub_fingerprint{0, 0, 4, 0}, 5},
		{sub_fingerprint{0, 1, 1, 0}, sub_fingerprint{0, 0, 1, 1}, 2},
		{sub_fingerprint{0, 0, 1, 1}, sub_fingerprint{0, 1, 1, 0}, 2}, // symmetrical
	}

	for i, fixture := range fixtures {
		got := fixture.left.hammingDistanceTo(fixture.right)
		if fixture.expected != got {
			t.Errorf("[%d] Expected %d but got %d", i, fixture.expected, got)
		}
	}
}

func TestSubFingerprintBitFlip(t *testing.T) {
	fixtures := []struct {
		sfp      sub_fingerprint
		i        uint
		expected sub_fingerprint
	}{
		{sub_fingerprint{0, 0, 0, 0}, 0*8 + 7, sub_fingerprint{1, 0, 0, 0}}, // index 07; byte 0, index 7; 00000000 (0) -> 00000001 (1)
		{sub_fingerprint{0, 0, 0, 0}, 1*8 + 7, sub_fingerprint{0, 1, 0, 0}}, // index 15; byte 1, index 7; 00000000 (0) -> 00000001 (1)
		{sub_fingerprint{0, 0, 0, 0}, 2*8 + 5, sub_fingerprint{0, 0, 4, 0}}, // index 22; byte 2, index 5; 00000000 (0) -> 00000100 (4)
	}

	for i, fixture := range fixtures {
		got, err := fixture.sfp.flipBit(fixture.i)
		if err != nil {
			t.Errorf("[%d] Unexpected error: %s", err)
		}
		if fixture.expected != got {
			t.Errorf("[%d] Expected %d but got %d", i, fixture.expected, got)
		}
	}
}

func TestBitErrorRateWith(t *testing.T) {

	left := fingerprint_block{
		sub_fingerprint{0, 0, 0, 0},
		sub_fingerprint{0, 0, 1, 0},
		sub_fingerprint{0, 0, 9, 0},
		sub_fingerprint{0, 0, 0, 0},
	}

	fixtures := []struct {
		expected float32
		right    fingerprint_block
	}{
		{
			0.25,
			fingerprint_block{
				sub_fingerprint{0, 0, 0, 0},
				sub_fingerprint{0, 0, 1, 0},
				sub_fingerprint{0, 0, 9, 0},
				sub_fingerprint{255, 255, 255, 255},
			},
		},
		{
			0.50,
			fingerprint_block{
				sub_fingerprint{255, 255, 255, 255},
				sub_fingerprint{0, 0, 1, 0},
				sub_fingerprint{0, 0, 9, 0},
				sub_fingerprint{255, 255, 255, 255},
			},
		},
	}

	for i, fixture := range fixtures {
		got, err := left.bitErrorRateWith(fixture.right)
		if err != nil {
			t.Fatalf("[%d] BER failed when it should not have: %s", i, err)
		}

		if fixture.expected != got {
			t.Errorf("[%d] Expected BER %f but was %f", i, fixture.expected, got)
		}
	}
}

func TestExtractFingerprintBlock(t *testing.T) {

	fp := fingerprint{
		"0001",
		[]sub_fingerprint{
			sub_fingerprint{0, 0, 0, 0},
			sub_fingerprint{0, 0, 1, 0},
			sub_fingerprint{0, 0, 9, 0},
			sub_fingerprint{1, 1, 0, 0},
		},
	}

	fixtures := []struct {
		start    int
		size     int
		expected fingerprint_block
	}{
		{
			1,
			2,
			fingerprint_block{
				sub_fingerprint{0, 0, 1, 0},
				sub_fingerprint{0, 0, 9, 0},
			},
		},
		{
			0,
			4,
			fingerprint_block{
				sub_fingerprint{0, 0, 0, 0},
				sub_fingerprint{0, 0, 1, 0},
				sub_fingerprint{0, 0, 9, 0},
				sub_fingerprint{1, 1, 0, 0},
			},
		},
	}

	for i, fixture := range fixtures {
		got, err := fp.extractFingerprintBlock(fixture.start, fixture.size)
		if err != nil {
			t.Fatalf("[%d] Extracting fingerprint block failed when it should not have: %s", i, err)
		}

		if len(fixture.expected) != len(got) {
			t.Errorf("[%d] Expected fingerprint of size %d but was of size %d", i, len(fixture.expected), len(got))
		}

		for j, expected := range fixture.expected {
			if expected != got[j] {
				t.Errorf("[%d][%d] Expected sub-fingerprint %v but was %v", i, j, expected, got[j])
			}
		}
	}
}

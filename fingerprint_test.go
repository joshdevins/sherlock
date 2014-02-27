package main

import "testing"

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

func TestSubFingerprintFlipBit(t *testing.T) {
	fixtures := []struct {
		sfp      sub_fingerprint
		i        int
		expected sub_fingerprint
	}{
		{sub_fingerprint{0, 0, 0, 0}, 0*8 + 7, sub_fingerprint{1, 0, 0, 0}}, // index 07; byte 0, index 7; 00000000 (0) -> 00000001 (1)
		{sub_fingerprint{0, 0, 0, 0}, 1*8 + 7, sub_fingerprint{0, 1, 0, 0}}, // index 15; byte 1, index 7; 00000000 (0) -> 00000001 (1)
		{sub_fingerprint{0, 0, 0, 0}, 2*8 + 5, sub_fingerprint{0, 0, 4, 0}}, // index 22; byte 2, index 5; 00000000 (0) -> 00000100 (4)
	}

	for i, fixture := range fixtures {
		got := fixture.sfp.flipBit(fixture.i)
		if fixture.expected != got {
			t.Errorf("[%d] Expected %d but got %d", i, fixture.expected, got)
		}
	}
}

func TestSubFingerprintFlipAllBits(t *testing.T) {
	fixture := sub_fingerprint{0, 0, 0, 0}
	expectations := []struct {
		index int
		sfp   sub_fingerprint
	}{
		{0*8 + 7, sub_fingerprint{1, 0, 0, 0}}, // index 07; byte 0, index 7; 00000000 (0) -> 00000001 (1)
		{1*8 + 7, sub_fingerprint{0, 1, 0, 0}}, // index 15; byte 1, index 7; 00000000 (0) -> 00000001 (1)
		{2*8 + 5, sub_fingerprint{0, 0, 4, 0}}, // index 22; byte 2, index 5; 00000000 (0) -> 00000100 (4)
	}

	flipped := fixture.flipAllBits()
	for i, e := range expectations {
		got := flipped[e.index]
		if e.sfp != got {
			t.Errorf("[%d] Expected %d but got %d", i, e, got)
		}
	}
}

func TestFingerprintBlockBitErrorRateWith(t *testing.T) {

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

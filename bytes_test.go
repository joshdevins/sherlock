package main

import "testing"

func TestHammingDistance(t *testing.T) {
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

func TestFlipBit(t *testing.T) {
	fixtures := []struct {
		b        byte
		i        int
		expected byte
	}{
		{0, 7, 1}, // index 7; 00000000 (0) -> 00000001 (1)
		{0, 6, 2}, // index 6; 00000000 (0) -> 00000010 (2)
		{0, 5, 4}, // index 6; 00000000 (0) -> 00000100 (4)
		{4, 5, 0}, // index 5; 00000100 (4) -> 00000000 (0)
	}

	for i, fixture := range fixtures {
		got := flipBit(fixture.b, fixture.i)
		if fixture.expected != got {
			t.Errorf("[%d] Expected %d but got %d", i, fixture.expected, got)
		}
	}
}

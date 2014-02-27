package main

import "testing"

func TestBuildIndex(t *testing.T) {
	corpus := []fingerprint{
		fingerprint{
			"0001",
			[]sub_fingerprint{
				sub_fingerprint{0, 0, 0, 0},
				sub_fingerprint{0, 0, 1, 0},
				sub_fingerprint{0, 0, 9, 0},
				sub_fingerprint{1, 8, 0, 0},
			},
		},
		fingerprint{
			"0002",
			[]sub_fingerprint{
				sub_fingerprint{0, 0, 0, 0},
				sub_fingerprint{0, 0, 1, 0},
				sub_fingerprint{0, 0, 9, 0},
				sub_fingerprint{1, 8, 0, 0},
			},
		},
		fingerprint{
			"0003",
			[]sub_fingerprint{
				sub_fingerprint{1, 0, 0, 0},
				sub_fingerprint{0, 0, 2, 0},
				sub_fingerprint{0, 7, 9, 0},
				sub_fingerprint{1, 8, 0, 1},
			},
		},
	}

	idx := buildIndex(corpus)

	// size of the index == total number of unique sub-fingerprints in the corpus
	if expected, got := 8, len(idx); expected != got {
		t.Errorf("Expected %d but got %d", expected, got)
	}

	// should have posting list with sound one and two
	pl := idx[sub_fingerprint{0, 0, 1, 0}]
	if expected, got := 2, len(pl); expected != got {
		t.Errorf("Expected %d but got %d", expected, got)
	}

	expectedPl := posting_list{
		posting{&corpus[0], 1},
		posting{&corpus[1], 1},
	}
	for i, _ := range pl {
		expectedP := expectedPl[i]
		p := pl[i]

		if expected, got := (*expectedP.fp).id, (*p.fp).id; expected != got {
			t.Errorf(
				"Expected posting at index %d did not have fingerprint with ID %s but was %s",
				i,
				expected,
				got,
			)
		}

		if expected, got := expectedP.offset, p.offset; expected != got {
			t.Errorf(
				"Expected posting at index %d did not have offset %d but was %d",
				i,
				expected,
				got,
			)
		}
	}
}

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
	postingsList := idx[sub_fingerprint{0, 0, 1, 0}]
	if expected, got := 2, len(postingsList); expected != got {
		t.Errorf("Expected %d but got %d", expected, got)
	}

	expectedPostingsList := postings_list{
		posting{&corpus[0], 1},
		posting{&corpus[1], 1},
	}
	for i, _ := range postingsList {
		expectedPosting := expectedPostingsList[i]
		posting := postingsList[i]

		if expected, got := (*expectedPosting.fp).id, (*posting.fp).id; expected != got {
			t.Errorf(
				"Expected posting at index %d did not have fingerprint with ID %s but was %s",
				i,
				expected,
				got,
			)
		}

		if expected, got := expectedPosting.offset, posting.offset; expected != got {
			t.Errorf(
				"Expected posting at index %d did not have offset %d but was %d",
				i,
				expected,
				got,
			)
		}
	}
}

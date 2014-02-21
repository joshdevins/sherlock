package main

import "testing"

func buildTestCorpus() []fingerprint {
	return []fingerprint{
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
}

func TestDeduplicateCandidates(t *testing.T) {
	corpus := buildTestCorpus()

	// three duplicates in candidates (0001, 1)
	candidates := []candidate{
		candidate{&corpus[0], 1},
		candidate{&corpus[0], 1},
		candidate{&corpus[1], 2},
		candidate{&corpus[0], 1},
		candidate{&corpus[2], 1},
		candidate{&corpus[2], 2},
		candidate{&corpus[2], 3},
	}

	deduped := deduplicateCandidates(candidates)

	if expected, got := len(candidates)-2, len(deduped); expected != got {
		t.Errorf("Expected %d but got %d", expected, got)
	}
}

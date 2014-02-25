package main

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

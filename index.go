package main

type posting struct {
	fp     *fingerprint
	offset int
}

type postings_list []posting

type index map[sub_fingerprint]postings_list

func buildIndex(corpus []fingerprint) index {
	idx := make(index)

	for i, fp := range corpus {
		for offset, sfp := range fp.sfps {
			posting := posting{&corpus[i], offset}

			// add posting to postings list for given sub-fingerprint
			postingsList, exists := idx[sfp]
			if exists {
				idx[sfp] = append(postingsList, posting)
			} else {
				postingsList = make(postings_list, 1)
				postingsList[0] = posting
				idx[sfp] = postingsList
			}
		}
	}

	return idx
}

package main

type posting struct {
	fp     *fingerprint
	offset int
}

type posting_list []posting

type index map[sub_fingerprint]posting_list

// Builds an index from a fingerprint corpus. Conserves memory by simply
// building against pointers to fingerprints in the corpus directly.
func buildIndex(corpus []fingerprint) index {
	idx := make(index)

	for i, fp := range corpus {
		for offset, sfp := range fp.sfps {

			// need to dereference the actual fp, can't just use `&fp`
			posting := posting{&corpus[i], offset}

			// add posting to posting list for given sub-fingerprint
			pl, exists := idx[sfp]
			if exists {
				idx[sfp] = append(pl, posting) // existing posting list
			} else {
				idx[sfp] = posting_list{posting} // new posting list
			}
		}
	}

	return idx
}

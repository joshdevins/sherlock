package main

type candidate struct {
	fp     *fingerprint
	offset int
}

type approximate_search_strategy func(sfp sub_fingerprint) ([]sub_fingerprint, error)

func noopApproximateSearchStrategy() approximate_search_strategy {
	return func(sfp sub_fingerprint) ([]sub_fingerprint, error) {
		return make([]sub_fingerprint, 0), nil
	}
}

func flipAllApproximateSearchStrategy(n uint) approximate_search_strategy {
	return func(sfp sub_fingerprint) ([]sub_fingerprint, error) {
		return sfp.flipAllBitsUntil(n)
	}
}

// Given a query fingerprint block, find candidates based on one or more of the
// sub-fingerprints. The candidates are not filtered in any way, for example on
// BER. This will always do an exact match search on the sub-fingerprints in the
// query fingerprint block, however you can optionally pass a strategy for
// approximate sub-fingerprint searching. This is usually a bit-flipping
// algorithm.
func searchByFingerprintBlock(
	queryFpb fingerprint_block,
	approxSearchStrategy approximate_search_strategy,
	idx index) ([]candidate, error) {

	candidates := make(map[candidate]bool)

	// find exact matches for sub-fingerprints in the fingerint block
	// BER threshold filtering will happen after generating all candidates (could
	// be optimised later)
	for queryOffset, querySfp := range queryFpb {
		newCandidates := searchBySubFingerprint(
			querySfp,
			queryOffset,
			idx,
		)

		if len(newCandidates) > 0 {
			for _, c := range newCandidates {
				candidates[c] = true
			}
		}
	}

	// try approximate searching, if a strategy was provided
	if approxSearchStrategy != nil {
		for queryOffset, querySfp := range queryFpb {
			approxQuerySfps, err := approxSearchStrategy(querySfp)
			if err != nil {
				return make([]candidate, 0), err
			}

			for _, approxQuerySfp := range approxQuerySfps {
				newCandidates := searchBySubFingerprint(
					approxQuerySfp,
					queryOffset,
					idx,
				)

				if len(newCandidates) > 0 {
					for _, c := range newCandidates {
						candidates[c] = true
					}
				}
			}
		}
	}

	// set to slice
	cslice := make([]candidate, len(candidates))
	i := 0
	for k, _ := range candidates {
		cslice[i] = k
		i++
	}

	return cslice, nil
}

// Given a sub-fingerprint and the offset of that sub-fingerprint in the query
// fingerprint block, find an exact match of the sub-fingerprint in the index
// and build a fingerprint block for each posting. The candidate fingerprint
// block is created such that the position in the fingerprint block of both the
// query and candidate sub-fingerprint are the same.
func searchBySubFingerprint(
	querySfp sub_fingerprint,
	queryOffset int,
	idx index) []candidate {

	postings, found := idx[querySfp]
	if !found {
		return make([]candidate, 0)
	}

	candidates := make([]candidate, len(postings))
	for i, posting := range postings {
		start := posting.offset - queryOffset
		candidates[i] = candidate{
			posting.fp,
			start,
		}
	}

	return candidates
}

package main

type candidate struct {
	fp     *fingerprint
	offset int
}

func deduplicateCandidates(candidates []candidate) []candidate {
	m := make(map[candidate]bool, 0)
	for _, c := range candidates {
		m[c] = true
	}

	deduped := make([]candidate, len(m))
	i := 0
	for k, v := range m {
		if v {
			deduped[i] = k
			i++
		}
	}

	return deduped
}

type approximate_search_strategy func(sfp sub_fingerprint) ([]sub_fingerprint, error)

func noopApproximateSearchStrategy() approximate_search_strategy {
	return func(sfp sub_fingerprint) ([]sub_fingerprint, error) {
		return make([]sub_fingerprint, 0), nil
	}
}

func flippingApproximateSearchStrategy(n int) approximate_search_strategy {
	return func(sfp sub_fingerprint) ([]sub_fingerprint, error) {
		flipped := make([]sub_fingerprint, SubFingerprintSizeBits)

		for i := 0; i < SubFingerprintSizeBits; i++ {
			flippedSfp, err := sfp.flipBit(uint(i))
			if err != nil {
				return make([]sub_fingerprint, 0), err
			}
			flipped[i] = flippedSfp
		}

		return flipped, nil
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

	candidates := make([]candidate, 0)

	// find exact matches for sub-fingerprints in the fingerint block
	// BER threshold filtering will happen after generating all candidates (could
	// be optimized later)
	for queryOffset, querySfp := range queryFpb {
		newCandidates := searchBySubFingerprint(
			querySfp,
			queryOffset,
			idx,
		)

		if len(newCandidates) > 0 {
			candidates = append(candidates, newCandidates...)
		}
	}

	// try approximate searching, if a strategy was provided
	// this will do a depth-first search but since we collect all candidates,
	// this should not matter
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
					candidates = append(candidates, newCandidates...)
				}
			}
		}
	}

	return deduplicateCandidates(candidates), nil
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

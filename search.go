package main

type candidate_fingerprint_block struct {
	fp     *fingerprint
	offset int
}

// TODO
func deduplicateCandidateFingerprintBlocks(candidates []candidate_fingerprint_block) []candidate_fingerprint_block {
	return candidates
}

type approximate_search_strategy func(sfp sub_fingerprint) []sub_fingerprint

func noopApproximateSearchStrategy() approximate_search_strategy {
	return func(sfp sub_fingerprint) []sub_fingerprint {
		return make([]sub_fingerprint, 0)
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
	idx index) []candidate_fingerprint_block {

	candidates := make([]candidate_fingerprint_block, 0)

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
			approxQuerySfps := approxSearchStrategy(querySfp)

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

	// deduplicate any candidates and return
	return deduplicateCandidateFingerprintBlocks(candidates)
}

// Given a sub-fingerprint and the offset of that sub-fingerprint in the query
// fingerprint block, find an exact match of the sub-fingerprint in the index
// and build a fingerprint block for each posting. The candidate fingerprint
// block is created such that the position in the fingerprint block of both the
// query and candidate sub-fingerprint are the same.
func searchBySubFingerprint(
	querySfp sub_fingerprint,
	queryOffset int,
	idx index) []candidate_fingerprint_block {

	postings, found := idx[querySfp]
	if !found {
		return make([]candidate_fingerprint_block, 0)
	}

	candidates := make([]candidate_fingerprint_block, len(postings))
	for i, posting := range postings {
		start := posting.offset - queryOffset
		candidates[i] = candidate_fingerprint_block{
			posting.fp,
			start,
		}
	}

	return candidates
}

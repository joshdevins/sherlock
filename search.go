package main

type candidate_fingerprint_block struct {
	fp     *fingerprint
	offset int
	fpb    fingerprint_block
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

	// find exact matches for sub-fingerprints in the fingerint block
	// note that only one needs to match in order to consider it as a candidate
	// BER threshold filtering will happen after generating all candidates (could be optimized later)
	for queryOffset, querySfp := range queryFpb {
		candidates := searchBySubFingerprint(
			querySfp,
			queryOffset,
			idx,
			len(queryFpb),
		)

		if len(candidates) > 0 {
			return candidates
		}
	}

	// none of the exact sub-fingerprints were found
	// try approximate searching, if a strategy was provided
	if approxSearchStrategy != nil {
		for queryOffset, querySfp := range queryFpb {
			approxQuerySfps := approxSearchStrategy(querySfp)

			for _, approxQuerySfp := range approxQuerySfps {
				candidates := searchBySubFingerprint(
					approxQuerySfp,
					queryOffset,
					idx,
					len(queryFpb),
				)

				if len(candidates) > 0 {
					return candidates
				}
			}
		}
	}

	// ok, seriously, still nothing was found, even after approximate search
	return make([]candidate_fingerprint_block, 0)
}

// Given a sub-fingerprint and the offset of that sub-fingerprint in the query
// fingerprint block, find an exact match of the sub-fingerprint in the index
// and build a fingerprint block for each posting. The candidate fingerprint
// block is created such that the position in the fingerprint block of both the
// query and candidate sub-fingerprint are the same.
func searchBySubFingerprint(
	querySfp sub_fingerprint,
	queryOffset int,
	idx index,
	fingerprintBlockSize int) []candidate_fingerprint_block {

	postings, found := idx[querySfp]
	if !found {
		return make([]candidate_fingerprint_block, 0)
	}

	candidates := make([]candidate_fingerprint_block, 0)
	for _, posting := range postings {
		start := posting.offset - queryOffset
		fpb, err := posting.fp.extractFingerprintBlock(start, fingerprintBlockSize)
		if err != nil {
			candidate := candidate_fingerprint_block{
				posting.fp,
				start,
				fpb,
			}
			candidates = append(candidates, candidate)
		}
	}

	return candidates
}

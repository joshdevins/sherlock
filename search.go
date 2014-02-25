package main

import "fmt"

// Candidate fingerprint block within a fingerprint. The offset here is the
// start of the block. The end of the block is determined by the offset plus the
// size of the blocks being used.
type candidate struct {
	fp     *fingerprint
	offset uint
}

func (c *candidate) extractFingerprintBlock(size uint) (fingerprint_block, error) {
	return c.fp.extractFingerprintBlock(c.offset, size)
}

func addCandidatesToSet(s []candidate, m map[candidate]bool) {
	for _, c := range s {
		m[c] = true
	}
}

func candidateSetToSlice(m map[candidate]bool) []candidate {
	s := make([]candidate, len(m))
	i := 0
	for k, _ := range m {
		s[i] = k
		i++
	}

	return s
}

func filterCandidatesByBER(
	queryFpb fingerprint_block,
	candidates []candidate,
	ber float32) []candidate {

	var filtered []candidate
	for _, candidate := range candidates {
		// FIXME: errors swallowed
		candidateFpb, _ := candidate.extractFingerprintBlock(uint(len(queryFpb)))
		actualBer, _ := queryFpb.bitErrorRateWith(candidateFpb)

		if actualBer <= ber {
			filtered = append(filtered, candidate)
		}
	}

	return filtered
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

// Given a sub-fingerprint and the offset of that sub-fingerprint in the query
// fingerprint block, find an exact match of the sub-fingerprint in the index.
// The candidate fingerprint block is created such that the position in the
// fingerprint block of both the query and candidate sub-fingerprint are the
// same. Only this offset is saved since the block can be regenerated from this.
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
		candidates[i] = candidate{
			posting.fp,
			uint(posting.offset - queryOffset),
		}
	}

	return candidates
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

		addCandidatesToSet(newCandidates, candidates)
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

				addCandidatesToSet(newCandidates, candidates)
			}
		}
	}

	return candidateSetToSlice(candidates), nil
}

// Given a query fingerprint, find candidates based on a sliding window query
// fingerprint block. The step size of the sliding window and the block size
// must be specified. Note that it is possible to provide a step size that is
// greater than or equal to the block size. This results in sub-fingerprints
// being searched no more than once from the query fingerprint.
func searchByFingerprint(
	queryFp fingerprint,
	blockSize uint,
	stepSize uint,
	approxSearchStrategy approximate_search_strategy,
	ber float32,
	idx index) ([]candidate, error) {

	if blockSize < 1 {
		err := fmt.Errorf("Block size must be greater than or equal to one: %d", blockSize)
		return make([]candidate, 0), err
	}

	if stepSize < 1 {
		err := fmt.Errorf("Step size must be greater than or equal to one: %d", stepSize)
		return make([]candidate, 0), err
	}

	if l := uint(len(queryFp.sfps)); l < blockSize {
		err := fmt.Errorf("Query fingerprint must be greater than or equal to a block (%d): %d", blockSize, l)
		return make([]candidate, 0), err
	}

	candidates := make(map[candidate]bool)

	// step through the fingerprint, taking steps as specified
	for offset := uint(0); offset+blockSize < uint(len(queryFp.sfps)); offset += stepSize {
		queryFpb, err := queryFp.extractFingerprintBlock(offset, blockSize)
		if err != nil {
			return make([]candidate, 0), err
		}

		newCandidates, err := searchByFingerprintBlock(queryFpb, approxSearchStrategy, idx)
		if err != nil {
			return make([]candidate, 0), err
		}
		addCandidatesToSet(filterCandidatesByBER(queryFpb, newCandidates, ber), candidates)
	}

	return candidateSetToSlice(candidates), nil
}

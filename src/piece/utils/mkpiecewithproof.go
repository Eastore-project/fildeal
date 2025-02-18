package pieceutils

import (
	"io"
	"math/bits"

	"github.com/filecoin-project/go-data-segment/datasegment"
	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/filecoin-project/go-state-types/abi"
)

func MakeDataSegmentPieceWithProof(subPieces []io.ReadSeeker) (*ProofResult, error) {
	readers := make([]io.Reader, 0)
	deals := make([]abi.PieceInfo, 0)

	for _, arg := range subPieces {
		readers = append(readers, arg)
		cp := new(commp.Calc)
		_, err := io.Copy(cp, arg)
		if err != nil {
			return nil, err
		}
		rawCommP, size, err := cp.Digest()
		if err != nil {
			return nil, err
		}
		_, err = arg.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
		c, err := commcid.DataCommitmentV1ToCID(rawCommP)
		if err != nil {
			return nil, err
		}
		subdeal := abi.PieceInfo{
			Size:     abi.PaddedPieceSize(size),
			PieceCID: c,
		}
		deals = append(deals, subdeal)
	}

	if len(deals) == 0 {
		return nil, nil
	}

	_, size, err := datasegment.ComputeDealPlacement(deals)
	if err != nil {
		return nil, err
	}

	overallSize := abi.PaddedPieceSize(size)
	// we need to make this the 'next' power of 2 in order to have space for the index
	next := 1 << (64 - bits.LeadingZeros64(uint64(overallSize+256)))

	agg, err := datasegment.NewAggregate(abi.PaddedPieceSize(next), deals)
	if err != nil {
		return nil, err
	}

	out, err := agg.AggregateObjectReader(readers)
	if err != nil {
		return nil, err
	}

	// Generate proofs for each piece
	inclProofs := make([]datasegment.InclusionProof, len(deals))
	ids := make([]uint64, len(deals))

	for i, piece := range deals {
		podsi, err := agg.ProofForPieceInfo(piece)
		if err != nil {
			return nil, err
		}
		ids[i] = uint64(i) // Using index as ID since we don't have offer IDs
		inclProofs[i] = *podsi
	}

	return &ProofResult{
		Reader:     out,
		InclProofs: inclProofs,
	}, nil
}

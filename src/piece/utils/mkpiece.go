package pieceutils

import (
	"bytes"
	"fmt"
	"io"
	"math/bits"

	"github.com/filecoin-project/go-data-segment/datasegment"
	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/filecoin-project/go-state-types/abi"
)

// ProofResult contains the reader for the data segment and inclusion proofs
type ProofResult struct {
	Reader     io.Reader
	InclProofs []datasegment.InclusionProof
}

func MakeDataSegmentPiece(subPieces []io.ReadSeeker) (io.Reader, error) {
	readers := make([]io.Reader, 0)
	deals := make([]abi.PieceInfo, 0)
	for _, arg := range subPieces {
		readers = append(readers, arg)
		cp := new(commp.Calc)
		_, err := io.Copy(cp, arg)
		if err != nil {
			return nil, fmt.Errorf("failed to copy data: %w", err)
		}
		rawCommP, size, err := cp.Digest()
		if err != nil {
			return nil, fmt.Errorf("failed to calculate digest: %w", err)
		}
		_, err = arg.Seek(0, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("failed to seek: %w", err)
		}
		c, err := commcid.DataCommitmentV1ToCID(rawCommP)
		if err != nil {
			return nil, fmt.Errorf("failed to create CID: %w", err)
		}
		subdeal := abi.PieceInfo{
			Size:     abi.PaddedPieceSize(size),
			PieceCID: c,
		}
		deals = append(deals, subdeal)
	}
	if len(deals) == 0 {
		return nil, fmt.Errorf("no deals provided")
	}

	_, size, err := datasegment.ComputeDealPlacement(deals)
	if err != nil {
		return nil, fmt.Errorf("failed to compute deal placement: %w", err)
	}

	overallSize := abi.PaddedPieceSize(size)
	// we need to make this the 'next' power of 2 in order to have space for the index
	next := 1 << (64 - bits.LeadingZeros64(uint64(overallSize+256)))

	a, err := datasegment.NewAggregate(abi.PaddedPieceSize(next), deals)
	if err != nil {
		return nil, fmt.Errorf("failed to create new aggregate: %w", err)
	}
	out, err := a.AggregateObjectReader(readers)
	if err != nil {
		return nil, fmt.Errorf("failed to create aggregate object reader: %w", err)
	}
	return out, nil
}

func ParseSegmentPieces(piece io.ReadSeeker) ([]io.Reader, error) {
	size, err := piece.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to end: %w", err)
	}
	offset := datasegment.DataSegmentIndexStartOffset(abi.UnpaddedPieceSize(size).Padded())
	_, err = piece.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to offset: %w", err)
	}
	index, err := datasegment.ParseDataSegmentIndex(piece)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data segment index: %w", err)
	}
	entries, err := index.ValidEntries()
	if err != nil {
		return nil, fmt.Errorf("failed to get valid entries: %w", err)
	}
	out := make([]io.Reader, 0)

	for _, e := range entries {
		buf := bytes.NewBuffer(nil)
		upoff := abi.PaddedPieceSize(e.Offset).Unpadded()
		_, err = piece.Seek(int64(upoff), io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("failed to seek to piece offset: %w", err)
		}
		upSize := abi.PaddedPieceSize(e.Size).Unpadded()
		_, err = io.CopyN(buf, piece, int64(upSize))
		if err != nil {
			return nil, fmt.Errorf("failed to copy piece data: %w", err)
		}
		out = append(out, buf)
	}
	return out, nil
}

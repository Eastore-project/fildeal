package index

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/filecoin-project/boost/extern/boostd-data/model"
	"github.com/filecoin-project/go-data-segment/datasegment"
	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/filecoin-project/go-state-types/abi"
	carv2 "github.com/ipld/go-car/v2"
	"golang.org/x/sync/errgroup"
)
const (
	MaxCachedReaders = 128
	// 20 MiB x 4 parallel deals is just 80MiB RAM overhead required
	PodsiBuffesrSize = 20e6
	// Concurrency is driven by the number of available cores. Set reasonable max and mins
	// to support multiple concurrent AddIndex operations
	PodsiMaxConcurrency = 32
	PodsiMinConcurrency = 4
)

func BoostIndex(filePath string) error { 
    r, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer r.Close()
	var cp commp.Calc
	io.Copy(&cp, r)
	rawCommP, size, err := cp.Digest()
	if err != nil {
		panic(err)
	}
	r.Seek(0, io.SeekStart)
	c, _ := commcid.DataCommitmentV1ToCID(rawCommP)

	fmt.Println("Unpadded piece size: ", size)
	fmt.Println("Padded piece size: ", abi.PaddedPieceSize(size).Unpadded())
	dsis := datasegment.DataSegmentIndexStartOffset(abi.PaddedPieceSize(size))

	// unnecessary, something about the way boost works
	if _, err = r.Seek(0, io.SeekEnd); err != nil {
		panic(err)
	}

	fmt.Printf("Seeking back to %d\n", dsis)
	if _, err := r.Seek(int64(dsis), io.SeekStart); err != nil {
		panic(err)
	}

	var readsCnt int32
	cr := &countingReader{
		Reader: r,
		cnt:    &readsCnt,
	}

	index, err := datasegment.ParseDataSegmentIndex(bufio.NewReaderSize(cr, 20e6))
	if err != nil {
		panic(err)
	}

	fmt.Println("CommP: ", c)

	jsonData, err := json.Marshal(index)
	if err != nil {
		panic(err)
	}
	fmt.Println("Found index data: " + string(jsonData))

	readsCnt = 0

	concurrency := len(index.Entries)

	chunkSize := len(index.Entries) / concurrency
	results := make([][]datasegment.SegmentDesc, concurrency)

	var eg errgroup.Group
	for i := 0; i < concurrency; i++ {
		i := i
		eg.Go(func() error {
			start := i * chunkSize
			end := start + chunkSize
			if i == concurrency-1 {
				end = len(index.Entries)
			}

			res, err := validateEntries(index.Entries[start:end])
			if err != nil {
				return err
			}

			results[i] = res

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		panic(err)
	}
	fmt.Println(len(index.Entries), "index entries found")
	validSegments := make([]datasegment.SegmentDesc, 0, len(index.Entries))
	for _, res := range results {
		validSegments = append(validSegments, res...)
	}

	if len(validSegments) == 0 {
		panic("no valid data segments found")
	}
	fmt.Println(len(validSegments), "valid entries found")

	for i, e := range validSegments {
		if err := e.Validate(); err != nil {
			fmt.Printf("Error validating entry %d: %s\n", i, err)
			// continue
		}

		segOffset := e.UnpaddedOffest()
		segSize := e.UnpaddedLength()

		lr := io.NewSectionReader(r, int64(segOffset), int64(segSize))

		// write the segment to args[1]_segment_<i>
		segmentPath := fmt.Sprintf("%s_segment_%d", filePath, i)
		sw, err := os.Create(segmentPath)
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(sw, lr); err != nil {
			panic(err)
		}
		if err := sw.Close(); err != nil {
			panic(err)
		}
		fmt.Printf("Segment #%d written to to %s\n", i, segmentPath)

		lr = io.NewSectionReader(r, int64(segOffset), int64(segSize))
		cr = &countingReader{
			Reader: lr,
			cnt:    &readsCnt,
		}
		recs := make([]model.Record, 0)

		subRecs, err := parseRecordsFromCar(bufio.NewReaderSize(cr, PodsiBuffesrSize))
		if err != nil {
			return fmt.Errorf("could not parse data segment #%d at offset %d: %w", len(recs), segOffset, err)
		}
		for i := range subRecs {
			subRecs[i].Offset += segOffset
		}
		recs = append(recs, subRecs...)

		// opts := []carv2.Option{carv2.ZeroLengthSectionAsEOF(true)}
		// blockReader, err := carv2.NewBlockReader(bufio.NewReaderSize(cr, 20e6), opts...)
		// if err != nil {
		// 	panic(e)
		// }

		// blockMetadata, err := blockReader.SkipNext()
		// for err == nil {
		// 	fmt.Printf("Segment #%d CAR Block: %s, Offset: %d, Size: %d\n", i, blockMetadata.Cid, blockMetadata.SourceOffset, blockMetadata.Size)
		// 	blockMetadata, err = blockReader.SkipNext()
		// }
		// if !errors.Is(err, io.EOF) {
		// 	fmt.Printf("Error reading blocks: %s\n", err)
		// }
	}
	fmt.Printf("Parsed PoDSI piece (with %d reads)\n", readsCnt)
	return nil
}

type countingReader struct {
	io.Reader

	cnt *int32
}

func (cr *countingReader) Read(p []byte) (n int, err error) {
	atomic.AddInt32(cr.cnt, 1)
	return cr.Reader.Read(p)
}

func validateEntries(entries []datasegment.SegmentDesc) ([]datasegment.SegmentDesc, error) {
	res := make([]datasegment.SegmentDesc, 0, len(entries))
	for i, e := range entries {

		if err := e.Validate(); err != nil {
			if errors.Is(err, datasegment.ErrValidation) {
				fmt.Printf("Error validating entry: %s\n", err)
				continue
			} else {
				return nil, fmt.Errorf("got unknown error for entry %d: %w", i, err)
			}
		}
		res = append(res, e)
	}
	return res, nil
}

func parseRecordsFromCar(reader io.Reader) ([]model.Record, error) {
	// Iterate over all the blocks in the piece to extract the index records
	recs := make([]model.Record, 0)
	opts := []carv2.Option{carv2.ZeroLengthSectionAsEOF(true)}
	blockReader, err := carv2.NewBlockReader(reader, opts...)
	if err != nil {
		return nil, fmt.Errorf("getting block reader over piece: %w", err)
	}

	blockMetadata, err := blockReader.SkipNext()
	for err == nil {
		recs = append(recs, model.Record{
			Cid: blockMetadata.Cid,
			OffsetSize: model.OffsetSize{
				Offset: blockMetadata.SourceOffset,
				Size:   blockMetadata.Size,
			},
		})

		blockMetadata, err = blockReader.SkipNext()
	}
	if !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("generating index for piece: %w", err)
	}
	return recs, nil
}
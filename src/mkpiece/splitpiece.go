package mkpiece

import (
	"fmt"
	"io"
	"os"

	"github.com/filecoin-project/go-data-segment/datasegment"
	"github.com/filecoin-project/go-state-types/abi"
)

func SplitPiece(filePath, outputDir string) error {
    file, err := os.Open(filePath)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    fi, err := file.Stat()
    if err != nil {
        return fmt.Errorf("failed to get file info: %w", err)
    }

    offset := datasegment.DataSegmentIndexStartOffset(abi.UnpaddedPieceSize(fi.Size()).Padded())
    file.Seek(int64(offset), io.SeekStart)
    index, err := datasegment.ParseDataSegmentIndex(file)
    if err != nil {
        return fmt.Errorf("failed to parse data segment index: %w", err)
    }
    entries, err := index.ValidEntries()
    if err != nil {
        return fmt.Errorf("failed to get valid entries: %w", err)
    }
    for _, e := range entries {
        file.Seek(0, io.SeekStart)
        strt := e.UnpaddedOffest()
        leng := e.UnpaddedLength()
        segment := io.NewSectionReader(file, int64(strt), int64(leng))
        seg, err := os.OpenFile(fmt.Sprintf("%s/%d.chunk", outputDir, strt), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
        if err != nil {
            return fmt.Errorf("failed to open segment file: %w", err)
        }
        n, err := io.Copy(seg, segment)
        if err != nil {
            return fmt.Errorf("failed to copy segment: %w", err)
        }
        if n != int64(leng) {
            return fmt.Errorf("didn't write enough: wrote %d, expected %d", n, leng)
        }
        seg.Close()
        fmt.Printf("Segment found: %d - %d\n", strt, strt+leng)
    }
    return nil
}
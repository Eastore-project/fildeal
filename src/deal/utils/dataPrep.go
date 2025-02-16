package dealutils

import (
	"fmt"
	"path/filepath"
)

type DataPrepResult struct {
	PieceCid   string
	PayloadCid string
	PieceSize  uint64
	CarSize    uint64
	LocalPath  string
	Hash       string
}

func PrepareData(inputPath, outDir string, buffer string, apiKey string) (*DataPrepResult, error) {
	output, err := ConvertToCar(inputPath, outDir, inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to car: %w", err)
	}

	result := &DataPrepResult{
		PieceCid:   output.PieceCid,
		PayloadCid: output.DataCid,
		PieceSize:  output.PieceSize,
		CarSize:    output.CarSize,
		LocalPath:  filepath.Join(outDir, fmt.Sprintf("%s.car", output.PieceCid)),
	}

	if buffer == "lighthouse" {
		lighthouseResp, err := UploadToLighthouse(result.LocalPath, apiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to upload to Lighthouse: %w", err)
		}
		result.Hash = lighthouseResp.Hash
	}

	return result, nil
}

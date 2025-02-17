package dealutils

import (
	"fmt"
	"path/filepath"

	"github.com/eastore-project/fildeal/src/buffer"
)

type DataPrepResult struct {
	PieceCid   string
	PayloadCid string
	PieceSize  uint64
	CarSize    uint64
	LocalPath  string
	BufferInfo *buffer.Response
}

func PrepareData(inputPath, outDir string, bufferConfig *buffer.Config) (*DataPrepResult, error) {
	output, err := ConvertToCar(inputPath, outDir, inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to car: %w", err)
	}

	carPath := filepath.Join(outDir, fmt.Sprintf("%s.car", output.PieceCid))
	result := &DataPrepResult{
		PieceCid:   output.PieceCid,
		PayloadCid: output.DataCid,
		PieceSize:  output.PieceSize,
		CarSize:    output.CarSize,
		LocalPath:  carPath,
	}

	var buf buffer.Buffer
	switch bufferConfig.Type {
	case "lighthouse":
		buf = buffer.NewLighthouseBuffer(bufferConfig.ApiKey, bufferConfig.BaseURL)
	default:
		buf = buffer.NewLocalBuffer() // No port needed for data prep
	}

	bufferResp, err := buf.Store(result.LocalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to store in buffer: %w", err)
	}
	result.BufferInfo = bufferResp

	return result, nil
}

package dealutils

import (
	"fmt"

	"github.com/eastore-project/fildeal/src/utils"
)

func ConvertToCar(path string, outputDir string, parentPath string) (*Result, error) {
	carParams := &CarParams{
		Input:     path,
		PieceSize: 0,
		OutDir:    outputDir,
		Parent:    parentPath,
		TmpDir:    "",
		Single:    true,
	}
	output, err := carParams.GenerateCarUtil()
	if err != nil {
		return nil, fmt.Errorf("failed to generate car file: %w", err)
	}

	carSize, err := utils.GetFileSize(fmt.Sprintf("%s/%s.car", outputDir, output.PieceCid))
	if err != nil {
		return nil, fmt.Errorf("failed to get car file size: %w", err)
	}

	output.CarSize = uint64(carSize)
	return &output, nil
}

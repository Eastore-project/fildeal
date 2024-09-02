package dealutils

import (
	"fildeal/src/types"
	"fmt"
)


func ConvertToCar(filePath string, outputDir string, parentPath string) ( *types.Result, error) {
	carParams := &CarParams{
		Input:     filePath,
		PieceSize: 0,
		OutDir:    outputDir,
		Parent:    parentPath,
		TmpDir:    "",
		Single:    true,
	}
	output, err := carParams.GenerateCar()
	if err != nil {
		return nil, fmt.Errorf("failed to generate car file: %w", err)
	}

	carSize, err := GetFileSize(fmt.Sprintf("%s/%s.car", outputDir, output.PieceCid))
	if err != nil {
		return nil, fmt.Errorf("failed to get car file size: %w", err)
	}

	output.CarSize = uint64(carSize)
	return  &output, nil
}


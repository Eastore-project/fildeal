package deal

import (
	configurations "fildeal/src/config"
	dealutils "fildeal/src/deal/utils"
	"fildeal/src/mkpiece"
	"fmt"
	"io"
	"os"
	"path/filepath"

	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/google/uuid"
)

func MakeDeal(inputFolder string, miner string) error {

// Convert each file of input folder to car file in dummy output folder
   files, err := os.ReadDir(inputFolder)
   if err != nil {
	return fmt.Errorf("failed to read input folder: %w", err)
   }
   
    // Define and create the output folder if it doesn't exist
    outputFolder := configurations.LoadConfigurations().GenerateCarPath
	// clear the output folder first if it exists
	err = os.RemoveAll(outputFolder)
	if err != nil {
		return fmt.Errorf("failed to clear generate car folder: %w", err)
	}

    err = os.MkdirAll(outputFolder, 0755)
    if err != nil {
        return fmt.Errorf("failed to create output folder: %w", err)
    }

   for i, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info: %w", err)
		}
		filePath := filepath.Join(inputFolder, fileInfo.Name())


		if fileInfo.IsDir() {
			// Handle directory case
			fmt.Printf("Processing directory: %s\n", fileInfo.Name())
			output, err := dealutils.ConvertToCar(filePath, outputFolder, inputFolder)
			if err != nil {
				return fmt.Errorf("failed to convert directory to CAR: %w", err)
			}
			// rename the directory to maintain order
			newFilePath := fmt.Sprintf("%s/%d_%s", outputFolder, i, output.PieceCid)
			err = os.Rename(fmt.Sprintf("%s/%s.car", outputFolder, output.PieceCid), newFilePath)
			if err != nil {
				return fmt.Errorf("failed to rename directory: %w", err)
			}
			fmt.Printf("Output: %v\n", output)
		} else {
			// Handle file case
			fmt.Printf("Processing file: %s\n", fileInfo.Name())
			output, err := dealutils.ConvertToCar(filePath, outputFolder, inputFolder)
			if err != nil {
				return fmt.Errorf("failed to convert file to CAR: %w", err)
			}
			// rename file to maintain order
			newFilePath := fmt.Sprintf("%s/%d_%s", outputFolder, i, output.PieceCid)
			err = os.Rename(fmt.Sprintf("%s/%s.car", outputFolder, output.PieceCid), newFilePath)
			if err != nil {
				return fmt.Errorf("failed to rename file: %w", err)
			}
			fmt.Printf("Output: %v\n", output)
		}
   	}
	// read each car file and make aggregate using mkpiece
	readers, err := dealutils.GetReaders(outputFolder)
    if err != nil {
        return err
    }
    defer func() {
        for _, r := range readers {
            r.(io.Closer).Close()
        }
    }()

	out := mkpiece.MakeDataSegmentPiece(readers)

	commpHasher := commp.Calc{}
	_, _ = io.CopyBuffer(&commpHasher, out, make([]byte, commpHasher.BlockSize()*128))
	commp, pieceSize, _ := commpHasher.Digest()
	pieceCid, err := commcid.PieceCommitmentV1ToCID(commp)
	if err != nil {
		return fmt.Errorf("failed to get piece CID: %w", err)
	}

	// Put aggregate in aggregateCar folder with uuid as name
	aggregateName := uuid.New().String()
	aggregateFolder := configurations.LoadConfigurations().AggregateCarPath
	err = os.MkdirAll(aggregateFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create aggregate folder: %w", err)
	}
	aggregatePath := fmt.Sprintf("%s%s.data", aggregateFolder, aggregateName)
	// copy out to aggretgatePath
	aggregateFile, err := os.Create(aggregatePath)
	if err != nil {
		return fmt.Errorf("failed to create aggregate file: %w", err)
	}
	defer aggregateFile.Close()

	// read each car file and make aggregate using mkpiece
	readers = make([]io.ReadSeeker, 0)

    aggreateReaders, err := dealutils.GetReaders(outputFolder)
    if err != nil {
        return err
    }
    defer func() {
        for _, r := range readers {
            r.(io.Closer).Close()
        }
    }()

	out = mkpiece.MakeDataSegmentPiece(aggreateReaders)

	_, err = io.Copy(aggregateFile, out)
	if err != nil {
		return fmt.Errorf("failed to copy aggregate: %w", err)
	}
	carSize, err := dealutils.GetFileSize(aggregatePath)
	if carSize == 0 || err != nil {
		return fmt.Errorf("failed to get aggregate file size: %w", err)
	}

	// Create deal with miner
	 err = dealutils.InitiateDeal(aggregateName, miner, pieceSize, pieceCid.String(), uint64(carSize))
	if err != nil {
		return fmt.Errorf("failed to initiate deal: %w", err)
	}
	return nil
}
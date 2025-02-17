package deal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/eastore-project/fildeal/src/buffer"
	utils "github.com/eastore-project/fildeal/src/deal/utils"
	"github.com/eastore-project/fildeal/src/mkpiece"

	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func MakePodsiDeal(ctx *cli.Context) error {
	inputFolder := ctx.String("input")
	miner := ctx.String("miner")
	bufferType := ctx.String("buffer")
	lighthouseApiKey := ctx.String("lighthouse-api-key")

	if bufferType == "lighthouse" && lighthouseApiKey == "" {
		return fmt.Errorf("lighthouse API key is required when using lighthouse buffer")
	}

	if ctx.Uint("duration") < 518400 || ctx.Uint("duration") > 1814400 {
		return fmt.Errorf("duration must be between 518400 (6 months) and 181440 (app. 3.5 years)")
	}
	// Process input folder
	files, err := os.ReadDir(inputFolder)
	if err != nil {
		return fmt.Errorf("failed to read input folder: %w", err)
	}

	// Use provided generateCarPath from context
	outputFolder := ctx.String("generate-car-path")
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
			fmt.Printf("Processing directory: %s\n", fileInfo.Name())
			output, err := utils.ConvertToCar(filePath, outputFolder, inputFolder)
			if err != nil {
				return fmt.Errorf("failed to convert directory to CAR: %w", err)
			}
			newFilePath := fmt.Sprintf("%s/%d_%s", outputFolder, i, output.PieceCid)
			err = os.Rename(fmt.Sprintf("%s/%s.car", outputFolder, output.PieceCid), newFilePath)
			if err != nil {
				return fmt.Errorf("failed to rename directory: %w", err)
			}
		} else {
			fmt.Printf("Processing file: %s\n", fileInfo.Name())
			output, err := utils.ConvertToCar(filePath, outputFolder, inputFolder)
			if err != nil {
				return fmt.Errorf("failed to convert file to CAR: %w", err)
			}
			newFilePath := fmt.Sprintf("%s/%d_%s", outputFolder, i, output.PieceCid)
			err = os.Rename(fmt.Sprintf("%s/%s.car", outputFolder, output.PieceCid), newFilePath)
			if err != nil {
				return fmt.Errorf("failed to rename file: %w", err)
			}
		}
	}
	// Get readers from outputFolder
	readers, err := utils.GetReaders(outputFolder)
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
	commpVal, pieceSize, _ := commpHasher.Digest()
	pieceCid, err := commcid.PieceCommitmentV1ToCID(commpVal)
	if err != nil {
		return fmt.Errorf("failed to get piece CID: %w", err)
	}

	aggregateName := uuid.New().String()
	// Use provided aggregateCarPath from context
	aggregateFolder := ctx.String("aggregate-car-path")
	err = os.MkdirAll(aggregateFolder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create aggregate folder: %w", err)
	}
	aggregatePath := fmt.Sprintf("%s%s.data", aggregateFolder, aggregateName)
	aggregateFile, err := os.Create(aggregatePath)
	if err != nil {
		return fmt.Errorf("failed to create aggregate file: %w", err)
	}
	defer aggregateFile.Close()

	// Generate aggregate
	readers, err = utils.GetReaders(outputFolder)
	if err != nil {
		return err
	}
	defer func() {
		for _, r := range readers {
			r.(io.Closer).Close()
		}
	}()
	out = mkpiece.MakeDataSegmentPiece(readers)
	_, err = io.Copy(aggregateFile, out)
	if err != nil {
		return fmt.Errorf("failed to copy aggregate: %w", err)
	}
	carSize, err := utils.GetFileSize(aggregatePath)
	if carSize == 0 || err != nil {
		return fmt.Errorf("failed to get aggregate file size: %w", err)
	}

	var buf buffer.Buffer
	var bufferResp *buffer.Response
	switch bufferType {
	case "lighthouse":
		buf = buffer.NewLighthouseBuffer(lighthouseApiKey, ctx.String("lighthouse-download-url"))
		bufferResp, err = buf.Store(aggregatePath)
	default:
		localBuf := buffer.NewLocalBuffer(ctx.Int("port")).(interface {
			buffer.Buffer
			StoreForServer(filePath string) (*buffer.Response, error)
		})
		bufferResp, err = localBuf.StoreForServer(aggregatePath)
	}

	if err != nil {
		return fmt.Errorf("failed to store in buffer: %w", err)
	}

	// Prepare deal parameters
	dealParams := utils.DealParams{
		FileName:        bufferResp.Hash,
		StorageProvider: miner,
		PieceSize:       pieceSize,
		CommpCid:        pieceCid.String(),
		CarFileSize:     uint64(carSize),
		PayloadCid:      ctx.String("payload-cid"),
		Duration:        uint64(ctx.Uint("duration")),
		StoragePrice:    uint64(ctx.Uint("storage-price")),
		Verified:        ctx.Bool("verified") || (ctx.Bool("testnet") && !ctx.IsSet("verified")),
		DownloadURL:     bufferResp.URL,
	}

	if err := utils.InitiateDeal(dealParams); err != nil {
		return fmt.Errorf("failed to initiate deal: %w", err)
	}
	return nil
}

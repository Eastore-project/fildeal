package deal

import (
	"fmt"
	"path/filepath"

	"github.com/eastore-project/fildeal/src/buffer"
	utils "github.com/eastore-project/fildeal/src/deal/utils"

	"github.com/urfave/cli/v2"
)

func MakeDeal(ctx *cli.Context) error {
	outDir := ctx.String("aggregate-car-path")
	path := ctx.String("input")
	miner := ctx.String("miner")
	bufferType := ctx.String("buffer")
	lighthouseApiKey := ctx.String("lighthouse-api-key")

	if bufferType == "lighthouse" && lighthouseApiKey == "" {
		return fmt.Errorf("lighthouse API key is required when using lighthouse buffer")
	}

	if ctx.Uint("duration") < 518400 || ctx.Uint("duration") > 1814400 {
		return fmt.Errorf("duration must be between 518400 (6 months) and 181440 (app. 3.5 years)")
	}

	output, err := utils.ConvertToCar(path, outDir, path)
	if err != nil {
		return fmt.Errorf("failed to convert to car: %w", err)
	}
	aggregatePath := filepath.Join(outDir, fmt.Sprintf("%s.car", output.PieceCid))

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
		PieceSize:       output.PieceSize,
		CommpCid:        output.PieceCid,
		CarFileSize:     output.CarSize,
		PayloadCid:      output.DataCid,
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

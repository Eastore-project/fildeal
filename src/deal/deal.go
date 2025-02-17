package deal

import (
	"fmt"
	"path/filepath"

	utils "github.com/eastore-project/fildeal/src/deal/utils"

	"github.com/urfave/cli/v2"
)

func MakeDeal(ctx *cli.Context, path string, miner string) error {
	outDir := ctx.String("aggregate-car-path")

	// directly convert the file/folder to car
	output, err := utils.ConvertToCar(path, outDir, path)
	if err != nil {
		return fmt.Errorf("failed to convert to car: %w", err)
	}
	aggregatePath := filepath.Join(outDir, fmt.Sprintf("%s.car", output.PieceCid))

	var aggregateName = output.PieceCid

	// Set the payload CID directly in the context's parent set
	if err := ctx.Set("payload-cid", output.DataCid); err != nil {
		return fmt.Errorf("failed to set payload-cid: %w", err)
	}

	if ctx.String("buffer") == "lighthouse" {
		// Upload the aggregate file to Lighthouse using api key from context
		lighthouseResp, err := utils.UploadToLighthouse(aggregatePath, ctx.String("lighthouse-api-key"))
		if err != nil {
			return fmt.Errorf("failed to upload to Lighthouse: %w", err)
		}
		fmt.Printf("File uploaded to Lighthouse. CID: %s, Name: %s, Size: %s\n",
			lighthouseResp.Hash, lighthouseResp.Name, lighthouseResp.Size)
		aggregateName = lighthouseResp.Hash
	}
	// Use the updated context
	if err := utils.InitiateDeal(aggregateName, miner, output.PieceSize, output.PieceCid, output.CarSize, ctx); err != nil {
		return fmt.Errorf("failed to initiate deal: %w", err)
	}
	return nil
}

package dealutils

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/urfave/cli/v2"
)

func InitiateDeal(fileName string, storageProvider string, pieceSize uint64, commpCid string, carFileSize uint64, ctx *cli.Context) error {
	var payloadCid = ctx.String("payload-cid")

	command := "boost"
	var url string
	var verified string

	// Set verified flag based on testnet and user input
	if ctx.Bool("verified") || (ctx.Bool("testnet") && !ctx.IsSet("verified")) {
		verified = "true"
	} else {
		verified = "false"
	}

	if ctx.String("buffer") == "lighthouse" {
		url = ctx.String("lighthouse-download-url") + fileName
	} else {
		url = fmt.Sprintf("http://localhost:%d/download/car?file_name=%s.data", ctx.Int("port"), fileName)
	}

	args := []string{
		"deal",
		"--provider=" + storageProvider,
		"--http-url=" + url,
		"--commp=" + commpCid,
		"--car-size=" + strconv.Itoa(int(carFileSize)),
		"--piece-size=" + strconv.Itoa(int(pieceSize)),
		"--payload-cid=" + payloadCid,
		"--duration=" + strconv.FormatUint(uint64(ctx.Uint("duration")), 10),
		"--storage-price=" + strconv.FormatUint(uint64(ctx.Uint("storage-price")), 10),
		"--verified=" + verified,
	}
	fmt.Println("Running command: ", command, args)
	dealResponse, err := exec.Command(command, args...).Output()
	if err != nil {
		fmt.Println(dealResponse)
		return fmt.Errorf("failed to initiate deal: %w", err)
	}
	fmt.Println("Deal initiated successfully for: " + fileName)
	fmt.Println("Deal Response: ", string(dealResponse))
	return nil
}

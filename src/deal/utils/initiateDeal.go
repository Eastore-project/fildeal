package dealutils

import (
	"fmt"
	"os/exec"
)

type DealParams struct {
	FileName        string
	StorageProvider string
	PieceSize       uint64
	CommpCid        string
	CarFileSize     uint64
	PayloadCid      string
	Duration        uint64
	StoragePrice    uint64
	Verified        bool
	DownloadURL     string
}

func InitiateDeal(params DealParams) error {
	verified := "false"
	if params.Verified {
		verified = "true"
	}

	// Construct command as a single string
	cmdStr := fmt.Sprintf("boost deal --provider=%s --http-url='%s' --commp=%s --car-size=%d --piece-size=%d --payload-cid=%s --duration=%d --storage-price=%d --verified=%s",
		params.StorageProvider,
		params.DownloadURL,
		params.CommpCid,
		params.CarFileSize,
		params.PieceSize,
		params.PayloadCid,
		params.Duration,
		params.StoragePrice,
		verified)

	fmt.Println("Running command:", cmdStr)

	cmd := exec.Command("bash", "-c", cmdStr)
	dealResponse, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(dealResponse))
		return fmt.Errorf("failed to initiate deal: %w", err)
	}

	fmt.Println("Deal initiated successfully for:", params.FileName)
	fmt.Println("Deal Response:\n", string(dealResponse))
	return nil
}

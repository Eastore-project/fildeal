package dealutils

import (
	"fmt"
	"os/exec"
	"strconv"
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
	command := "boost"

	verified := "false"
	if params.Verified {
		verified = "true"
	}

	args := []string{
		"deal",
		"--provider=" + params.StorageProvider,
		"--http-url=" + params.DownloadURL,
		"--commp=" + params.CommpCid,
		"--car-size=" + strconv.Itoa(int(params.CarFileSize)),
		"--piece-size=" + strconv.Itoa(int(params.PieceSize)),
		"--payload-cid=" + params.PayloadCid,
		"--duration=" + strconv.FormatUint(params.Duration, 10),
		"--storage-price=" + strconv.FormatUint(params.StoragePrice, 10),
		"--verified=" + verified,
	}
	fmt.Println("Running command: ", command, args)
	dealResponse, err := exec.Command(command, args...).Output()
	if err != nil {
		fmt.Println(dealResponse)
		return fmt.Errorf("failed to initiate deal: %w", err)
	}
	fmt.Println("Deal initiated successfully for: " + params.FileName)
	fmt.Println("Deal Response: ", string(dealResponse))
	return nil
}

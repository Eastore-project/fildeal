package main

import (
	"fildeal/src/deal"
	"fildeal/src/server"
	"fmt"

	"github.com/urfave/cli/v2"
)

func getDealCommand() *cli.Command {
	return &cli.Command{
		Name:  "podsi-deal",
		Usage: "Make a deal with a miner using podsi-aggregate for folder aggregation",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "Input folder containing files to make deal with",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "miner",
				Aliases:  []string{"m"},
				Usage:    "Miner ID to make the deal with",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "generate-car-path",
				Value:   "generated_car/",
				EnvVars: []string{"GENERATE_CAR_PATH"},
				Usage:   "Path for generated CAR files",
			},
			&cli.StringFlag{
				Name:    "aggregate-car-path",
				Value:   "aggregate_car_file/",
				EnvVars: []string{"AGGREGATE_CAR_PATH"},
				Usage:   "Path for aggregate CAR files",
			},
			&cli.StringFlag{
				Name:    "lighthouse-api-key",
				Value:   "",
				EnvVars: []string{"LIGHTHOUSE_API_KEY"},
				Usage:   "Lighthouse API key",
			},
			&cli.StringFlag{
				Name:    "lighthouse-download-url",
				Value:   "https://gateway.lighthouse.storage/ipfs/",
				EnvVars: []string{"LIGHTHOUSE_DOWNLOAD_URL"},
				Usage:   "Lighthouse download URL",
			},
			&cli.StringFlag{
				Name:  "buffer",
				Value: "localhost",
				Usage: "Buffer to use (localhost or lighthouse)",
			},
			&cli.StringFlag{
				Name:  "payload-cid",
				Value: "bafkreibtkdcncmofmavpdsar6msrmb2h4d7oetwtwtkz5cv3zsnwoyrrfq",
				Usage: "Payload CID for the deal",
			},
			&cli.UintFlag{
				Name:  "duration",
				Value: 518400,
				Usage: "Deal duration in epochs (minimum 518400 [6 months], maximum 1814400 [3.5 years])",
			},
			&cli.UintFlag{
				Name:  "storage-price",
				Value: 0,
				Usage: "Storage price in attoFIL per epoch per GiB",
			},
			&cli.BoolFlag{
				Name:  "verified",
				Usage: "Whether the deal is verified (default: true for testnet, false otherwise)",
			},
			&cli.BoolFlag{
				Name:  "server",
				Usage: "Start the server after initiating the deal",
			},
			&cli.BoolFlag{
				Name:  "testnet",
				Usage: "make deal on public testnet",
			},
		},
		Action: func(c *cli.Context) error {
			inputFolder := c.String("input")
			miner := c.String("miner")
			buffer := c.String("buffer")
			lighthouseApiKey := c.String("lighthouse-api-key")

			if buffer == "lighthouse" && lighthouseApiKey == "" {
				return fmt.Errorf("lighthouse API key is required when using lighthouse buffer")
			}

			if c.Uint("duration") < 518400 || c.Uint("duration") > 1814400 {
				return fmt.Errorf("duration must be between 518400 (6 months) and 181440 (app. 3.5 years)")
			}

			if err := deal.MakeDeal(c, inputFolder, miner); err != nil {
				return err
			}

			if c.Bool("server") {
				handler := server.SetupRouter()
				server.StartServer(c.Int("port"), handler)
			}
			return nil
		},
	}
}

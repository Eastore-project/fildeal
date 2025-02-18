package main

import (
	"github.com/eastore-project/fildeal/src/piece"
	"github.com/urfave/cli/v2"
)

func getPieceCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "podsi-aggregate",
			Usage: "Generate a data segment piece from all files in the input folder and write to output file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "input",
					Aliases:  []string{"i"},
					Usage:    "Input folder containing files to aggregate",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "output",
					Aliases:  []string{"o"},
					Usage:    "Output file path where the aggregated piece will be written",
					Required: true,
				},
				&cli.StringFlag{
					Name:    "proof-dir",
					Usage:   "Directory where inclusion proofs will be stored",
					Value:   "proofs-dir/",
					EnvVars: []string{"PROOFS_DIR"},
				},
			},
			Action: func(c *cli.Context) error {
				inputFolder := c.String("input")
				outputFile := c.String("output")
				proofDir := c.String("proof-dir")

				return piece.AggregateWithProofs(inputFolder, outputFile, proofDir)
			},
		},
		{
			Name:  "splitpiece",
			Usage: "Split a podsi-aggregate output file into pieces and save them in the output directory",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "input",
					Aliases:  []string{"i"},
					Usage:    "Input file (podsi-aggregate output) to split",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "output",
					Aliases:  []string{"o"},
					Usage:    "Output directory where pieces will be saved",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				inputFile := c.String("input")
				outputDir := c.String("output")
				return piece.SplitPiece(inputFile, outputDir)
			},
		},
	}
}

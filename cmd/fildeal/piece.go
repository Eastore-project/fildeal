package main

import (
	dealutils "fildeal/src/deal/utils"
	mkpiece "fildeal/src/mkpiece"
	"fmt"
	"io"
	"os"

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
			},
			Action: func(c *cli.Context) error {
				inputFolder := c.String("input")
				outputFile := c.String("output")

				readers, err := dealutils.GetReaders(inputFolder)
				if err != nil {
					return fmt.Errorf("failed to get readers from input folder: %w", err)
				}
				defer func() {
					for _, r := range readers {
						if closer, ok := r.(io.Closer); ok {
							closer.Close()
						}
					}
				}()

				out := mkpiece.MakeDataSegmentPiece(readers)

				f, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer f.Close()

				if _, err := io.Copy(f, out); err != nil {
					return fmt.Errorf("failed to write to output file: %w", err)
				}

				for _, reader := range readers {
					if _, err := reader.Read(make([]byte, 1)); err != io.EOF {
						return fmt.Errorf("reader not fully consumed")
					}
				}
				return nil
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
				return mkpiece.SplitPiece(inputFile, outputDir)
			},
		},
	}
}

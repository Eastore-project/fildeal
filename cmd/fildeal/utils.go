package main

import (
	"fildeal/src/index"
	"fildeal/src/utils"
	"fmt"

	"github.com/urfave/cli/v2"
)

func getUtilCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "cmp",
			Usage: "Compare two files and find the offset of the child file in the parent",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "parent",
					Aliases:  []string{"p"},
					Usage:    "Path to the parent file",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "child",
					Aliases:  []string{"c"},
					Usage:    "Path to the child file",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				parentPath := c.String("parent")
				childPath := c.String("child")

				offset, err := utils.FindOffset(parentPath, childPath)
				if err != nil {
					return err
				}
				fmt.Printf("Child file starts at offset %d in parent\n", offset)
				return nil
			},
		},
		{
			Name:  "boost-index",
			Usage: "Parse and index a file similar to Boost",
			Action: func(c *cli.Context) error {
				if c.NArg() != 1 {
					fmt.Println("Usage: fildeal boost-index <file>")
					return nil
				}
				filePath := c.Args().Get(0)
				return index.BoostIndex(filePath)
			},
		},
	}
}

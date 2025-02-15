package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "fildeal",
		Usage: "Filecoin Deals CLI",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Value:   8000,
				EnvVars: []string{"PORT"},
				Usage:   "Port for server",
			},
		},
		Commands: append(
			[]*cli.Command{getDealCommand()},
			append(getPieceCommands(), getUtilCommands()...)...,
		),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"github.com/eaardal/dig/commands"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(*cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
		Commands: []*cli.Command{
			commands.MsgCommand,
			commands.JobCommand,
			commands.SyncCommand,

			// TODO: Remove temporary command
			commands.TmpCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

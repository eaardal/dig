package commands

import (
	"github.com/eaardal/dig/digfile"
	"github.com/urfave/cli/v2"
)

var TmpCommand = &cli.Command{
	Name:  "tmp",
	Usage: "Temporary command",
	Action: func(c *cli.Context) error {
		digf, err := digfile.Read()
		if err != nil {
			return err
		}

		for _, job := range digf.Jobs {
			job.Print()
			println()
		}

		return nil
	},
}

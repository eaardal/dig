package commands

import (
	"fmt"
	"github.com/eaardal/dig/digfile"
	"github.com/urfave/cli/v2"
)

var JobListCommand = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "List jobs",
	Action: func(c *cli.Context) error {
		digf, err := digfile.Read()
		if err != nil {
			return err
		}

		for i, job := range digf.Jobs {
			println(fmt.Sprintf("%d - %s", i, job.Name))
		}

		return nil
	},
}

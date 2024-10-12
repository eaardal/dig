package commands

import (
	"github.com/eaardal/dig/digfile"
	"github.com/urfave/cli/v2"
)

var JobGetDefaultCommand = &cli.Command{
	Name:  "get-default",
	Usage: "Get the default job",
	Action: func(c *cli.Context) error {
		digf, err := digfile.Read()
		if err != nil {
			return err
		}

		job := digf.GetDefaultJob()
		if job != nil {
			println(job.Name)
			return nil
		}

		println("No default job set")

		return nil
	},
}

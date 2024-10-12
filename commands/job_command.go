package commands

import "github.com/urfave/cli/v2"

var JobCommand = &cli.Command{
	Name:  "job",
	Usage: "Manage jobs",
	Subcommands: []*cli.Command{
		JobCreateCommand,
		JobListCommand,
		JobRemoveCommand,
	},
}

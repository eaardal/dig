package commands

import (
	"fmt"
	"github.com/eaardal/dig/digfile"
	"github.com/urfave/cli/v2"
)

var JobSetDefaultCommand = &cli.Command{
	Name:  "set-default",
	Usage: "Set default job",
	Action: func(c *cli.Context) error {
		args, err := parseJobSetDefaultCommandArgs(c.Args())
		if err != nil {
			return err
		}

		digf, err := digfile.Read()
		if err != nil {
			return err
		}

		if args.jobIndex != nil && *args.jobIndex >= len(digf.Jobs) {
			return cli.Exit(fmt.Sprintf("job index %d out of bounds: %d jobs total (index %d-%d)", *args.jobIndex, len(digf.Jobs), 0, len(digf.Jobs)-1), 1)
		}

		if args.jobIndex != nil {
			for i, job := range digf.Jobs {
				job.IsDefault = i == *args.jobIndex
			}
		}

		if args.jobName != nil {
			for _, job := range digf.Jobs {
				job.IsDefault = job.Name == *args.jobName
			}
		}

		if err = digfile.Write(*digf); err != nil {
			return err
		}

		if args.jobIndex != nil {
			println(fmt.Sprintf("Job at index %d set as default", *args.jobIndex))
		} else if args.jobName != nil {
			println(fmt.Sprintf("Job %s set as default", *args.jobName))
		}

		return nil
	},
}

type jobSetDefaultCommandArgs struct {
	jobIndex *int
	jobName  *string
}

func parseJobSetDefaultCommandArgs(args cli.Args) (*jobSetDefaultCommandArgs, error) {
	jobIndex, jobName, err := parseJobNameOrIndex(args.Get(0), true)
	if err != nil {
		return nil, err
	}

	if jobIndex == nil && jobName == nil {
		return nil, cli.Exit("job index or name is required", 1)
	}

	return &jobSetDefaultCommandArgs{
		jobIndex: jobIndex,
		jobName:  jobName,
	}, nil
}

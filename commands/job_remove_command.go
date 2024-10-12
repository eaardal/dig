package commands

import (
	"fmt"
	"github.com/eaardal/dig/digfile"
	"github.com/urfave/cli/v2"
	"log"
	"strconv"
)

var JobRemoveCommand = &cli.Command{
	Name:    "remove",
	Usage:   "Remove a job",
	Aliases: []string{"rm"},
	Action: func(c *cli.Context) error {
		args, err := parseJobRemoveCommandArgs(c.Args())
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
			digf.Jobs = append(digf.Jobs[:*args.jobIndex], digf.Jobs[*args.jobIndex+1:]...)
		} else if args.jobName != nil {
			for i, job := range digf.Jobs {
				if job.Name == *args.jobName {
					digf.Jobs = append(digf.Jobs[:i], digf.Jobs[i+1:]...)
					break
				}
			}
		}

		if err = digfile.Write(*digf); err != nil {
			return err
		}

		if args.jobName != nil {
			log.Printf("Job %s removed", *args.jobName)
		} else if args.jobIndex != nil {
			log.Printf("Job at index %d removed", *args.jobIndex)
		}

		return nil
	},
}

type jobRemoveCommandArgs struct {
	jobIndex *int
	jobName  *string
}

func parseJobRemoveCommandArgs(args cli.Args) (*jobRemoveCommandArgs, error) {
	nameOrIndex := args.Get(0)

	if nameOrIndex == "" {
		return nil, cli.Exit("job name or index is required", 1)
	}

	var jobName *string
	var jobIndex *int

	index, err := strconv.Atoi(nameOrIndex)
	if err != nil {
		jobName = &nameOrIndex
	} else {
		jobIndex = &index
	}

	return &jobRemoveCommandArgs{
		jobIndex: jobIndex,
		jobName:  jobName,
	}, nil
}

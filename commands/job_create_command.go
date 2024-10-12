package commands

import (
	"fmt"
	"github.com/eaardal/dig/digfile"
	"github.com/urfave/cli/v2"
	"log"
)

var JobCreateCommand = &cli.Command{
	Name:  "create",
	Usage: "Create a new job",
	Action: func(c *cli.Context) error {
		args, err := parseJobCreateCommandArgs(c.Args())
		if err != nil {
			return fmt.Errorf("invalid args: %w", err)
		}

		log.Printf("Creating job %s", args.jobName)

		newJob := &digfile.Job{
			Name:       args.jobName,
			Kubernetes: nil,
		}

		digf, err := digfile.Read()
		if err != nil {
			return fmt.Errorf("failed to read digfile: %w", err)
		}

		digf.Jobs = append(digf.Jobs, newJob)

		if err = digfile.Write(*digf); err != nil {
			return fmt.Errorf("failed to write digfile: %w", err)
		}

		log.Printf("Job %s created", args.jobName)

		return nil
	},
}

type jobCreateCommandArgs struct {
	jobName string
}

func parseJobCreateCommandArgs(args cli.Args) (*jobCreateCommandArgs, error) {
	jobName := args.Get(0)

	if jobName == "" {
		return nil, cli.Exit("job name is required", 1)
	}

	return &jobCreateCommandArgs{
		jobName: jobName,
	}, nil
}

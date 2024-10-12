package commands

import (
	"github.com/urfave/cli/v2"
	"strconv"
)

func parseJobNameOrIndex(nameOrIndex string) (*int, *string, error) {
	if nameOrIndex == "" {
		return nil, nil, cli.Exit("job name or index is required", 1)
	}

	var jobName *string
	var jobIndex *int

	index, err := strconv.Atoi(nameOrIndex)
	if err != nil {
		jobName = &nameOrIndex
	} else {
		jobIndex = &index
	}

	return jobIndex, jobName, nil
}

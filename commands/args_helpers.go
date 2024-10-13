package commands

import (
	"github.com/urfave/cli/v2"
	"strconv"
)

func parseJobNameOrIndex(nameOrIndex string, required bool) (*int, *string, error) {
	if nameOrIndex == "" && required {
		return nil, nil, cli.Exit("job name or index is required", 1)
	}

	if nameOrIndex == "" {
		return nil, nil, nil
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

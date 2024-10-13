package commands

import (
	"fmt"
	"github.com/eaardal/dig/digfile"
	"github.com/eaardal/dig/localstorage"
	"github.com/eaardal/dig/logentry"
	"github.com/eaardal/dig/logparser"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

var MsgCommand = &cli.Command{
	Name:  "msg",
	Usage: "Search for something in a log entry's message",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "job",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		args, err := parseMsgCommandArgs(c)
		if err != nil {
			return err
		}

		dig, err := digfile.Read()
		if err != nil {
			return err
		}

		job, err := dig.GetJob(args.jobName, args.jobIndex)
		if err != nil {
			return err
		}

		if job == nil {
			return cli.Exit("No job found. Try setting a job as default or specify a job name or index when invoking the command", 1)
		}

		fileCh := make(chan *localstorage.CacheFile)
		logFileCh := bridgeCacheFileToLogFile(fileCh)
		logEntriesCh := make(chan *logentry.LogEntry)

		group, gctx := errgroup.WithContext(c.Context)

		group.Go(func() error {
			return localstorage.StreamCacheFiles(job.Name, fileCh)
		})

		group.Go(func() error {
			return logparser.ParseLogFile(gctx, logFileCh, logEntriesCh)
		})

		if err := group.Wait(); err != nil {
			return err
		}

		return nil
	},
}

func bridgeCacheFileToLogFile(fileCh chan *localstorage.CacheFile) chan logparser.LogFile {
	logFileCh := make(chan logparser.LogFile)

	go func() {
		defer close(logFileCh)

		for file := range fileCh {
			logFileCh <- file // localstorage.CacheFile struct implements logparser.LogFile interface so we can just forward it as is.
		}
	}()

	return logFileCh
}

type msgCommandArgs struct {
	jobName  *string
	jobIndex *int
	query    string
}

func parseMsgCommandArgs(c *cli.Context) (*msgCommandArgs, error) {
	args := c.Args()
	if args.Len() == 0 {
		return nil, fmt.Errorf("no query provided")
	}

	query := args.Get(0)
	if query == "" {
		return nil, fmt.Errorf("empty query provided")
	}

	jobNameOrIndex := c.String("job")
	jobIndex, jobName, err := parseJobNameOrIndex(jobNameOrIndex, false)
	if err != nil {
		return nil, err
	}

	return &msgCommandArgs{
		jobName:  jobName,
		jobIndex: jobIndex,
		query:    query,
	}, nil
}

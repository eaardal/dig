package commands

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eaardal/dig/digfile"
	"github.com/eaardal/dig/localstorage"
	"github.com/eaardal/dig/logentry"
	"github.com/eaardal/dig/logparser"
	"github.com/eaardal/dig/search"
	"github.com/eaardal/dig/ui/logentrieslist"
	"github.com/eaardal/dig/viewcontroller"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

var MsgCommand = &cli.Command{
	Name:  "msg",
	Usage: "Search for something in a log entry's message",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "job",
			Usage:    "Name or index of the job to search in. If not specified, the default job will be used",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "before",
			Usage:    "Number of log entries before the search result to display",
			Value:    10,
			Required: false,
		},
		&cli.IntFlag{
			Name:     "after",
			Usage:    "Number of log entries after the search result to display",
			Value:    10,
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "sync",
			Usage:    "Sync logs before searching",
			Value:    false,
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "case-sensitive",
			Usage:    "Perform case-sensitive search",
			Value:    false,
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "exact",
			Usage:    "Perform exact match search. If not specified, the search query will be treated as a substring",
			Value:    false,
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

		if args.syncFirst {
			if err := runSyncCommand(c, job.Name); err != nil {
				return err
			}
		}

		searchParams := search.Params{
			Query:         args.query,
			InMessage:     true,
			CaseSensitive: args.caseSensitive,
			Exact:         args.exactMatch,
		}

		fileCh := make(chan *localstorage.CacheFile)
		logFileCh := bridgeCacheFileToLogFile(fileCh)
		logEntriesCh := make(chan *logentry.LogEntry)
		searchResultsCh := make(chan *search.Result)

		group, gctx := errgroup.WithContext(c.Context)

		group.Go(func() error {
			return localstorage.StreamCacheFiles(job.Name, fileCh)
		})

		group.Go(func() error {
			return logparser.ParseLogFile(gctx, logFileCh, logEntriesCh)
		})

		group.Go(func() error {
			return search.Search(gctx, searchParams, logEntriesCh, searchResultsCh)
		})

		group.Go(func() error {
			results, err := search.GroupSearchResults(gctx, searchResultsCh)
			if err != nil {
				return err
			}

			opts := viewcontroller.Options{
				NumLogEntriesBefore: args.numLogEntriesBefore,
				NumLogEntriesAfter:  args.numLogEntriesAfter,
			}

			viewEntries, err := viewcontroller.PrepareSearchResultsForDisplay(results, opts)
			if err != nil {
				return err
			}

			model := logentrieslist.NewModel(viewEntries)

			_, err = tea.NewProgram(model, tea.WithAltScreen()).Run()
			if err != nil {
				return err
			}

			return nil
		})

		if err := group.Wait(); err != nil {
			return err
		}

		return nil
	},
}

func runSyncCommand(c *cli.Context, jobName string) error {
	syncFlags := &flag.FlagSet{}
	syncFlags.String("job", jobName, "")

	syncCtx := cli.NewContext(c.App, syncFlags, c)
	return SyncCommand.Action(syncCtx)
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
	jobName             *string
	jobIndex            *int
	query               string
	numLogEntriesBefore int
	numLogEntriesAfter  int
	syncFirst           bool
	caseSensitive       bool
	exactMatch          bool
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
		jobName:             jobName,
		jobIndex:            jobIndex,
		query:               query,
		numLogEntriesBefore: c.Int("before"),
		numLogEntriesAfter:  c.Int("after"),
		syncFirst:           c.Bool("sync"),
		caseSensitive:       c.Bool("case-sensitive"),
		exactMatch:          c.Bool("exact"),
	}, nil
}

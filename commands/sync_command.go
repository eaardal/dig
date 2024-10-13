package commands

import (
	"fmt"
	"github.com/eaardal/dig/digfile"
	"github.com/eaardal/dig/k8s"
	"github.com/eaardal/dig/localstorage"
	"github.com/eaardal/dig/ui"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

var SyncCommand = &cli.Command{
	Name:  "sync",
	Usage: "Syncs Kubernetes logs to files on your local machine",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "job",
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		args, err := parseSyncCommandArgs(c)
		if err != nil {
			return err
		}

		digf, err := digfile.Read()
		if err != nil {
			return err
		}

		job, err := digf.GetJob(args.jobName, args.jobIndex)
		if err != nil {
			return err
		}

		if job == nil {
			return cli.Exit("No job found. Try setting a job as default or specify a job name or index when invoking the command", 1)
		}

		if args.jobName == nil && args.jobIndex == nil {
			ui.Write(fmt.Sprintf("Using the default job: %s", job.Name))
		}

		if err := digfile.ValidateKubernetesConfigExistsForJob(job); err != nil {
			return cli.Exit(fmt.Sprintf("Job is invalid: %v", err), 1)
		}

		client, err := k8s.Client(job.Kubernetes.ContextName, job.Kubernetes.Namespace)
		if err != nil {
			return err
		}

		k8sLogChunks := make(chan *k8s.LogChunk)
		cacheFiles := make(chan *localstorage.CacheFile)

		group, gctx := errgroup.WithContext(c.Context)

		group.Go(func() error {
			return k8s.ReadLogs(gctx, client, job.Kubernetes.Namespace, job.Kubernetes.DeploymentNames, k8sLogChunks)
		})

		group.Go(func() error {
			return mapKubernetesLogChunksToCacheFile(k8sLogChunks, cacheFiles)
		})

		group.Go(func() error {
			return localstorage.SaveFileToCache(job.Name, cacheFiles)
		})

		if err := group.Wait(); err != nil {
			return err
		}

		ui.Write("Sync complete")

		return nil
	},
}

func mapKubernetesLogChunksToCacheFile(sourceCh <-chan *k8s.LogChunk, sinkCh chan<- *localstorage.CacheFile) error {
	defer close(sinkCh)

	for logMsg := range sourceCh {
		cacheFile := mapLogMsgToCacheFile(logMsg)
		sinkCh <- cacheFile
	}

	return nil
}

func mapLogMsgToCacheFile(logMsg *k8s.LogChunk) *localstorage.CacheFile {
	return &localstorage.CacheFile{
		FileName:    fmt.Sprintf("%s.log", logMsg.Origin),
		FileContent: logMsg.LogChunk,
	}
}

type syncCommandArgs struct {
	jobName  *string
	jobIndex *int
}

func parseSyncCommandArgs(c *cli.Context) (*syncCommandArgs, error) {
	args := c.Args()
	if args.Len() == 0 {
		return &syncCommandArgs{
			jobName:  nil,
			jobIndex: nil,
		}, nil
	}

	jobIndex, jobName, err := parseJobNameOrIndex(args.Get(0), false)
	if err != nil {
		return nil, err
	}

	return &syncCommandArgs{
		jobName:  jobName,
		jobIndex: jobIndex,
	}, nil
}

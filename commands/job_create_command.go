package commands

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eaardal/dig/config"
	"github.com/eaardal/dig/digfile"
	"github.com/eaardal/dig/k8s"
	"github.com/eaardal/dig/ui"
	"github.com/eaardal/dig/ui/interactiveselectlist"
	"github.com/urfave/cli/v2"
)

var JobCreateCommand = &cli.Command{
	Name:  "create",
	Usage: "Create a new job",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "namespace",
			EnvVars:  []string{config.KubernetesNamespaceEnvVar},
			Required: false,
		},
		&cli.StringFlag{
			Name:     "context",
			EnvVars:  []string{config.KubernetesContextEnvVar},
			Required: false,
		},
	},
	Action: func(c *cli.Context) error {
		args, err := parseJobCreateCommandArgs(c)
		if err != nil {
			return fmt.Errorf("invalid args: %w", err)
		}

		newJob := &digfile.Job{
			Name:       args.jobName,
			Kubernetes: nil,
		}

		digf, err := digfile.Read()
		if err != nil {
			return fmt.Errorf("failed to read digfile: %w", err)
		}

		namespace := args.k8sNamespace
		if namespace == "" {
			namespace, err = k8s.ResolveNamespace()
			if err != nil {
				return err
			}
		}

		if namespace == "" {
			return fmt.Errorf("namespace is required but could not be resolved via args, environment variables or kubens")
		}

		k8sContext := args.k8sContext
		if k8sContext == "" {
			k8sContext, err = k8s.ResolveContext()
			if err != nil {
				return err
			}
		}

		if k8sContext == "" {
			return fmt.Errorf("context is required but could not be resolved via args, environment variables or kubectx")
		}

		client, err := k8s.Client(k8sContext, namespace)
		if err != nil {
			return err
		}

		deployments, err := k8s.Deployments(c.Context, client, namespace)
		if err != nil {
			return err
		}

		listItems := make([]interactiveselectlist.ListItem, 0)
		for _, deployment := range deployments {
			listItems = append(listItems, interactiveselectlist.ListItem{
				Value:      deployment.Name,
				IsSelected: false,
			})
		}

		model := interactiveselectlist.NewModel(listItems, fmt.Sprintf("Which deployments should be included in this job?"))

		appState, err := tea.NewProgram(model).Run()
		if err != nil {
			return err
		}

		selectedDeployments := appState.(interactiveselectlist.Model).GetSelectedChoices()

		newJob.Kubernetes = &digfile.KubernetesJob{
			ContextName:     k8sContext,
			Namespace:       namespace,
			DeploymentNames: selectedDeployments,
		}

		digf.Jobs = append(digf.Jobs, newJob)

		if err = digfile.Write(*digf); err != nil {
			return fmt.Errorf("failed to write digfile: %w", err)
		}

		ui.Write("Job %s created", args.jobName)
		ui.Write(" - Context: %s", k8sContext)
		ui.Write(" - Namespace: %s", namespace)
		ui.Write(" - Deployments:")
		for _, deployment := range newJob.Kubernetes.DeploymentNames {
			ui.Write("    - %s", deployment)
		}

		return nil
	},
}

type jobCreateCommandArgs struct {
	jobName      string
	k8sNamespace string
	k8sContext   string
}

func parseJobCreateCommandArgs(c *cli.Context) (*jobCreateCommandArgs, error) {
	args := c.Args()
	if args.Len() == 0 {
		return nil, cli.Exit("job name is required", 1)
	}

	jobName := args.Get(0)
	if jobName == "" {
		return nil, cli.Exit("job name is required", 1)
	}

	return &jobCreateCommandArgs{
		jobName:      jobName,
		k8sNamespace: c.String("namespace"),
		k8sContext:   c.String("context"),
	}, nil
}

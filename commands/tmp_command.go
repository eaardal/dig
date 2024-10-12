package commands

import (
	"fmt"
	"github.com/eaardal/dig/k8s"
	"github.com/urfave/cli/v2"
)

var TmpCommand = &cli.Command{
	Name:  "tmp",
	Usage: "Temporary command",
	Action: func(c *cli.Context) error {
		client, err := k8s.Client2()
		if err != nil {
			return err
		}

		deployments, err := k8s.Deployments(c.Context, client)
		if err != nil {
			return err
		}

		for _, deployment := range deployments {
			println(fmt.Sprintf("Deployment: %s", deployment.Name))
		}

		return nil
	},
}

package commands

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eaardal/dig/k8s"
	"github.com/eaardal/dig/ui/interactiveselectlist"
	"github.com/urfave/cli/v2"
)

var TmpCommand = &cli.Command{
	Name:  "tmp",
	Usage: "Temporary command",
	Action: func(c *cli.Context) error {

		namespace, err := k8s.ResolveNamespace()
		if err != nil {
			return err
		}

		k8sContext, err := k8s.ResolveContext()
		if err != nil {
			return err
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

		// Cast the final model back and get selected choices
		selectedDeployments := appState.(interactiveselectlist.Model).GetSelectedChoices()
		fmt.Println("Selected deployments:", selectedDeployments)

		return nil
	},
}

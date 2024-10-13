package commands

import (
	"github.com/eaardal/dig/ui/logentrieslist"
	"github.com/urfave/cli/v2"
)

var TmpCommand = &cli.Command{
	Name:  "tmp",
	Usage: "Temporary command",
	Action: func(c *cli.Context) error {
		logentrieslist.Example()
		return nil
	},
}

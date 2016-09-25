package get

import (
	"io"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//NewGetCmd forms the basis of the Get command set
func NewGetCmd(in io.Reader, out io.Reader, store storage.Storer) cli.Command {
	apps := appsCmd{}
	return cli.Command{
		Name:  "get",
		Usage: "get <resource>",
		Subcommands: []cli.Command{
			cli.Command{
				Name: "projects",
			},
			cli.Command{
				Name:   "apps",
				Usage:  "",
				Action: apps.appsAction,
			},
		},
	}
}

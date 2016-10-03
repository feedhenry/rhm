package get

import (
	"io"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//NewGetCmd forms the basis of the Get command set
func NewGetCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "get <resource>",
		Subcommands: []cli.Command{
			NewProjectsCmd(in, out, store),
			NewProjectCmd(in, out, store),
			NewTemplatesCmd(in, out, store),
		},
	}
}

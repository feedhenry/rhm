package delete

import (
	"io"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

// NewDeleteCmd configures the different types of delete cmd
func NewDeleteCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	return cli.Command{
		Name: "delete",
		Subcommands: []cli.Command{
			NewProjectDeleteCmd(in, out, store),
		},
	}
}

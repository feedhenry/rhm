package new

import (
	"io"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//Aggregator command for creating new business objects like:
// - Project
// - App

//CreateNewCmd returns the definition of the new command
func CreateNewCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	return cli.Command{
		Name:        "new",
		Usage:       "new <resource>",
		Description: "Create new <resource>",
		Subcommands: []cli.Command{
			ProjectCreateCmd(in, out, store),
		},
	}
}

package use

import (
	"io"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//use contains the various subcommands of the command for example rhm use project <guid>
//when you use something it becomes set as your current working context alot like when a use a db

//NewUseCmd return the definition of the use command
func NewUseCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	projectSub := setProjectCmd{in: in, out: out, storage: store}
	return cli.Command{
		Name:        "use",
		Description: "use certain contexts such as project for other actions",
		Subcommands: []cli.Command{
			{
				Name:        "project",
				Usage:       "<guid>",
				Description: "use a specified project for all following project actions",
				Action:      projectSub.SetProjectAction,
			},
		},
	}
}

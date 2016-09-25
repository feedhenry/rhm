package commands

import (
	"io"

	"text/template"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//The context command will tell you which domain you're currently targeting, which user you're logged in with and which project you're currently using.

//NewContextCmd will return a configured contextCmd
func NewContextCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	cc := contextCmd{
		in:    in,
		out:   out,
		store: store,
	}
	return cc.contextCmd()
}

type contextCmd struct {
	in    io.Reader
	out   io.Writer
	store storage.Storer
}

//define the cli.Command
func (cc contextCmd) contextCmd() cli.Command {
	return cli.Command{
		Name:        "context",
		Description: "context will give you the context you are currently working in. Domain, Project, User",
		Action:      cc.contextAction,
	}
}

//this is the template for outputting the context with

var contextTemplate = `
| Domain  : {{.Domain}} 
| Host    : {{.Host}}   
| User    : {{.UserName}}
| Project : {{.ActiveProject}}
`

//main handler for the command
func (cc contextCmd) contextAction(ctx *cli.Context) error {
	userData, err := cc.store.ReadUserData()
	if err != nil {
		return cli.NewExitError("failed to get context "+err.Error(), 1)
	}
	t := template.New("contextTemplate")
	t, err = t.Parse(contextTemplate)
	if err != nil {
		return cli.NewExitError("failed to parse contextTemplate "+err.Error(), 1)
	}
	if err := t.Execute(cc.out, userData); err != nil {
		return cli.NewExitError("failed to output template result "+err.Error(), 1)
	}
	return nil
}

package use

import (
	"io"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//Project is used to set the active project context. This allows for the user to use rhm from within a project's context without needing
//to keep passing the obscure guid.

type setProjectCmd struct {
	in      io.Reader
	out     io.Writer
	storage storage.Storer
}

//SetProject defines the interface for the cli tool for the project context command
func (sp setProjectCmd) SetProject() cli.Command {
	return cli.Command{
		Name:        "project",
		Action:      sp.SetProjectAction,
		Description: "sets your currently active project. Can be overriden by using the --project flag",
		Usage:       "<guid>",
	}
}

//SetProjectAction is the main logic for the set project command
func (sp setProjectCmd) SetProjectAction(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		return cli.NewExitError("a project guid is required", 1)
	}
	projectGUID := ctx.Args()[0]
	if len(projectGUID) != 24 {
		return cli.NewExitError("a guid should be 24 chars long", 1)
	}
	uData, err := sp.storage.ReadUserData()
	if err != nil {
		return cli.NewExitError("failed to set project "+err.Error(), 1)
	}
	uData.ActiveProject = projectGUID
	if err := sp.storage.WriteUserData(uData); err != nil {
		return cli.NewExitError("failed to save project context : "+err.Error(), 1)
	}
	return nil
}

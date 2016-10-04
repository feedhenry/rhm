package delete

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/ui"
	"github.com/urfave/cli"
)

type projectDelete struct {
	in          io.Reader
	out         io.Writer
	store       storage.Storer
	force       bool
	projectGUID string
	deleter     func(*http.Request) (*http.Response, error)
}

// NewProjectDeleteCmd returns a configured projectDeleteCmd
func NewProjectDeleteCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	pd := &projectDelete{
		in:      in,
		out:     out,
		store:   store,
		deleter: http.DefaultClient.Do,
	}
	return pd.project()
}

func (pd *projectDelete) project() cli.Command {
	return cli.Command{
		Name:        "project",
		Action:      pd.deleteAction,
		Description: "deletes a project from rhmap",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Destination: &pd.projectGUID,
				Usage:       "set the project to delete --project=<guid>",
			},
			cli.BoolFlag{
				Name:        "f",
				Destination: &pd.force,
				Usage:       "force the project delete. You will not be prompted before the delete.",
			},
		},
	}
}

func (pd *projectDelete) deleteAction(ctx *cli.Context) error {
	var urlTemplate = "%s/box/api/projects/%s"
	userData, err := pd.store.ReadUserData()
	if err != nil {
		return err
	}
	var guid = userData.ActiveProject
	if pd.projectGUID == "" && guid == "" {
		return cli.NewExitError("no project GUID passed. Use --project=<GUID>", 1)
	}
	if pd.projectGUID == "" {
		pd.projectGUID = guid
	}
	fullURL := fmt.Sprintf(urlTemplate, userData.Host, pd.projectGUID)
	if !pd.force {
		answer, err := ui.WaitForAnswer("Are you sure you want to delete project with GUID: "+pd.projectGUID, pd.out, pd.in)
		if err != nil {
			return err
		}
		answer = strings.ToLower(answer)
		if answer != "y" && answer != "yes" {
			return cli.NewExitError("project not deleted.", 1)
		}
	}
	req, err := http.NewRequest("delete", fullURL, nil)
	if err != nil {
		return cli.NewExitError("failed to create request", 1)
	}
	cookie := http.Cookie{Name: "feedhenry", Value: userData.Auth}
	req.AddCookie(&cookie)
	//make delete request
	res, err := pd.deleter(req)
	if err != nil {
		return cli.NewExitError("failed to make delete request "+err.Error(), 1)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return cli.NewExitError("failed to make delete request "+err.Error(), 1)
	}
	if res.StatusCode != http.StatusOK {
		pd.out.Write(data)
		return cli.NewExitError(fmt.Sprintf("\n unexpected response %d \n", res.StatusCode), 1)
	}
	return nil
}

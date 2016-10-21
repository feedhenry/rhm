package get

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/ui"
	"github.com/urfave/cli"
)

var projectTemplate = `
==== Project Data ====
| {{PadRight 14 " " "Title"}}| {{PadRight 30 " " "Email"}}| {{PadRight 28 " " "GUID"}}| {{PadRight 14 " " "Type"}}|
+-{{PadRight 14 "-" ""}}+-{{PadRight 30 "-" ""}}+-{{PadRight 28 "-" ""}}+-{{PadRight 14 "-" ""}}+
| {{PadRight 14 " " .Title}}| {{PadRight 30 " " .AuthorEmail}}| {{PadRight 28 " " .GUID}}| {{PadRight 14 " " .Type}}|
`

type projectCmd struct {
	in            io.Reader
	out           io.Writer
	store         storage.Storer
	project       string
	getter        commands.HTTPGetter
	projectFinder commands.ProjectFinder
}

func (pc *projectCmd) Project() cli.Command {
	return cli.Command{
		Name:        "project",
		Action:      pc.projectAction,
		Description: "get the project definition. If you have set rhm use project <guid> it will use that project. Or you can pass --project=<guid>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Destination: &pc.project,
				Usage:       "the project guid ",
			},
		},
	}
}

func (pc *projectCmd) projectAction(ctx *cli.Context) error {
	var urlTemplate = "%s/box/api/projects/%s"
	uData, err := pc.store.ReadUserData()
	if err != nil {
		return err
	}
	if pc.project == "" && uData.ActiveProject != "" {
		pc.project = uData.ActiveProject
	}
	if pc.project == "" {
		return cli.NewExitError("expeced a project guid or name. Use --project", 1)
	}
	guid, err := pc.projectFinder(pc.project, uData, pc.getter)
	if err == nil {
		pc.project = guid
	}
	url := fmt.Sprintf(urlTemplate, uData.Host, pc.project)
	newrequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	// create a cookie
	cookie := http.Cookie{Name: "feedhenry", Value: uData.Auth}
	newrequest.AddCookie(&cookie)
	res, err := pc.getter(newrequest)
	if err != nil {
		return cli.NewExitError("failed to make request to get project", 1)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return cli.NewExitError("failed to read response body "+err.Error(), 1)
		}
		pc.out.Write(data)
		return cli.NewExitError(fmt.Sprintf("\n unexpected response %d \n", res.StatusCode), 1)
	}

	var dataStructure commands.Project
	return ui.NewPrinter(ctx.GlobalString("o"), res.Body, pc.out, projectTemplate, &dataStructure).Print()
}

// NewProjectCmd configures a new project command
func NewProjectCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	pc := &projectCmd{
		in:            in,
		out:           out,
		store:         store,
		getter:        http.DefaultClient.Do,
		projectFinder: ProjectNameToGUID,
	}
	return pc.Project()
}

// ProjectNameToGUID will take a project title as an argument and find it in the core and return the guid for that name
func ProjectNameToGUID(title string, userData *storage.UserData, getter commands.HTTPGetter) (string, error) {
	url := "%s/box/api/projects?apps=false"
	fullURL := fmt.Sprintf(url, userData.Host)
	newrequest, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	// create a cookie
	cookie := http.Cookie{Name: "feedhenry", Value: userData.Auth}
	newrequest.AddCookie(&cookie)
	//do request
	res, err := getter(newrequest)
	if err != nil {
		return "", cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	//check if not authed)

	//handle Output
	var projects []*commands.Project

	dec := json.NewDecoder(res.Body)
	if err := dec.Decode(&projects); err != nil {
		return "", err
	}

	for _, project := range projects {
		if project.Title == title {
			return project.GUID, nil
		}
	}

	return "", errors.New("Project with title '" + title + "' not found.")
}

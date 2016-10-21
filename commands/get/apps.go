package get

//handle the apps list for rhmap.

import (
	"fmt"
	"io"
	"net/http"

	"io/ioutil"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/ui"
	"github.com/urfave/cli"
)

var appsTemplate = `
== Apps in project: {{.Title}} ==

| {{PadRight 14 " " "App Name"}}|  {{PadRight 26 " " "GUID"}}|
|{{PadRight 15 "-" ""}}+{{PadRight 28 "-" ""}}|{{range .Apps }}
| {{PadRight 14 " " .Title}}|  {{PadRight 26 " " .GUID}}|{{end}}
`

//appsCmd constructs the required writer in order to send the response to the right place.
type appsCmd struct {
	out           io.Writer
	in            io.Reader
	getter        func(*http.Request) (*http.Response, error)
	project       string
	store         storage.Storer
	projectFinder commands.ProjectFinder
}

//Apps Defines our cli command including its flags and usage then returns the command to allow a user to do specific operations on projects
func (ac *appsCmd) Apps() cli.Command {
	return cli.Command{
		Name:        "apps",
		Action:      ac.appsAction,
		Usage:       "apps",
		Description: "apps allows listing apps in current project in rhm",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "project",
				Destination: &ac.project,
				Usage:       "the project guid ",
			},
		},
	}
}

//appsAction is where the logic is pulled together to perform the command. This function conforms to the cli action
func (ac *appsCmd) appsAction(ctx *cli.Context) error {
	url := "%s/box/api/projects/%s"

	userData, err := ac.store.ReadUserData()
	if err != nil {
		return cli.NewExitError("could not read userData "+err.Error(), 1)
	}
	if ac.project == "" && userData.ActiveProject != "" {
		ac.project = userData.ActiveProject
	}
	guid, err := ac.projectFinder(ac.project, userData, ac.getter)
	if err == nil {
		ac.project = guid
	}

	fullURL := fmt.Sprintf(url, userData.Host, ac.project)
	newrequest, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	// create a cookie
	cookie := http.Cookie{Name: "feedhenry", Value: userData.Auth}
	newrequest.AddCookie(&cookie)
	//do request
	res, err := ac.getter(newrequest)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}

	if res.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return cli.NewExitError("failed to read response body "+err.Error(), 1)
		}
		ac.out.Write(data)
		return cli.NewExitError(fmt.Sprintf("\n unexpected response %d \n", res.StatusCode), 1)
	}

	var dataStructure commands.Project
	return ui.NewPrinter(ctx.GlobalString("o"), res.Body, ac.out, appsTemplate, &dataStructure).Print()
}

//NewAppsCmd configures the appsCmd for use with the client
func NewAppsCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	var client http.Client
	ac := &appsCmd{out: out, in: in, getter: client.Do, store: store, projectFinder: ProjectNameToGUID}
	return ac.Apps()
}

package get

//handle the project list for rhmap.

import (
	"fmt"
	"io"
	"net/http"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/ui"
	"github.com/urfave/cli"
)

var projectsTemplate = `
| {{PadRight 14 " " "Project Name"}}|  {{PadRight 26 " " "GUID"}}|
|{{PadRight 15 "-" ""}}+{{PadRight 28 "-" ""}}|{{range . }}
| {{PadRight 14 " " .Title}}|  {{PadRight 26 " " .GUID}}|{{end}}
`

//projectsCmd constructs the required writer in order to send the response to the right place.
type projectsCmd struct {
	out    io.Writer
	in     io.Reader
	getter func(*http.Request) (*http.Response, error)
	store  storage.Storer
}

//Project Defines our cli command including its flags and usage then returns the command to allow a user to do specific operations on projects
func (pc *projectsCmd) Projects() cli.Command {
	return cli.Command{
		Name:        "projects",
		Action:      pc.projectsAction,
		Usage:       "project",
		Description: "projects allows listing projects in rhm",
	}
}

//projectAction is where the logic is pulled together to perform the command. This funtion conforms to the cli action
func (pc *projectsCmd) projectsAction(ctx *cli.Context) error {
	var (
		url = "%s/box/api/projects?apps=false"
	)
	userData, err := pc.store.ReadUserData()
	if err != nil {
		return cli.NewExitError("could not read userData "+err.Error(), 1)
	}
	fullURL := fmt.Sprintf(url, userData.Host)
	newrequest, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	// create a cookie
	cookie := http.Cookie{Name: "feedhenry", Value: userData.Auth}
	newrequest.AddCookie(&cookie)
	//do request
	res, err := pc.getter(newrequest)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}

	var dataStructure []*commands.Project
	return ui.NewPrinter(ctx.GlobalString("o"), res.Body, pc.out, projectsTemplate, &dataStructure).Print()
}

func handleProjectsResponseStatus(status int) error {
	if status == http.StatusOK {
		return nil
	}
	return cli.NewExitError(fmt.Sprintf("\n response %d \n", status), 1)
}

//NewProjectsCmd configures the projectsCmd for use with the client
func NewProjectsCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	var client http.Client
	pc := &projectsCmd{out: out, in: in, getter: client.Do, store: store}
	return pc.Projects()
}

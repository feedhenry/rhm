package get

//handle the project list for rhmap.

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

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
	resp, err := pc.getter(newrequest)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.NewExitError("failed to read response body "+err.Error(), 1)
	}
	//check if not authed)
	if err := handleProjectsResponseStatus(resp.StatusCode); err != nil {
		pc.out.Write(ret)
		return err
	}

	var resJSON []*commands.Project
	if err := json.Unmarshal(ret, &resJSON); err != nil {
		return cli.NewExitError("failed to decode response :"+err.Error(), 1)
	}

	t := template.New("project list template")
	t, _ = t.Parse("{{range . }} |  Project | {{.Title}}  | GUID | {{.GUID}} \n\n  {{end}}")
	if err := t.Execute(pc.out, resJSON); err != nil {
		return cli.NewExitError("failed to execute template "+err.Error(), 1)
	}

	return nil
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

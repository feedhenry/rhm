package commands

//handle the project logic for rhmap.

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	//"os"
	"text/template"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//ProjectCmd constructs the required writer in order to send the response to the right place.
type projectCmd struct {
	out      io.Writer
	in       io.Reader
	response func(*http.Request) (*http.Response, error)
	store    storage.Storer
}

//Project Defines our cli command including its flags and usage then returns the command to allow a user to do specific operations on projects
func (pc *projectCmd) Project() cli.Command {
	return cli.Command{
		Name:        "project",
		Action:      pc.projectAction,
		Usage:       "project",
		Description: "project will allow a user to create,update,delete and list projects in rhm",
	}
}

//projectAction is where the logic is pulled together to perform the command. This funtion conforms to the cli action
func (pc *projectCmd) projectAction(ctx *cli.Context) error {
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

	// do the request :)
	resp, err := pc.response(newrequest)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.NewExitError("project list request failed "+err.Error(), 1)
	}

	var resJSON []*Project
	if err := json.Unmarshal(ret, &resJSON); err != nil {
		return cli.NewExitError("failed to decode response", 1)
	}

	t := template.New("project list template")
	t, _ = t.Parse("Project : {{.Title}}  Guid : {{.Guid}} \n\n")
	for _, v := range resJSON {
		t.Execute(pc.out, v)
	}

	return nil
}

//NewProjectCmd configures the ProjectCmd for use with the client
func NewProjectCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	var client http.Client
	pc := &projectCmd{out: out, in: in, response: client.Do, store: store}
	return pc.Project()
}

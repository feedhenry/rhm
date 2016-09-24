package commands

//handle the project logic for rhmap.

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//ProjectCmd constructs the required writer in order to send the response to the right place.
type projectCmd struct {
	out      io.Writer
	in       io.Reader
	response func(*http.Request) (*http.Response, error)
	//response   http.Client
	newrequest func(string, string, io.Reader) (*http.Request, error)
	store      storage.Storer
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

//this is the data structure for posting to the server
type projectParams struct {
	Domain string `json:"d"`
}

//projectAction is where the logic is pulled together to perform the command. This funtion conforms to the cli action
func (pc *projectCmd) projectAction(ctx *cli.Context) error {
	var (
		url = "%s/box/api/projects"
	)
	userData, err := pc.store.ReadUserData()
	if err != nil {
		return cli.NewExitError("could not read userData "+err.Error(), 1)
	}
	fullURL := fmt.Sprintf(url, userData.Host)
	newrequest, err := pc.newrequest("GET", fullURL, nil)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}

	// create a cookie
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "feedhenry", Value: userData.Auth, Expires: expiration}
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

	// some juggling with the json return string
	var resJSON interface{}
	if err := json.Unmarshal(ret, &resJSON); err != nil {
		return cli.NewExitError("failed to decode response", 1)
	}
	k := resJSON.([]interface{})
	for _, v := range k {
		x := v.(map[string]interface{})
		fmt.Printf("Project %+v\n", x["title"])
	}

	return nil
}

//NewProjectCmd configures the ProjectCmd for use with the client
func NewProjectCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	var client http.Client
	pc := &projectCmd{out: out, in: in, response: client.Do, newrequest: http.NewRequest, store: store}
	return pc.Project()
}

package new

//handle the project list for rhmap.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/request"
	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/ui"
	"github.com/urfave/cli"
)

var projectTemplate = `
| Title |  {{.Title}} 
| Email |  {{.AuthorEmail}} 
| Guid  |  {{.GUID}}
| Type  |  {{.Type}}
        {{if .Apps}}
        | Apps |
                {{range .Apps }}
               | Title | {{.Title}}
               | Guid  | {{.GUID}}
            
                {{end}}
        {{end}}        
`

type projectCmd struct {
	out        io.Writer
	in         io.Reader
	store      storage.Storer
	httpClient func(*http.Request) (*http.Response, error)
	title      string
	template   string
}

//this is the data structure for posting to the server
type projectParams struct {
	Title    string `json:"title"`
	Template string `json:"template"`
}

// ProjectCreateCmd - creates command
func ProjectCreateCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	var client http.Client
	pc := &projectCmd{out: out, in: in, store: store, httpClient: client.Do}
	return pc.Project()
}

func (pc *projectCmd) Project() cli.Command {
	return cli.Command{
		Name:        "project",
		Action:      pc.projectAction,
		Usage:       "new project",
		Description: "Create new project",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "title",
				Destination: &pc.title,
				Usage:       "the project title ",
			},
			cli.StringFlag{
				Name:        "template",
				Destination: &pc.template,
				Usage:       "Name of the template for project",
			},
		},
	}
}

func (pc *projectCmd) projectAction(ctx *cli.Context) error {
	var (
		url = "%s/box/api/projects"
	)
	userData, err := pc.store.ReadUserData()
	if err != nil {
		return cli.NewExitError("could not read userData "+err.Error(), 1)
	}
	fullURL := fmt.Sprintf(url, userData.Host)
	projectData, err := pc.createRequestBody()
	postData, err := request.PrepareJSONBody(projectData)
	newrequest, err := http.NewRequest("POST", fullURL, postData)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	//newrequest.Header.Set("Content-Type", "application/json")
	cookie := http.Cookie{Name: "feedhenry", Value: userData.Auth}
	newrequest.AddCookie(&cookie)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	resp, err := pc.httpClient(newrequest)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.NewExitError("failed to read response body "+err.Error(), 1)
	}
	if err := handleProjectsResponseStatus(resp.StatusCode); err != nil {
		pc.out.Write(ret)
		return err
	}
	var resJSON []*commands.Project
	if err := json.Unmarshal(ret, &resJSON); err != nil {
		return cli.NewExitError("failed to decode response :"+err.Error(), 1)
	}

	t := template.New("project")
	t, err = t.Parse(projectTemplate)
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

func (pc *projectCmd) createRequestBody() (projectParams, error) {
	body := projectParams{}
	if pc.title == "" {
		title, err := ui.WaitForAnswer("Enter project name", pc.out, pc.in)
		if err != nil {
			return body, err
		}
		body.Title = title
	} else {
		body.Title = pc.title
	}
	if pc.template == "" {
		body.Template = "bare_project"
	} else {
		body.Template = pc.template
	}
	return body, nil
}

func prepareJSONBody(b interface{}) (io.Reader, error) {
	body, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(body), nil
}

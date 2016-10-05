package get

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

// NewGetTemplateCmd forms the basis of getting templates of various types
func (tc *templatesCmd) Templates() cli.Command {
	return cli.Command{
		Name:   "templates",
		Usage:  "templates",
		Action: tc.templatesAction,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "type",
				Destination: &tc.templateType,
				Usage:       "The type of the templates (e.g. projects, apps etc.) ",
				Value:       "projects",
			},
		},
	}
}

// templatesCmd constructs the required writer in order to send the response to the right place.
type templatesCmd struct {
	out          io.Writer
	in           io.Reader
	templateType string
	getter       func(*http.Request) (*http.Response, error)
	store        storage.Storer
}

// ListTemplates gets a list of templates of the supplied templateType (e.g. projects)
func (tc *templatesCmd) templatesAction(ctx *cli.Context) error {
	var url = "%s/box/api/templates/%s"

	userData, err := tc.store.ReadUserData()
	if err != nil {
		return cli.NewExitError("could not read userData "+err.Error(), 1)
	}

	fullURL := fmt.Sprintf(url, userData.Host, tc.templateType)

	newrequest, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}

	// create a cookie
	cookie := http.Cookie{Name: "feedhenry", Value: userData.Auth}
	newrequest.AddCookie(&cookie)
	//do request

	resp, err := tc.getter(newrequest)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	defer resp.Body.Close()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cli.NewExitError("failed to read response body "+err.Error(), 1)
	}

	//check if not authed)
	if err := handleTemplatesResponseStatus(resp.StatusCode); err != nil {
		tc.out.Write(ret)
		return err
	}

	var resJSON []*commands.ProjectTemplate
	if err := json.Unmarshal(ret, &resJSON); err != nil {
		return cli.NewExitError("failed to decode response :"+err.Error(), 1)
	}

	t := template.New("project templates list template")
	t, _ = t.Parse("{{range . }} | ID | {{.ID}} | Title | {{.Title}} | Category | {{.Category}} \n\n  {{end}}")
	if err := t.Execute(tc.out, resJSON); err != nil {
		return cli.NewExitError("failed to execute template "+err.Error(), 1)
	}

	return nil
}

// handleTemplatesResponseStatus checks whether the API request returned an ok response
func handleTemplatesResponseStatus(status int) error {
	if status == http.StatusOK {
		return nil
	}
	return cli.NewExitError(fmt.Sprintf("\n response %d \n", status), 1)
}

//NewTemplatesCmd configures the templatesCmd for use with the client
func NewTemplatesCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	var client http.Client
	tc := &templatesCmd{out: out, in: in, getter: client.Do, store: store}
	return tc.Templates()
}

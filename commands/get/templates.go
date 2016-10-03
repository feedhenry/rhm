package get

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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
			cli.StringFlag{
				Name:        "id",
				Destination: &tc.templateID,
				Usage:       "The ID of the template to retrieve, only returns for an exact match",
			},
			cli.StringFlag{
				Name:        "name",
				Destination: &tc.templateName,
				Usage:       "The name of the template to retrieve, will return any template which contains the provided string",
			},
		},
	}
}

// templatesCmd constructs the required writer in order to send the response to the right place.
type templatesCmd struct {
	out          io.Writer
	in           io.Reader
	getter       func(*http.Request) (*http.Response, error)
	store        storage.Storer
	templateType string
	templateID   string
	templateName string
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

	var resJSON []*commands.Template
	if err := json.Unmarshal(ret, &resJSON); err != nil {
		return cli.NewExitError("failed to decode response :"+err.Error(), 1)
	}

	if tc.templateID != "" {
		resJSON, err = filterTemplates(resJSON, tc.templateID, func(template *commands.Template, templateID string) bool {
			return template.ID == templateID
		})
	}
	if err != nil {
		return cli.NewExitError("Failed to find template with ID :"+err.Error(), 1)
	}

	if tc.templateName != "" {
		resJSON, err = filterTemplates(resJSON, tc.templateName, func(template *commands.Template, templateName string) bool {
			return strings.Contains(strings.ToLower(template.Name), strings.ToLower(templateName))
		})
	}
	if err != nil {
		return cli.NewExitError("Failed to find template with Title :"+err.Error(), 1)
	}

	t := template.New("project templates list template")
	t, _ = t.Parse("{{range . }} | ID | {{.ID}} | Name | {{.Name}} | Category | {{.Category}} \n\n  {{end}}")
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

// filterTemplates returns a list of templates depending on a supplied comparison function
func filterTemplates(templates []*commands.Template, templateName string, filterFunction func(*commands.Template, string) bool) ([]*commands.Template, error) {
	var match []*commands.Template
	for _, v := range templates {
		if filterFunction(v, templateName) {
			match = append(match, v)
		}
	}
	if len(match) == 0 {
		return nil, cli.NewExitError("No matches found", 1)
	}

	return match, nil
}

//NewTemplatesCmd configures the templatesCmd for use with the client
func NewTemplatesCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	var client http.Client
	tc := &templatesCmd{out: out, in: in, getter: client.Do, store: store}
	return tc.Templates()
}

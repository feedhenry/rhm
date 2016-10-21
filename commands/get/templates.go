package get

import (
	"fmt"
	"io"
	"net/http"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/ui"
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

var templatesTemplate = `
| {{PadRight 14 " " "Id"}}| {{PadRight 14 " " "Name"}}| {{PadRight 14 " " "Category"}}|
|-{{PadRight 14 "-" ""}}+-{{PadRight 14 "-" ""}}+-{{PadRight 14 "-" ""}}|{{range .}}
| {{PadRight 14 " " .ID}}| {{PadRight 14 " " .Name}}| {{PadRight 14 " " .Category}}|{{end}}`

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

	res, err := tc.getter(newrequest)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	defer res.Body.Close()
	var dataStructure []*commands.Template
	return ui.NewPrinter(ctx.GlobalString("o"), res.Body, tc.out, templatesTemplate, &dataStructure).Print()
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

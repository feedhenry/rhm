package get

//handle the environments list for rhmap.

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

//environmentsCmd constructs the required writer in order to send the response to the right place.
type environmentsCmd struct {
	out    io.Writer
	in     io.Reader
	getter func(*http.Request) (*http.Response, error)
	store  storage.Storer
}

var environmentListTemplate = `
{{range . }} 
Id      |  {{ .ID}} 
Label   |  {{ .Label}} 
Enabled |  {{ .Enabled}}
Target  |

  -- Id     | {{ .Target.ID}}
  -- Label  | {{ .Target.Label}}
  -- Env    | {{ .Target.Env}}


{{end}}        
`

//Environment Defines our cli command including its flags and usage then returns the command to allow a user to do specific operations on environments
func (ec *environmentsCmd) Environments() cli.Command {
	return cli.Command{
		Name:        "environments",
		Action:      ec.environmentsAction,
		Usage:       "environments",
		Description: "environments:where allows listing environments in rhm",
	}
}

// environmentAction is where the logic is pulled together to perform the command. This funtion conforms to the cli action
func (ec *environmentsCmd) environmentsAction(ctx *cli.Context) error {
	var (
		url = "%s/api/v2/environments/%s"
	)
	userData, err := ec.store.ReadUserData()
	if err != nil {
		return cli.NewExitError("could not read userData "+err.Error(), 1)
	}
	fullURL := fmt.Sprintf(url, userData.Host, "all")

	newrequest, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}

	// create a cookie
	//cookieFeedhenry := http.Cookie{Name: "feedhenry", Value: userData.Auth}
	//newrequest.AddCookie(&cookieFeedhenry)

	newrequest.Header.Add("X-FH-AUTH-USER", "9052dfadcf8994d997e6d1fa7c395d760e448312")

	// do request
	resp, err := ec.getter(newrequest)

	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return cli.NewExitError("failed to read response body "+err.Error(), 1)
	}
	// check if not authed)
	if err := handleEnvironmentsResponseStatus(resp.StatusCode); err != nil {
		ec.out.Write(ret)
		return err
	}

	var resJSON []*commands.Environment

	if err := json.Unmarshal(ret, &resJSON); err != nil {
		return cli.NewExitError("failed to decode response :"+err.Error(), 1)
	}

	t := template.New("Environments")
	t, err = t.Parse(environmentListTemplate)
	//t, _ = t.Parse("{{range . }}Id | {{.ID}} \nLabel | {{ .Label}} \nToken | {{ .Token}} \n\n\t-- Target ID | {{ .Target.ID}} \n\t-- Target Label | {{ .Target.Label}} \n\t-- Target MBaaS Host | {{ .Target.FhMbaasHost}} \n\t-- Target Env | {{ .Target.Env}} \n\n\n{{end}}")
	if err := t.Execute(ec.out, resJSON); err != nil {
		return cli.NewExitError("failed to execute template "+err.Error(), 1)
	}

	return nil
}

// handleEnvironmentsResponseStatus checks whether the API request returned an ok response
func handleEnvironmentsResponseStatus(status int) error {
	if status == http.StatusOK {
		return nil
	}
	return cli.NewExitError(fmt.Sprintf("\n response %d \n", status), 1)
}

// NewEnvironmentsCmd configures the environmentsCmd for use with the client
func NewEnvironmentsCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	var client http.Client
	ec := &environmentsCmd{out: out, in: in, getter: client.Do, store: store}
	return ec.Environments()
}

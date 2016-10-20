package get

//handle the environments list for rhmap.

import (
	"fmt"
	"io"
	"net/http"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/ui"
	"github.com/urfave/cli"
)

//environmentsCmd constructs the required writer in order to send the response to the right place.
type environmentsCmd struct {
	out    io.Writer
	in     io.Reader
	getter func(*http.Request) (*http.Response, error)
	store  storage.Storer
}

var environmentTemplate = `
| {{PadRight 14 " " "Id"}}| {{PadRight 14 " " "Label"}}| {{PadRight 14 " " "Enabled"}}| {{PadRight 14 " " "Target.id"}}| {{PadRight 14 " " "Target.Label"}}| {{PadRight 14 " " "Target.Env"}}|
|-{{PadRight 14 "-" ""}}+-{{PadRight 14 "-" ""}}+-{{PadRight 14 "-" ""}}+-{{PadRight 14 "-" ""}}+-{{PadRight 14 "-" ""}}+-{{PadRight 14 "-" ""}}|{{range . }}
| {{PadRight 14 " " .ID}}| {{PadRight 14 " " .Label}}| {{if .Enabled}}{{PadRight 14 " " "Yes"}}{{else}}{{PadRight 14 " " "No"}}{{end}}| {{PadRight 14 " " .Target.ID}}| {{PadRight 14 " " .Target.Label}}| {{PadRight 14 " " .Target.Env}}|{{end}}
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

	// This ensures that millicore proceeds with the userData.Auth setting
	newrequest.Header.Set("User-Agent", "FHC/Client")

	// create a cookie
	cookie := http.Cookie{Name: "feedhenry", Value: userData.Auth}
	newrequest.AddCookie(&cookie)

	// do request
	res, err := ec.getter(newrequest)
	if err != nil {
		return cli.NewExitError("could not create new request object "+err.Error(), 1)
	}
	defer res.Body.Close()

	var dataStructure []*commands.Environment
	return ui.NewOutPutter(ctx.GlobalString("o"), res.Body, ec.out, environmentTemplate, &dataStructure).Output()

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

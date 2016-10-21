package commands

//login handle the login logic for rhmap.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/feedhenry/rhm/request"
	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/ui"
	"github.com/urfave/cli"
)

//LoginCmd constructs the required writer in order to send the response to the right place.
type loginCmd struct {
	out      io.Writer
	in       io.Reader
	poster   func(string, string, io.Reader) (*http.Response, error)
	host     string
	username string
	password string
	store    storage.Storer
}

//Login Defines our cli command including its flags and usage then returns the command to allow a user to login
func (lc *loginCmd) Login() cli.Command {
	return cli.Command{
		Name:        "login",
		Action:      lc.loginAction,
		Usage:       "login <full_host>",
		Description: "login will authenticate you against the <host>. If no host is provided it will attempt to authenticate you against localhost",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "username",
				Destination: &lc.username,
				Usage:       "set username from cli  --username=<test@test.com>",
			},
			cli.StringFlag{
				Name:        "password",
				Destination: &lc.password,
				Usage:       "set password from cli  --password=<mypass>",
			},
		},
	}
}

//this is the data structure for posting to the server
type loginParams struct {
	UserName string `json:"u"`
	Password string `json:"p"`
	Domain   string `json:"d"`
}

//loginAction is where the logic is pulled together to perform the command. This funtion conforms to the cli action
func (lc *loginCmd) loginAction(ctx *cli.Context) error {
	var (
		url = "%s/box/srv/1.1/act/sys/auth/login"
	)
	if len(ctx.Args()) != 1 {
		return cli.NewExitError("missing argument: "+ctx.Command.Usage, 1)
	}
	lc.host = ctx.Args()[0]
	login, err := lc.getUsernamePassword()
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	domain, err := extractDomainFromHost(lc.host)
	if err != nil {
		return cli.NewExitError("error extracting domain "+err.Error(), 1)
	}
	login.Domain = domain

	postData, err := request.PrepareJSONBody(login)
	if err != nil {
		return cli.NewExitError("error preparing json "+err.Error(), 1)
	}
	fullURL := fmt.Sprintf(url, lc.host)
	res, err := lc.poster(fullURL, "application/json", postData)
	if err != nil {
		return cli.NewExitError("login request failed "+err.Error(), 1)
	}
	defer res.Body.Close()
	//no need for full type here map is fine
	resJSON := make(map[string]interface{})
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&resJSON); err != nil {
		return cli.NewExitError("failed to decode response", 1)
	}
	if !loginSuccess(resJSON) {
		return cli.NewExitError("authentication failed ", 1)
	}
	//store our data locally for use with other commands
	userData := storage.NewUserData(getFeedHenryCookie(res.Cookies()), login.UserName, lc.host, domain)
	if err := lc.store.WriteUserData(userData); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

//helper for prompting for user pass or getting it from the flags
func (lc *loginCmd) getUsernamePassword() (loginParams, error) {
	login := loginParams{}
	if lc.username == "" {
		username, err := ui.WaitForAnswer("Enter your username", lc.out, lc.in)
		if err != nil {
			return login, err
		}
		login.UserName = username
	} else {
		login.UserName = lc.username
	}
	if lc.password == "" {
		password, err := ui.WaitForAnswer("Enter your password", lc.out, lc.in)
		if err != nil {
			return login, err
		}
		login.Password = password
	} else {
		login.Password = lc.password
	}
	return login, nil
}

//frustrating but the api always return 200
func loginSuccess(res map[string]interface{}) bool {
	if v, ok := res["result"]; ok {
		return v.(string) == "ok"
	}
	return false
}

//helper to get the cookie value
func getFeedHenryCookie(cookies []*http.Cookie) string {
	for _, c := range cookies {
		if c.Name == "feedhenry" {
			return c.Value
		}
	}
	return ""
}

//helper to get the rhmap domain
func extractDomainFromHost(host string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}
	hostParts := strings.Split(u.Host, ".")
	return hostParts[0], nil
}

//NewLoginCmd configures the LoginCmd for use with the client
func NewLoginCmd(in io.Reader, out io.Writer, store storage.Storer) cli.Command {
	lc := &loginCmd{out: out, in: in, host: "http://localhost", poster: http.Post, store: store}
	return lc.Login()
}

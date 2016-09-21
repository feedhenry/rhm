package commands

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

//maybe should be in shared place
func createMockPoster(t *testing.T, responseStatus int, path, response string) func(string, string, io.Reader) (*http.Response, error) {
	return func(api, contentType string, body io.Reader) (*http.Response, error) {
		u, err := url.Parse(api)
		if err != nil {
			t.Fatal(err.Error())
		}
		if u.Path != "/box/srv/1.1/act/sys/auth/login" {
			t.Fatal("incorrect api path")
		}
		resBody := response
		bodyRC := ioutil.NopCloser(bytes.NewReader([]byte(resBody)))
		headers := http.Header{}
		headers.Add("Set-Cookie", "feedhenry=test;")
		res := &http.Response{StatusCode: 200, Body: bodyRC, Header: headers}
		return res, nil
	}
}

//maybe should be put in shared place
type MockLoginStore struct {
	writeAssert func(ud *storage.UserData)
}

func (ml MockLoginStore) ReadUserData() (*storage.UserData, error) {
	return nil, nil
}
func (ml MockLoginStore) WriteUserData(ud *storage.UserData) error {
	ml.writeAssert(ud)
	return nil
}

func TestLoginActionOk(t *testing.T) {
	var reader bytes.Buffer
	var writer bytes.Buffer
	poster := createMockPoster(t, 200, "/box/srv/1.1/act/sys/auth/login", `{"result":"ok"}`)
	store := MockLoginStore{
		writeAssert: func(ud *storage.UserData) {
			if ud.UserName != "test@test.com" {
				t.Fatal("expected the username to match")
			}
			if err := ud.Validate(); err != nil {
				t.Fatal(err.Error())
			}
		},
	}
	lcmd := &loginCmd{out: &writer, in: &reader, host: "http://localhost", poster: poster, store: store}
	lcmd.username = "test@test.com"
	lcmd.password = "password"
	if err := lcmd.loginAction(&cli.Context{}); err != nil {
		t.Fatalf("expected login to succeed %s", err.Error())
	}
	if len(writer.Bytes()) > 0 {
		//something was written to the writer not successful
		t.Fatalf("didnt expect output %s", string(writer.Bytes()))
	}
}

func TestLoginActionFail(t *testing.T) {
	var reader bytes.Buffer
	var writer bytes.Buffer
	poster := createMockPoster(t, 200, "/box/srv/1.1/act/sys/auth/login", `{"result":"fail"}`)
	//dont need the store should never get there
	lcmd := &loginCmd{out: &writer, in: &reader, host: "http://localhost", poster: poster, store: nil}
	lcmd.username = "test@test.com"
	lcmd.password = "password"
	err := lcmd.loginAction(&cli.Context{})
	if err == nil {
		t.Fatalf("expected login to fail ")
	}
	t.Log(err.Error())
}

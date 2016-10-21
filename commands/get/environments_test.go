package get

import (
	"bytes"
	"strings"
	"testing"

	"flag"

	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
	"github.com/urfave/cli"
)

var outPutType string

func TestEnvironmentsAction(t *testing.T) {
	cxt := cli.NewContext(cli.NewApp(), flag.NewFlagSet("o", 0), nil)
	cxt.GlobalSet("o", "json")
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `[{"id": "some-env", "label": "some-env-do-not-delete" , "target" : {"id" : "targetId" , "label": "targetLabel"} }]`
	eCmd := &environmentsCmd{
		in:     &in,
		out:    &out,
		store:  mockStore,
		getter: mock.CreateRequest(t, 200, "testing.feedhenry.me/api/v2/environments/all", response),
	}
	if err := eCmd.environmentsAction(cxt); err != nil {
		t.Fatal("failed to exectute environments cmd" + err.Error())
	}
	content := string(out.Bytes())
	if !strings.Contains(content, "some-env") {
		t.Fatal("expected the environment id to be present")
	}
	if !strings.Contains(content, "some-env-do-not-delete") {
		t.Fatal("expected the environment label to be present")
	}
	if !strings.Contains(content, "targetId") {
		t.Fatal("expected the environment target id to be present")
	}
	if !strings.Contains(content, "targetLabel") {
		t.Fatal("expected the environment target label to be present")
	}
}

func TestEnvironmentsAction401Error(t *testing.T) {
	cxt := cli.NewContext(cli.NewApp(), flag.NewFlagSet("o", 0), nil)
	cxt.GlobalSet("o", "json")
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `{"status":"error"}`
	eCmd := &environmentsCmd{
		in:     &in,
		out:    &out,
		store:  mockStore,
		getter: mock.CreateRequest(t, 401, "testing.feedhenry.me/api/v2/environments/all", response),
	}
	if err := eCmd.environmentsAction(cxt); err == nil {
		t.Fatal("expected an error executing command")
	}
}

func TestProjectsTestEnvironmentsAction401Error(t *testing.T) {
	cxt := cli.NewContext(cli.NewApp(), flag.NewFlagSet("o", 0), nil)
	cxt.GlobalSet("o", "json")
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `{"status":"error"}`
	eCmd := &environmentsCmd{
		in:     &in,
		out:    &out,
		store:  mockStore,
		getter: mock.CreateRequest(t, 500, "testing.feedhenry.me/api/v2/environments/all", response),
	}
	if err := eCmd.environmentsAction(cxt); err == nil {
		t.Fatal("expected an error executing command ")
	}
}

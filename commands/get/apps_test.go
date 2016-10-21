package get

import (
	"bytes"
	"errors"
	"flag"
	"strings"
	"testing"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
	"github.com/urfave/cli"
)

func mockAppsActionProjectFinder(title string, userData *storage.UserData, getter commands.HTTPGetter) (string, error) {
	return "", errors.New("Project not found")
}

func TestAppsAction(t *testing.T) {
	var (
		in, out   bytes.Buffer
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
	)
	//setup the flags to be passed through
	fSet := new(flag.FlagSet)
	ctx := cli.NewContext(nil, fSet, nil)
	t.Run("200ok", func(t *testing.T) {
		mockResponse := `{
			"title": "cordova-test",
			"guid": "scqswfv56m7fktyijkfw6tkd",
			"apps": [
			{
				"title": "app",
				"guid": "c36tnuxw4emjxkrxhkbgtg4x",
				"scmUrl": "git@gitlab-shell:rhmap/phils2-app.git",
				"scmKey": null,
				"scmBranch": "master"
			}]
		}`
		getter := mock.CreateRequest(t, 200, "testing.feedhenry.me/box/api/projects/scqswfv56m7fktyijkfw6tkd", mockResponse)
		mockStore.Data.ActiveProject = "scqswfv56m7fktyijkfw6tkd"
		aCommand := appsCmd{
			in:            &in,
			out:           &out,
			store:         mockStore,
			getter:        getter,
			projectFinder: mockAppsActionProjectFinder,
		}

		if err := aCommand.appsAction(ctx); err != nil {
			t.Fatal("did not expect an error ", err.Error())
		}
		content := string(out.Bytes())
		if !strings.Contains(content, "c36tnuxw4emjxkrxhkbgtg4x") {
			t.Fatalf("expected to find the app guid in the output")
		}
	})

	t.Run("500fail", func(t *testing.T) {
		mockResponse := `{"status": "error", "message": "unexpected error"}`
		mockStore.Data.ActiveProject = "scqswfv56m7fktyijkfw6tkd"
		getter := mock.CreateRequest(t, 500, "testing.feedhenry.me/box/api/projects/scqswfv56m7fktyijkfw6tkd", mockResponse)
		aCommand := appsCmd{
			in:            &in,
			out:           &out,
			store:         mockStore,
			getter:        getter,
			projectFinder: mockAppsActionProjectFinder,
		}
		if err := aCommand.appsAction(ctx); err == nil {
			t.Fatal("expected an error got nil")
		}
	})

	t.Run("401", func(t *testing.T) {
		mockResponse := `{"status": "error", "message": "unexpected error"}`
		mockStore.Data.ActiveProject = "scqswfv56m7fktyijkfw6tkd"
		getter := mock.CreateRequest(t, 401, "testing.feedhenry.me/box/api/projects/scqswfv56m7fktyijkfw6tkd", mockResponse)
		aCommand := appsCmd{
			in:            &in,
			out:           &out,
			store:         mockStore,
			getter:        getter,
			projectFinder: mockAppsActionProjectFinder,
		}
		if err := aCommand.appsAction(ctx); err == nil {
			t.Fatal("expected an error got nil")
		}
	})
}

func TestApps(t *testing.T) {
	appsCmd := appsCmd{}
	cliCmd := appsCmd.Apps()
	if cliCmd.Name != "apps" {
		t.Fatal("expected cliCmd.Name to be: 'apps' but got: " + cliCmd.Name)
	}
}

func TestNewAppsCmd(t *testing.T) {
	var (
		in, out   bytes.Buffer
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
	)
	aCommand := NewAppsCmd(&in, &out, mockStore)
	if aCommand.Name != "apps" {
		t.Fatal("expected aCommand.Name to be: 'apps' but got: " + aCommand.Name)
	}
}

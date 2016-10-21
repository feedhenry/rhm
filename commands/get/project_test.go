package get

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	"errors"

	"github.com/feedhenry/rhm/commands"
	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
	"github.com/urfave/cli"
)

func mockProjectFinder(title string, userData *storage.UserData, getter commands.HttpGetter) (string, error) {
	return "", errors.New("Project not found")
}

func TestProjectAction(t *testing.T) {
	var (
		in, out   bytes.Buffer
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
	)
	//setup the flags to be passed through
	fSet := new(flag.FlagSet)
	ctx := cli.NewContext(nil, fSet, nil)
	t.Run("200ok", func(t *testing.T) {
		mockResponse := `{"title": "cordova-test", "guid": "scqswfv56m7fktyijkfw6tkd"}`
		getter := mock.CreateRequest(t, 200, "testing.feedhenry.me/box/api/projects/scqswfv56m7fktyijkfw6tkd", mockResponse)
		pCommand := projectCmd{
			in:            &in,
			out:           &out,
			store:         mockStore,
			getter:        getter,
			project:       "scqswfv56m7fktyijkfw6tkd",
			projectFinder: mockProjectFinder,
		}

		if err := pCommand.projectAction(ctx); err != nil {
			t.Fatal("did not expect an error ", err.Error())
		}
		content := string(out.Bytes())
		if !strings.Contains(content, "scqswfv56m7fktyijkfw6tkd") {
			t.Fatalf("expected to find the guid in the output")
		}
	})

	t.Run("500fail", func(t *testing.T) {
		mockResponse := `{"status": "error", "message": "unexpected error"}`
		getter := mock.CreateRequest(t, 500, "testing.feedhenry.me/box/api/projects/scqswfv56m7fktyijkfw6tkd", mockResponse)
		pCommand := projectCmd{
			in:            &in,
			out:           &out,
			store:         mockStore,
			getter:        getter,
			project:       "scqswfv56m7fktyijkfw6tkd",
			projectFinder: mockProjectFinder,
		}
		if err := pCommand.projectAction(ctx); err == nil {
			t.Fatal("expected an error ", err.Error())
		}
	})

	t.Run("401", func(t *testing.T) {
		mockResponse := `{"status": "error", "message": "unexpected error"}`
		getter := mock.CreateRequest(t, 401, "testing.feedhenry.me/box/api/projects/scqswfv56m7fktyijkfw6tkd", mockResponse)
		pCommand := projectCmd{
			in:            &in,
			out:           &out,
			store:         mockStore,
			getter:        getter,
			project:       "scqswfv56m7fktyijkfw6tkd",
			projectFinder: mockProjectFinder,
		}
		if err := pCommand.projectAction(ctx); err == nil {
			t.Fatal("expected an error ", err.Error())
		}
	})

}

func TestProjectNameToGuid(t *testing.T) {
	mockResponse := `[{"guid": "347bkfnjoemm6cunjr2fbb6w", "title": "project_name"}]`
	getter := mock.CreateRequest(t, 200, "testing.feedhenry.me/box/api/projects", mockResponse)

	guid, err := ProjectNameToGuid("project_name", storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"), getter)
	if err != nil {
		t.Fatal("unexpected error: " + err.Error())
	}

	if guid != "347bkfnjoemm6cunjr2fbb6w" {
		t.Fatal("expected guid: 347bkfnjoemm6cunjr2fbb6w got: " + guid)
	}
}

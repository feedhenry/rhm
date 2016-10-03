package get

import (
	"bytes"
	"strings"
	"testing"

	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
)

func TestProjectsAction(t *testing.T) {
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `[{"title": "cordova-test", "guid": "scqswfv56m7fktyijkfw6tkd"}]`
	pCmd := &projectsCmd{
		in:     &in,
		out:    &out,
		store:  mockStore,
		getter: mock.CreateMockProjectGetter(t, 200, "testing.feedhenry.me/box/api/projects", response),
	}
	if err := pCmd.projectsAction(nil); err != nil {
		t.Fatal("failed to exectute projects cmd" + err.Error())
	}
	content := string(out.Bytes())
	if !strings.Contains(content, "cordova-test") {
		t.Fatal("expected the project title to be present")
	}
	if !strings.Contains(content, "scqswfv56m7fktyijkfw6tkd") {
		t.Fatal("expected the project guid to be present")
	}
}

func TestProjectsAction401Error(t *testing.T) {
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `{"status":"error"}`
	pCmd := &projectsCmd{
		in:     &in,
		out:    &out,
		store:  mockStore,
		getter: mock.CreateMockProjectGetter(t, 401, "testing.feedhenry.me/box/api/projects", response),
	}
	if err := pCmd.projectsAction(nil); err == nil {
		t.Fatal("expected an error executing command")
	}
}

func TestProjectsAction500Error(t *testing.T) {
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `{"status":"error"}`
	pCmd := &projectsCmd{
		in:     &in,
		out:    &out,
		store:  mockStore,
		getter: mock.CreateMockProjectGetter(t, 500, "testing.feedhenry.me/box/api/projects", response),
	}
	if err := pCmd.projectsAction(nil); err == nil {
		t.Fatal("expected an error executing command ")
	}
}

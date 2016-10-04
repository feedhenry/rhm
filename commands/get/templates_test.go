package get

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
)

func TestTemplatesAction(t *testing.T) {
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me/", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `[{"name": "example-template", "id": "scqswfv56m7fktyijkfw6tkd", "category": "this is a category"}]`
	tCmd := &templatesCmd{
		in:           &in,
		out:          &out,
		store:        mockStore,
		templateType: "projects",
		getter:       mock.CreateRequest(t, 200, "testing.feedhenry.me/box/api/templates/projects", response),
	}
	if err := tCmd.templatesAction(nil); err != nil {
		t.Fatal("failed to exectute templates cmd" + err.Error())
	}
	content := string(out.Bytes())
	fmt.Printf(content)
	if !strings.Contains(content, "example-template") {
		t.Fatal("expected the template title to be present")
	}
	if !strings.Contains(content, "scqswfv56m7fktyijkfw6tkd") {
		t.Fatal("expected the template guid to be present")
	}
	if !strings.Contains(content, "this is a category") {
		t.Fatal("expected the template description to be present")
	}
}

func TestTemplatesAction401Error(t *testing.T) {
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me/", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `{"status":"error"}`
	tCmd := &templatesCmd{
		in:           &in,
		out:          &out,
		store:        mockStore,
		templateType: "projects",
		getter:       mock.CreateRequest(t, 401, "testing.feedhenry.me/box/api/templates/projects", response),
	}
	if err := tCmd.templatesAction(nil); err == nil {
		t.Fatal("expected an error executing command")
	}
}

func TestTemplatesAction500Error(t *testing.T) {
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me/", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `{"status":"error"}`
	tCmd := &templatesCmd{
		in:           &in,
		out:          &out,
		store:        mockStore,
		templateType: "projects",
		getter:       mock.CreateRequest(t, 500, "testing.feedhenry.me/box/api/templates/projects", response),
	}
	if err := tCmd.templatesAction(nil); err == nil {
		t.Fatal("expected an error executing command ")
	}
}

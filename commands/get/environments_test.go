package get

import (
	"bytes"
	"strings"
	"testing"

	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
)

func TestEnvironmentsAction(t *testing.T) {
	var (
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
		in        bytes.Buffer
		out       bytes.Buffer
	)
	response := `[{"title": "cordova-test", "guid": "scqswfv56m7fktyijkfw6tkd"}]`
	eCmd := &environmentsCmd{
		in:     &in,
		out:    &out,
		store:  mockStore,
		getter: mock.CreateRequest(t, 200, "testing.feedhenry.me/api/v2/environments/all", response),
	}
	if err := eCmd.environmentsAction(nil); err != nil {
		t.Fatal("failed to exectute environments cmd" + err.Error())
	}
	content := string(out.Bytes())
	if !strings.Contains(content, "cordova-test") {
		t.Fatal("expected the environment title to be present")
	}
	if !strings.Contains(content, "scqswfv56m7fktyijkfw6tkd") {
		t.Fatal("expected the environment guid to be present")
	}
}

func TestEnvironmentsAction401Error(t *testing.T) {
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
	if err := eCmd.environmentsAction(nil); err == nil {
		t.Fatal("expected an error executing command")
	}
}

func TestProjectsTestEnvironmentsAction401Error(t *testing.T) {
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
	if err := eCmd.environmentsAction(nil); err == nil {
		t.Fatal("expected an error executing command ")
	}
}

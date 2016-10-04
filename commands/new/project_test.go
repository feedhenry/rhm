package new

import (
	"bytes"
	"strings"
	"testing"

	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
)

func TestProjectAction(t *testing.T) {
	var (
		in, out   bytes.Buffer
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
	)
	mockResponse := `{"title": "test", "guid": "someGuid"}`
	getter := mock.CreateMockProjectCreate(t, 200, "testing.feedhenry.me/box/api/projects", mockResponse)
	pCommand := projectCmd{
		in:         &in,
		out:        &out,
		store:      mockStore,
		httpClient: getter,
		title:      "test",
	}
	if err := pCommand.projectAction(nil); err != nil {
		t.Fatal("did not expect an error ", err.Error())
	}
	content := string(out.Bytes())
	if !strings.Contains(content, "test") {
		t.Fatalf("expected to find project title in the output")
	}
}

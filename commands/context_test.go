package commands

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
)

func TestContextAction(t *testing.T) {
	var (
		in          bytes.Buffer
		out         bytes.Buffer
		projectGUID = "ckexemysczrrv4qbc2qwankj"
		domain      = "mydomain"
		host        = "host.test.com"
		userName    = "test@test.com"
	)
	ud := &storage.UserData{
		ActiveProject: projectGUID,
		Domain:        domain,
		Host:          host,
		UserName:      userName,
	}
	store := mock.UserDataStore(ud)
	contextCmd := contextCmd{in: &in, out: &out, store: store}
	if err := contextCmd.contextAction(nil); err != nil {
		t.Fatal("did not expect error from context command ", err.Error())
	}
	outPut := string(out.Bytes())
	contains := []string{projectGUID, domain, host, userName}
	for _, v := range contains {
		t.Run("test output contains "+v, func(t *testing.T) {
			if !strings.Contains(outPut, v) {
				t.Fatal("expected to find ", v, " in the output")
			}
		})
	}
}

func TestContextActionError(t *testing.T) {
	store := mock.UserDataStore(nil)
	store.ReadError = errors.New("failed to read data")
	contextCmd := contextCmd{in: nil, out: nil, store: store}
	if err := contextCmd.contextAction(nil); err == nil {
		t.Fatal("Expected an error from context command ")
	}
}

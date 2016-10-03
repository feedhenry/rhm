package delete

import (
	"bytes"
	"flag"
	"net/http"
	"strings"
	"testing"

	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
	"github.com/urfave/cli"
)

func TestDeleteStopsOnNegativeAnswer(t *testing.T) {
	var (
		in        bytes.Buffer
		out       bytes.Buffer
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
	)
	pd := &projectDelete{
		in:          &in,
		out:         &out,
		projectGUID: "myguid",
		store:       mockStore,
	}
	in.Write([]byte("no\n")) //answer no to the prompt
	err := pd.deleteAction(nil)
	if err == nil {
		t.Fatal("expected an error to be returned")
	}
	if !strings.Contains(err.Error(), "not deleted") {
		t.Fatal("expected project not to be deleted")
	}
}

func TestDeleteOkOnYes(t *testing.T) {
	var (
		in        bytes.Buffer
		out       bytes.Buffer
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
	)
	pd := &projectDelete{
		in:          &in,
		out:         &out,
		projectGUID: "myguid",
		store:       mockStore,
		deleter:     mock.CreateRequest(t, http.StatusOK, "testing.feedhenry.me/box/api/projects/myguid", ""),
	}
	in.Write([]byte("yes\n")) //answer yes to the prompt
	err := pd.deleteAction(nil)
	if err != nil {
		t.Fatal("did not expect an error to be returned")
	}
}

func TestDeleteOkOnForce(t *testing.T) {
	var (
		in        bytes.Buffer
		out       bytes.Buffer
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
	)
	pd := &projectDelete{
		in:          &in,
		out:         &out,
		projectGUID: "myguid",
		store:       mockStore,
		deleter:     mock.CreateRequest(t, http.StatusOK, "testing.feedhenry.me/box/api/projects/myguid", ""),
		force:       false, //set force to false override with flag later
	}

	//setup the flags to be passed through
	fSet := new(flag.FlagSet)
	fSet.BoolVar(&pd.force, "f", false, "")
	fSet.Parse([]string{"-f", "true"})
	ctx := cli.NewContext(nil, fSet, nil)
	err := pd.deleteAction(ctx)
	if err != nil {
		t.Fatal("did not expect an error to be returned")
	}
}

func TestDeleteOkWithProjectFlag(t *testing.T) {
	var (
		in        bytes.Buffer
		out       bytes.Buffer
		mockStore = mock.UserDataStore(storage.NewUserData("test", "test@test.com", "testing.feedhenry.me", "testing"))
	)
	pd := &projectDelete{
		in:      &in,
		out:     &out,
		store:   mockStore,
		deleter: mock.CreateRequest(t, http.StatusOK, "testing.feedhenry.me/box/api/projects/myotherguid", ""),
	}
	//setup the flags to be passed through
	fSet := new(flag.FlagSet)
	fSet.StringVar(&pd.projectGUID, "project", "", "")
	fSet.Parse([]string{"--project", "myotherguid"})
	ctx := cli.NewContext(nil, fSet, nil)

	in.Write([]byte("yes\n")) //answer yes to the prompt
	err := pd.deleteAction(ctx)
	if err != nil {
		t.Fatal("did not expect an error to be returned")
	}
}

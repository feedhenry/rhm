package use

import (
	"flag"
	"testing"

	"github.com/feedhenry/rhm/storage"
	"github.com/feedhenry/rhm/test/mock"
	"github.com/urfave/cli"
)

func TestProjectAction(t *testing.T) {
	var projectGUID = "g6czgkxts27apu35nj6pztqm"
	t.Run("fails when not enough args", func(t *testing.T) {
		cmd := setProjectCmd{in: nil, out: nil, storage: mock.UserDataStore(nil)}
		var flagSet flag.FlagSet
		fset := &flagSet
		fset.Parse([]string{})
		ctx := cli.NewContext(nil, fset, nil)
		if err := cmd.SetProjectAction(ctx); err == nil {
			t.Fatal("expected an error when not enough args passed")
		}
	})

	t.Run("saves data to store ", func(t *testing.T) {
		ud := &storage.UserData{
			ActiveProject: projectGUID,
		}
		var store = mock.UserDataStore(ud)
		var writeCalled = 0
		store.WriteAssert = func(d *storage.UserData) {
			writeCalled++
			if d.ActiveProject != projectGUID {
				t.Fatal("expected the project guid to match")
			}
		}
		cmd := setProjectCmd{in: nil, out: nil, storage: store}
		var flagSet flag.FlagSet
		fset := &flagSet
		fset.Parse([]string{projectGUID})
		ctx := cli.NewContext(nil, fset, nil)
		if err := cmd.SetProjectAction(ctx); err != nil {
			t.Fatal("expected no error when setting project")
		}
		if writeCalled != 1 {
			t.Fatal("expected write to be called once")
		}
	})
}

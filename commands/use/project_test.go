package use

import (
	"flag"
	"testing"

	"github.com/feedhenry/rhm/storage"
	"github.com/urfave/cli"
)

type mockStore struct {
	writeAssert func(*storage.UserData)
	data        *storage.UserData
}

func (ms mockStore) WriteUserData(ud *storage.UserData) error {
	ms.writeAssert(ud)
	return nil
}

func (ms mockStore) ReadUserData() (*storage.UserData, error) {
	return ms.data, nil
}

func projectMockStore(toReturn *storage.UserData) mockStore {
	return mockStore{data: toReturn}
}

func TestProjectAction(t *testing.T) {
	var projectGUID = "g6czgkxts27apu35nj6pztqm"
	t.Run("fails when not enough args", func(t *testing.T) {
		cmd := setProjectCmd{in: nil, out: nil, storage: projectMockStore(nil)}
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
		var store = projectMockStore(ud)
		var writeCalled = 0
		store.writeAssert = func(d *storage.UserData) {
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

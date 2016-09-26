package mock

import "github.com/feedhenry/rhm/storage"

//this package holds mocks for tests

//Store implements storer
type Store struct {
	WriteAssert func(*storage.UserData)
	Data        *storage.UserData
	ReadError   error
	WriteError  error
}

func (ms Store) WriteUserData(ud *storage.UserData) error {
	if ms.WriteError != nil {
		return ms.WriteError
	}
	ms.WriteAssert(ud)
	return nil
}

func (ms Store) ReadUserData() (*storage.UserData, error) {
	if ms.ReadError != nil {
		return nil, ms.ReadError
	}
	return ms.Data, nil
}

//UserDataStore mocks out a storer
func UserDataStore(toReturn *storage.UserData) Store {
	return Store{Data: toReturn}
}

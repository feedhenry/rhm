package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/feedhenry/rhm/storage"
)

//this package holds mocks for tests

//Store implements storer
type Store struct {
	WriteAssert func(*storage.UserData)
	Data        *storage.UserData
	ReadError   error
	WriteError  error
}

// WriteUserData will store the UserData to disk
func (ms Store) WriteUserData(ud *storage.UserData) error {
	if ms.WriteError != nil {
		return ms.WriteError
	}
	ms.WriteAssert(ud)
	return nil
}

// ReadUserData retrieves UserData
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

// CreateMockProjectGetter returns a mocked get request
func CreateMockProjectGetter(t *testing.T, responseStatus int, path, response string) func(*http.Request) (*http.Response, error) {
	return func(res *http.Request) (*http.Response, error) {
		if res.URL.Path != path {
			t.Fatal("incorrect api path "+res.URL.Path, res.URL.Path)
		}
		bodyRC := ioutil.NopCloser(bytes.NewReader([]byte(response)))
		resp := &http.Response{StatusCode: responseStatus, Body: bodyRC, Header: http.Header{}}
		return resp, nil
	}
}

// CreateMockProjectCreate returns a mocked get request
func CreateMockProjectCreate(t *testing.T, responseStatus int, path, response string) func(*http.Request) (*http.Response, error) {
	return func(res *http.Request) (*http.Response, error) {
		if res.URL.Path != path {
			t.Fatal("incorrect api path "+res.URL.Path, res.URL.Path)
		}
		bodyRC := ioutil.NopCloser(bytes.NewReader([]byte(response)))
		resp := &http.Response{StatusCode: responseStatus, Body: bodyRC, Header: http.Header{}}
		return resp, nil
	}
}

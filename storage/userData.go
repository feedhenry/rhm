package storage

import (
	"encoding/json"
	"io"
	"os"
	"os/user"

	"github.com/pkg/errors"
)

//UserData represents the current logged in user
type UserData struct {
	Auth          string `json:"auth"`
	Host          string `json:"host"`
	UserName      string `json:"userName"`
	Domain        string `json:"domain"`
	ActiveProject string `json:"activeProject"` //this is the guid of the project the user is currently working with.
}

//NewUserData return new UserData
func NewUserData(auth, user, host, domain string) *UserData {
	return &UserData{
		Auth:     auth,
		UserName: user,
		Host:     host,
		Domain:   domain,
	}
}

//Validate ensures keys and values are present
func (ud *UserData) Validate() error {
	if ud.Auth == "" {
		return errors.New("missing key in user data auth")
	}
	if ud.Host == "" {
		return errors.New("missing key in user data host")
	}
	if ud.UserName == "" {
		return errors.New("missing key in user data userName")
	}
	if ud.Domain == "" {
		return errors.New("missing key in user data domain")
	}
	return nil
}

const (
	//StorageDefaultLocation The default file for storing user data
	StorageDefaultLocation = "/.rhm"
)

func getHomeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "failed to get os user")
	}
	return u.HomeDir, nil
}

//TODO maybe these should just implemt io.Writer and io.Reader?

//Storer defines a user data store
type Storer interface {
	Writer
	Reader
}

//Writer is responsible for writing the UserData to the store
type Writer interface {
	WriteUserData(ud *UserData) error
}

//Reader is responsible for reading the UserData to the store
type Reader interface {
	ReadUserData() (*UserData, error)
}

//Store implements the Storer interface
type Store struct{}

//ReadUserData reads the users data from the default location on disk
func (s Store) ReadUserData() (*UserData, error) {
	f, err := openUserData()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	userData := &UserData{}
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(userData); err != nil {
		return nil, errors.Wrap(err, "failed to decode user data")
	}
	return userData, nil
}

//WriteUserData writes the current user data to the default location on disk
func (s Store) WriteUserData(ud *UserData) error {
	if err := ud.Validate(); err != nil {
		return err
	}
	f, err := openUserData()
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	if err := encoder.Encode(ud); err != nil {
		return errors.Wrap(err, "failed to write json data to disk ")
	}
	return nil
}

func openUserData() (io.ReadWriteCloser, error) {
	home, err := getHomeDir()
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(home+StorageDefaultLocation, os.O_RDWR|os.O_CREATE, 0655)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open storage ")
	}
	return f, nil
}

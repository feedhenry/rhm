package request

import (
	"bytes"
	"encoding/json"
	"io"
)

//helps with common request actions

//PrepareJSONBody marhals the passed type into json and returns a reader ready to use with http.Post
func PrepareJSONBody(b interface{}) (io.Reader, error) {
	body, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(body), nil
}

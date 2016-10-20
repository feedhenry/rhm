package ui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"strings"
)

// OutPutter handles output types
type OutPutter interface {
	Output() error
}

type Throughput struct {
	from io.ReadCloser
	to   io.Writer
}

type JsonOutPutter struct {
	Throughput
	pretty bool
}

type PlainOutPutter struct {
	Throughput
	template      string
	dataStructure interface{}
}

// NewOutPutter returns a configured OutPutter
func NewOutPutter(format string, from io.ReadCloser, to io.Writer, template string, dataStructure interface{}) OutPutter {
	switch strings.ToLower(format) {
	case "json":
		return &JsonOutPutter{Throughput: Throughput{from: from, to: to}, pretty: true}
	default:
		return &PlainOutPutter{Throughput: Throughput{from: from, to: to}, template: template, dataStructure: dataStructure}
	}
}

// PrettyPrint enables or disables pretty json output
func (j *JsonOutPutter) PrettyPrint(enabled bool) {
	j.pretty = enabled
}

// Output outputs raw json from the reader
func (j *JsonOutPutter) Output() error {
	defer j.from.Close()
	var dest bytes.Buffer
	data, err := ioutil.ReadAll(j.from)
	if err != nil {
		return err
	}
	if j.pretty {
		if err := json.Indent(&dest, data, "", "\t"); err != nil {
			return err
		}
	}
	_, err = j.to.Write(dest.Bytes())
	return err
}

// OutputTemplate takes a template string and outputs it based on the data in the reader. It exepects it to be JSON data
func (p PlainOutPutter) Output() error {
	dec := json.NewDecoder(p.from)
	if err := dec.Decode(p.dataStructure); err != nil {
		fmt.Println(err)
		return err
	}

	funcMap := map[string]interface{}{
		"PadLeft":     PadLeft,
		"PadRight":    PadRight,
		"PadCentered": PadCentered,
	}
	template := template.New("cli output")
	template.Funcs(funcMap)
	template.Parse(p.template)
	if err := template.Execute(p.to, p.dataStructure); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

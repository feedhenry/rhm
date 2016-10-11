package ui

import (
	"bufio"
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
)

//interface helpers help with the user interface

//WaitForAnswer asks the question writing it to the writer and waiting on the reader for a response
func WaitForAnswer(question string, out io.Writer, in io.Reader) (string, error) {
	if _, err := out.Write([]byte(question + ":\n")); err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(in)
	scanner.Scan() //block on reading the reader until it hits the next carriage return
	return scanner.Text(), nil
}

// OutPutter handles output types
type OutPutter struct {
	in  io.ReadCloser
	out io.Writer
}

// NewOutPutter returns a configured OutPutter
func NewOutPutter(in io.ReadCloser, out io.Writer) *OutPutter {
	return &OutPutter{
		in:  in,
		out: out,
	}
}

// OutputJSON outputs raw json from the reader
func (o OutPutter) OutputJSON() error {
	defer o.in.Close()
	var dest bytes.Buffer
	data, err := ioutil.ReadAll(o.in)
	if err != nil {
		return err
	}
	if err := json.Indent(&dest, data, "", "\t"); err != nil {
		return err
	}
	_, err = o.out.Write(dest.Bytes())
	return err
}

// OutputTemplate takes a template string and outputs it based on the data in the reader. It exepects it to be JSON data
func (o OutPutter) OutputTemplate(templateDef string, dataType interface{}) error {
	dec := json.NewDecoder(o.in)
	if err := dec.Decode(dataType); err != nil {
		return err
	}
	t := template.New("cli template")
	t, _ = t.Parse(templateDef)
	if err := t.Execute(o.out, dataType); err != nil {
		return err
	}
	return nil
}

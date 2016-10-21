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

// Printer handles output types
type Printer interface {
	Print() error
}

// JSONPrinter reads from the from source and sends to the to destination in JSON encoding
type JSONPrinter struct {
	from io.ReadCloser
	to   io.Writer
}

// TemplatePrinter reads from the from source, and sends to the destination. The output is determined by the template and the data is read into the
// dataStructure
type TemplatePrinter struct {
	from          io.ReadCloser
	to            io.Writer
	template      string
	dataStructure interface{}
}

// NewPrinter returns a configured Printer
func NewPrinter(format string, from io.ReadCloser, to io.Writer, template string, dataStructure interface{}) Printer {
	switch strings.ToLower(format) {
	case "json":
		return &JSONPrinter{from: from, to: to}
	default:
		return &TemplatePrinter{from: from, to: to, template: template, dataStructure: dataStructure}
	}
}

// Print outputs raw json from the reader
func (j *JSONPrinter) Print() error {
	defer j.from.Close()
	var dest bytes.Buffer
	data, err := ioutil.ReadAll(j.from)
	if err != nil {
		return err
	}
	if err := json.Indent(&dest, data, "", "\t"); err != nil {
		return err
	}
	_, err = j.to.Write(dest.Bytes())
	return err
}

// Print takes a template string and outputs it based on the data in the reader. It exepects it to be JSON data
func (p TemplatePrinter) Print() error {
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

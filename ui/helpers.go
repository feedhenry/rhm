package ui

import (
	"bufio"
	"io"
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

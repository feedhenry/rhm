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

func PadRight(length int, pad, str string) string {
	for i := len(str); i < length; i += len(pad) {
		str += pad
	}
	return str
}

func PadLeft(length int, pad, str string) string {
	for i := len(str); i < length; i += len(pad) {
		str = pad + str
	}
	return str
}

func PadCentered(length int, pad, str string) string {
	for i := len(str); i < length; i += len(pad) {
		if i%2 == 0 {
			str = pad + str
		} else {
			str += pad
		}
	}
	return str
}

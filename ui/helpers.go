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

// PadRight takes a string and appends the pad character to the right of it, until it is as long as length
func PadRight(length int, pad, str string) string {
	for i := len(str); i < length; i += len(pad) {
		str += pad
	}
	return str
}

// PadLeft takes a string and appends the pad character to the left of it, until it is as long as length
func PadLeft(length int, pad, str string) string {
	for i := len(str); i < length; i += len(pad) {
		str = pad + str
	}
	return str
}

// PadCentered takes a string and appends the pad character to the left and right of it, until it is as long as length
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

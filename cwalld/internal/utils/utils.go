package utils

import (
	"errors"
)

type Operation int8

const ( 
	Unknown Operation = iota
	Read
	Write
	ReadWrite
	Metadata
)

type Object struct {
	Name string
	Label string
}

func (o Operation) String() string {
	switch o {
	case Read:
		return "Read"
	case Write:
		return "Write"
	case ReadWrite:
		return "ReadWrite"
	case Metadata:
		return "Metadata"
	}
	return "Unknown"
}

func RegexErr(s []string, regex_type string) (string, error) { // this is to streamline regex error checks and getting the exact variable from the slice
	if s == nil {
		return "", errors.New("Regex failed on " + regex_type)
	}
	return s[1], nil
}

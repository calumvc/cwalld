package utils

import (
	"cwalld/internal/logger"
	"os"
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

func (o Operation) ToString() string {
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

func CheckErr(err error) {
	if err != nil {
		logger.Log(("ERROR:\t" + err.Error()))
		os.Exit(1)
	}
}

func RegexErr(s []string, regex_type string) string {
	if s == nil {
		logger.Log("Regex failed on " + regex_type)
		os.Exit(1)
	}
	return s[1]
}

package utils

import (
	"log"
	"os"
	"os/exec"
)

type Operation int8

const ( 
	Unknown Operation = iota
	Read
	Write
	ReadWrite
	Metadata
)

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
		log.Fatal(err) 
	}
}

func RegexErr(s []string, regex_type string) string {
	if s == nil {
		log.Fatal("Regex failed on ", regex_type)
	}
	return s[1]
}

func GetOS() string {
	os_type := ""
	cmd := "cat /etc/os-release | grep 'Red Hat'"

	out, err := exec.Command("bash", "-c", cmd).Output()

	CheckErr(err)
	if len(out) != 0 {
		os_type = "Red Hat"
	} else {
		os_type = "Arch"
	}

	return os_type
}

func GetArgs() []string {
	return os.Args[1:]
}

package utils

import (
	"fmt"
	"log"
	"os/exec"
)

type Subject struct {
	Pid string 
	Name string
}

type Audit struct {
	Id string
	Subject *Subject
	Directory string
	Operation Operation 
}

type Operation int8

const (
	Unknown Operation = iota
	Read
	Write
	ReadWrite
	Metadata
)

func (s *Subject) ToString() {
	fmt.Printf("New Subject Registered\npid=%s\tcomm=%s\n", s.Pid, s.Name)
}

func (a *Audit) ToString() {
	fmt.Printf("New Audit Registered\nid=%s\tsubject=%s\tpath=%s\topertation=%s\n", a.Id, a.Subject.Pid, a.Directory, a.Operation.ToString())
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

func CheckErr(err error, important bool) {
	if err != nil {
		if important { 
			log.Fatal(err) 
		} else { 
			log.Println(err) 
		}
	}
}

func GetOS() string {
	cmd := exec.Command("cat", "/etc/os-release", "|", "grep", "Red Hat")

	out, err := cmd.CombinedOutput()
	CheckErr(err, true)
	println(out)

	return ""
}

package utils

import (
	"fmt"
	"log"
	"os"
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
	fmt.Printf("New Audit Registered\nid=%s\tsubject=%s\tpath=%s\toperation=%s\n", a.Id, a.Subject.Pid, a.Directory, a.Operation.ToString())
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
	os_type := ""
	cmd := "cat /etc/os-release | grep 'Red Hat'"

	out, err := exec.Command("bash", "-c", cmd).Output()

	CheckErr(err, true)
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

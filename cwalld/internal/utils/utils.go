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
	Object string
	Operation Operation 
}

type Denial struct {
	Subject *Subject
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
	fmt.Printf("New Subject Registered:\tpid=%s\tcomm=%s\n", s.Pid, s.Name)
}

func (a *Audit) ToString() {
	// fmt.Printf("id=%s\tsubject=%s\toperation=%s\tobject=%s\n", a.Id, a.Subject.Name, a.Operation.ToString(), a.Object)
	fmt.Printf("subject=%s\toperation=%s\tobject=%s\n", a.Subject.Name, a.Operation.ToString(), a.Object)
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

func LogDenial(s string, op string, obj string) { // operation here is just text because its reprented in string form by AVC already
	fmt.Printf("<!DENIAL!>:\tsubject=%s\toperation=%s\tobject=%s\n", s, op, obj)
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err) 
	}
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

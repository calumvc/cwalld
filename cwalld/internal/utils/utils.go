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
	Label string
}

type Audit struct {
	Id string
	Subject *Subject
	Object *Object
	Operation Operation 
}

type Object struct {
	Name string
	Label string
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
	fmt.Printf("New Subject Registered:\tpid=%s\tcomm=%s\tlabel=%s\n\n", s.Pid, s.Name, s.Label)
}

func (a *Audit) ToString() {
	// fmt.Printf("id=%s\tsubject=%s\toperation=%s\tobject=%s\n\n", a.Id, a.Subject.Name, a.Operation.ToString(), a.Object)
	// if a.Subject.Name != "setroubleshootd" {
		fmt.Printf("subject=%s : %s\toperation=%s\tobject=%s : %s\n\n", a.Subject.Name, a.Subject.Label, a.Operation.ToString(), a.Object.Name, a.Object.Label)
	// }
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
	fmt.Printf("<!DENIAL!>:\t%s\tattempted { %s }\ton %s\n\n", s, op, obj)
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

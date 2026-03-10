package subject

import (
	"cwalld/internal/utils"
	"fmt"
)

type Subject struct {
	Pid string 
	Name string
	Label string
	Entrypoint string
}

type CoI struct {
	A string
 	B string
}

type Conflicts []CoI

var conflicts = Conflicts{ 
}

func (s *Subject) ToString() {
	fmt.Printf("New Subject Registered:\tpid=%s\tcomm=%s\tlabel=%s\tentrypoint=%s\n\n", s.Pid, s.Name, s.Label, s.Entrypoint)
}

func (s *Subject) AlterLabel(l string, op utils.Operation) {
	fmt.Printf("Considering subject %s { %s } on object %s\n\n", s.Label, op.ToString(), l)
}

func (conflicts Conflicts) inConflict(a string, b string) bool {
	for _, c := range conflicts {
		if (c.A == a && c.B == b) || (c.B == a && c.A == b){
			return true
		}
	}
	return false
}


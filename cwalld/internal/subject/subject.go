package subject

import "fmt"

type Subject struct {
	Pid string 
	Name string
	Label string
}

func (s *Subject) ToString() {
	fmt.Printf("New Subject Registered:\tpid=%s\tcomm=%s\tlabel=%s\n\n", s.Pid, s.Name, s.Label)
}

func (s *Subject) AlterLabel() {

}

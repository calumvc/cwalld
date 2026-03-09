package subject

import "fmt"

type Subject struct {
	Pid string 
	Name string
	Label string
	Entrypoint string
}

func (s *Subject) ToString() {
	fmt.Printf("New Subject Registered:\tpid=%s\tcomm=%s\tlabel=%s\tentrypoint=%s\n\n", s.Pid, s.Name, s.Label, s.Entrypoint)
}

func (s *Subject) AlterLabel(l string) {
	fmt.Printf("Considering subject %s and object %s\n\n", s.Label, l)
}

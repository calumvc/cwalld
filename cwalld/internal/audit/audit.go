package audit

import (
	"cwalld/internal/object"
	"cwalld/internal/subject"
	"cwalld/internal/utils"
	"fmt"
)

type Audit struct {
	Id string
	Subject *subject.Subject
	Object *object.Object
	Operation utils.Operation 
	Success bool
}

func (a *Audit) ToString() {
	fmt.Printf("subject=%s : %s\toperation=%s : %t\tobject=%s : %s\n\n", a.Subject.Name, a.Subject.Label, a.Operation.ToString(), a.Success, a.Object.Name, a.Object.Label)
}

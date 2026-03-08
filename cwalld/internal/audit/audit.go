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
}

func (a *Audit) ToString() {
	fmt.Printf("subject=%s : %s\toperation=%s\tobject=%s : %s\n\n", a.Subject.Name, a.Subject.Label, a.Operation.ToString(), a.Object.Name, a.Object.Label)
}

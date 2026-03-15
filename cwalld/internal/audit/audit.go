package audit

import (
	"cwalld/internal/decorator"
	"cwalld/internal/subject"
	"cwalld/internal/utils"
	"fmt"
)

type Audit struct {
	Id string
	Subject *subject.Subject
	Object *utils.Object
	Operation utils.Operation 
	Success bool
}

func (a *Audit) ToString() {
	line := fmt.Sprintf("subject=%s : %s\toperation=%s : %t\tobject=%s : %s", a.Subject.Name, a.Subject.Label, a.Operation.ToString(), a.Success, a.Object.Name, a.Object.Label)
	decorator.DecorateAndLog(line, "audit")
}

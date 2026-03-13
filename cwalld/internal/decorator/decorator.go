package decorator

import (
	"cwalld/internal/logger"
	"fmt"
)

func DecorateAndLog(s string, log_type string) {
	new_s := ""
	switch log_type {
		case "audit" : {
			new_s = fmt.Sprintf("|Audit:\t%s\n", s)
		}
		case "register" : {
			new_s = fmt.Sprintf("|New Subject:\t%s\n", s)
		}
		case "reregister" : {
			new_s = fmt.Sprintf("|Reregistered:\t%s\n", s)
		}
		case "denial" : {
			new_s = fmt.Sprintf("<!DENIAL!>:\t%s\t<!DENIAL!>\n", s)
		}
		case "relabel" : {
			new_s = fmt.Sprintf("|Relabel: %s\n", s)
		}
		case "relabelcode" : {
			new_s = fmt.Sprintf("Daemon restart successful: %s\n", s)
		}
	}

	logger.Log(new_s)
}

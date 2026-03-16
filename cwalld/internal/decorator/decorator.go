package decorator

import (
	"cwalld/internal/logger"
	"fmt"
)

type Decor int8

const (
	Audit Decor = iota
	Register
	Reregister
	Denial
	Relabel
	Dbus
	Error
)

func DecorateAndLog(s string, d Decor) {
	line := ""
	switch d {
		case Audit: {
			line = fmt.Sprintf("|Audit:\t%s\n", s)
		}
		case Register: {
			line = fmt.Sprintf("|New Subject:\t%s\n", s)
		}
		case Reregister: {
			line = fmt.Sprintf("|Reregistered:\t%s\n", s)
		}
		case Denial: {
			line = fmt.Sprintf("<!DENIAL!>:\t%s\t<!DENIAL!>\n", s)
		}
		case Relabel: {
			line = fmt.Sprintf("|Relabel: %s\n", s)
		}
		case Dbus: {
			line = fmt.Sprintf("Daemon restart successful: %s\n", s)
		}
		case Error: {
			line = fmt.Sprintf("ERROR: %s\n", s)
		}
	default:
		line = "ERROR: Decorator bad argument"
	}

	logger.Log(line)
}

package subject

import (
	"context"
	"cwalld/internal/decorator"
	"cwalld/internal/utils"
	"fmt"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/opencontainers/selinux/go-selinux"
)

type Subject struct {
	Pid string 
	Name string
	Label string
	Entrypoint string
}

func (s *Subject) Log() {
	line := fmt.Sprintf("pid=%s\tcomm=%s\tlabel=%s\tentrypoint=%s", s.Pid, s.Name, s.Label, s.Entrypoint)
	decorator.DecorateAndLog(line, decorator.Register)
}

func (s *Subject) ReLog() {
	line := fmt.Sprintf("%s under label %s", s.Name, s.Label)
	decorator.DecorateAndLog(line, decorator.Reregister)
}

func (s *Subject) AlterLabel(l string, op utils.Operation) error {
	label_change := false
	if s.Label == "unconfined_service_t" || s.Label == "init_t" {
		if op.String() == "Read" || op.String() == "ReadWrite" {
			switch l {
				case "alpha_t" : {
					s.Label = "alpha_rw_exec_t"
					label_change = true
				}
				case "beta_t" : {
					s.Label = "beta_rw_exec_t"
					label_change = true
				}
				case "gamma_t" : {
					s.Label = "gamma_rw_exec_t"
					label_change = true
				}
			}
		}
	}

	if s.Label == "alpha_rw_t" && l == "gamma_t" && (op.String() == "Read" || op.String() == "ReadWrite") {
		label_change = true
		s.Label = "alpha_gamma_r_exec_t"
	}

	if s.Label == "beta_rw_t" && l == "gamma_t" && (op.String() == "Read" || op.String() == "ReadWrite") {
		label_change = true
		s.Label = "beta_gamma_r_exec_t"
	}

	if s.Label == "gamma_rw_t" && l == "alpha_t" && (op.String() == "Read" || op.String() == "ReadWrite") {
		label_change = true
		s.Label = "alpha_gamma_r_exec_t"
	}

	if s.Label == "gamma_rw_t" && l == "beta_t" && (op.String() == "Read" || op.String() == "ReadWrite") {
		label_change = true
		s.Label = "beta_gamma_r_exec_t"
	}

	if label_change {
		err := s.restartSubject()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Subject) restartSubject() error { // subject needs to be restarted to actually get its new label
	label := fmt.Sprintf("system_u:object_r:%s:s0", s.Label)
	line := fmt.Sprintf("attempting: %s to %s", s.Name, s.Label)

	err := selinux.Chcon(s.Entrypoint, label, false)

	if err != nil {
		return err
	}

	decorator.DecorateAndLog(line, decorator.Audit)

	conn, err := dbus.NewSystemConnectionContext(context.Background())

	if err != nil {
		return err
	}

	response_channel := make(chan string, 1)
	conn.RestartUnitContext(context.Background(), fmt.Sprintf("%s.service", s.Name), "replace", response_channel)
	result := <- response_channel
	decorator.DecorateAndLog(result, decorator.Dbus)
	conn.Close()

	return nil
}

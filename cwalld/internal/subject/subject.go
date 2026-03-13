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

func (s *Subject) ToString() {
	line := fmt.Sprintf("pid=%s\tcomm=%s\tlabel=%s\tentrypoint=%s", s.Pid, s.Name, s.Label, s.Entrypoint)
	decorator.DecorateAndLog(line, "register")
}

func (s *Subject) AlterLabel(l string, op utils.Operation) {
	label_change := false
	if s.Label == "unconfined_service_t" || s.Label == "init_t" {
		if op.ToString() == "Read" || op.ToString() == "ReadWrite" {
			label_change = true

			switch l {
				case "alpha_t" : {
					s.Label = "alpha_rw_exec_t"
				}
				case "beta_t" : {
					s.Label = "beta_rw_exec_t"
				}
				case "gamma_t" : {
					s.Label = "gamma_rw_exec_t"
				}
			}
		}
	}

	if s.Label == "alpha_rw_t" && l == "gamma_t" && (op.ToString() == "Read" || op.ToString() == "ReadWrite") {
		label_change = true
		s.Label = "alpha_gamma_r_exec_t"
	}

	if s.Label == "beta_rw_t" && l == "gamma_t" && (op.ToString() == "Read" || op.ToString() == "ReadWrite") {
		label_change = true
		s.Label = "beta_gamma_r_exec_t"
	}

	if s.Label == "gamma_rw_t" && l == "alpha_t" && (op.ToString() == "Read" || op.ToString() == "ReadWrite") {
		label_change = true
		s.Label = "alpha_gamma_r_exec_t"
	}

	if s.Label == "gamma_rw_t" && l == "beta_t" && (op.ToString() == "Read" || op.ToString() == "ReadWrite") {
		label_change = true
		s.Label = "beta_gamma_r_exec_t"
	}

	if label_change {
		s.restart_subject()
	}
}

func (s *Subject) restart_subject() { // subject needs to be restarted to actually get its new label
	label := fmt.Sprintf("system_u:object_r:%s:s0", s.Label)
	line := fmt.Sprintf("%s to %s", s.Name, s.Label)

	err := selinux.Chcon(s.Entrypoint, label, false)
	utils.CheckErr(err)

	decorator.DecorateAndLog(line, "relabel")

	conn, err := dbus.NewSystemConnectionContext(context.Background())
	utils.CheckErr(err)

	response_channel := make(chan string, 1)
	conn.RestartUnitContext(context.Background(), fmt.Sprintf("%s.service", s.Name), "replace", response_channel)
	result := <- response_channel
	decorator.DecorateAndLog(result, "relabelcode")
	conn.Close()
}

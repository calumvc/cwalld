package subject

import (
	"context"
	"cwalld/internal/decorator"
	"cwalld/internal/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/opencontainers/selinux/go-selinux"
)

type Subject struct {
	Pid string 
	Name string
	Label string
	Entrypoint string
}

func (s *Subject) String() string { // log when we find a new subject
	return fmt.Sprintf("pid=%s\tcomm=%s\tlabel=%s\tentrypoint=%s", s.Pid, s.Name, s.Label, s.Entrypoint)
}

func (s *Subject) ReString() string { // relog when we find an old subject with new properties
	return fmt.Sprintf("%s under label %s and new pid %s", s.Name, s.Label, s.Pid)
}

func (s *Subject) AlterLabel(l string, op utils.Operation) error {
	label_change := "false"

	regex := regexp.MustCompile(`r:([^:]+)`)
	regex_label_type := regex.FindStringSubmatch(s.Label)
	label_type, err := utils.RegexErr(regex_label_type, "label type")

	if err != nil {
		return err
	}

	if label_type == "unconfined_service_t" || label_type == "init_t" { // if the subject hasn't been restricted yet

		if op == utils.Read || op == utils.ReadWrite { // if they read from an object, align them with it
			switch l {
				case "alpha_t" : {
					label_change = "alpha_rw_exec_t"
				}
				case "beta_t" : {
					label_change = "beta_rw_exec_t"
				}
				case "gamma_t" : {
					label_change = "gamma_rw_exec_t"
				}
			}
		}
	}

	if label_type == "alpha_rw_t" && l == "gamma_t" && (op == utils.Read || op == utils.ReadWrite) {
		label_change = "alpha_gamma_r_exec_t"
	}

	if label_type == "beta_rw_t" && l == "gamma_t" && (op == utils.Read || op == utils.ReadWrite) {
		label_change = "beta_gamma_r_exec_t"
	}

	if label_type == "gamma_rw_t" && l == "alpha_t" && (op == utils.Read || op == utils.ReadWrite) {
		label_change = "alpha_gamma_r_exec_t"
	}

	if label_type == "gamma_rw_t" && l == "beta_t" && (op == utils.Read || op == utils.ReadWrite) {
		label_change = "beta_gamma_r_exec_t"
	}

	if label_change != "false" {
		err := s.restartSubject(label_change)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Subject) restartSubject(new_label string) error { // subject needs to be restarted to actually get its new label from entrypoint
	l_label := strings.Split(s.Label, "u:")
	r_label := strings.Split(s.Label, "t:")

	label := fmt.Sprintf("%s%s%s", l_label[0] + "u:object_r:", new_label, ":" + r_label[1]) // piece the label back together, it has to have object_r because its an object we're changing
	line := fmt.Sprintf("attempting: %s to %s", s.Name, label)
	
	err := selinux.Chcon(s.Entrypoint, label, false)

	if err != nil {
		return err
	}

	decorator.DecorateAndLog(line, decorator.Audit)

	conn, err := dbus.NewSystemConnectionContext(context.Background())

	if err != nil {
		return err
	}

	response_channel := make(chan string, 1) // function requires a response channel
	conn.RestartUnitContext(context.Background(), fmt.Sprintf("%s.service", s.Name), "replace", response_channel)
	result := <- response_channel
	decorator.DecorateAndLog(result, decorator.Dbus)
	conn.Close()

	return nil
}

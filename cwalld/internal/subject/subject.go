package subject

import (
	"context"
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

var conflicts = [][]int8{
	{1, 0}, // alpha
	{1, 0}, // beta
	{0, 0}, // gamma
}

type faction int8

const (
	alpha faction = iota
	beta
	gamma
)

func (s *Subject) ToString() {
	fmt.Printf("New Subject Registered:\tpid=%s\tcomm=%s\tlabel=%s\tentrypoint=%s\n\n", s.Pid, s.Name, s.Label, s.Entrypoint)
}

func (s *Subject) AlterLabel(l string, op utils.Operation) {
	// fmt.Printf("Considering subject %s { %s } on object %s\n\n", s.Label, op.ToString(), l)
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

func (s *Subject) restart_subject() { // subject needs to be restarted to take the necessary label
	fmt.Printf("attempting to alter %s %s", s.Entrypoint, s.Label)
	label := fmt.Sprintf("system_u:object_r:%s:s0", s.Label)
	fmt.Println("relabelled to %s", label)
	selinux.Chcon(s.Entrypoint, label, false)

	conn, err := dbus.NewSystemConnectionContext(context.Background())
	utils.CheckErr(err)

	response_channel := make(chan string, 1)
	conn.RestartUnitContext(context.Background(), fmt.Sprintf("%s.service", s.Name), "replace", response_channel)
	result := <- response_channel
	println(result)

}

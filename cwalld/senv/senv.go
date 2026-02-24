package senv 

import (
	"cwalld/utils"
	"fmt"
	"os"
	"os/exec"
)

func Setup(DIR string) { // make sure audit is configured
	// TODO: make sure user is using sudo
	// cmd := exec.Command("sudo", "auditctl", "-w", DIR, "-p", "rwa", "-k", "cwalld") // add a rule to auditd to watch all reads and writes and operations and give them a label
	rule := fmt.Sprintf("-D\n-w %s -p rwa -k cwalld", DIR)
	rule_path := "/etc/audit/rules.d/audit.rules"
	
	err := os.WriteFile(rule_path, []byte(rule), 0640)
	utils.CheckErr(err, true)

	cmd := exec.Command("augenrules", "--load") // daemons must be reloaded after rule is added

	err = cmd.Run()
	utils.CheckErr(err, true)

	println("-- Audit Rule Successfully Added --")

	// install selinux modules for the directory next
}


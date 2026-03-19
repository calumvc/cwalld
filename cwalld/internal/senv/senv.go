package senv

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Setup(DIR string) error { // make sure audit is configured
	rule_path := "/etc/audit/rules.d/audit.rules"
	rule := fmt.Sprintf("-D\n-w %s -p rwa -k cwalld", DIR) // this rule will add the key 'cwalld' to every logged event within the specified directory
	
	err := os.WriteFile(rule_path, []byte(rule), 0640)

	if err != nil {
		return err
	}

	cmd := exec.Command("augenrules", "--load") // auditd rules must be refreshed so new daemons follow them

	res, err := cmd.CombinedOutput()
	
	if err != nil {
		return err
	}

	cmd = exec.Command("sudo", "sh", "-c", "echo 0 > /sys/fs/selinux/avc/cache_threshold") // this will let us see repeats in the cache so we can log every single denial
	res, err = cmd.CombinedOutput()

	if bytes.Contains(res, []byte("No such file")) {
		return fmt.Errorf("Error: is selinux enabled?")
	}

	if err != nil {
		return err
	}

	println("-- Audit Rule Successfully Added --")
	return nil
}


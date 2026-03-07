package senv

import (
	"cwalld/internal/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/opencontainers/selinux/go-selinux"
)

func Setup(DIR string) { // make sure audit is configured
	// TODO: make sure user is using sudo

	setupAuditd(DIR)

	labels := []string{}

	scanLabels(DIR, labels)

	for i := range labels {
		println(labels[i])
	}

	// setupSEmodules()
}

func setupAuditd(DIR string){
	rule_path := "/etc/audit/rules.d/audit.rules"
	rule := fmt.Sprintf("-D\n-w %s -p rwa -k cwalld", DIR)
	
	err := os.WriteFile(rule_path, []byte(rule), 0640)
	utils.CheckErr(err)

	cmd := exec.Command("augenrules", "--load") // auditd rules must be refreshed so new daemons follow them

	err = cmd.Run()
	utils.CheckErr(err)

	cmd = exec.Command("sudo", "sh", "-c", "echo 0 > /sys/fs/selinux/avc/cache_threshold") // this will let us see repeats in the cache and therefore stop every single denial
	err = cmd.Run()
	utils.CheckErr(err)

	println("-- Audit Rule Successfully Added --")
}

func scanLabels(DIR string, labels []string) {
	filepath.Walk(DIR, func(file_path string, info os.FileInfo, err error) error {
		res, err := selinux.FileLabel(file_path)

		regex := regexp.MustCompile(`r:([^:]+)`)
		label := regex.FindStringSubmatch(res)
		// fmt.Printf("File %s has label %s\n", file_path, label[1])

		dupe := false
		for i := range labels { 
			if labels[i] == label[1] { 
				dupe = true
			}
		}
		
		if !dupe {
			labels = append(labels, label[1])
		}

		return err
	})
}

func setupSEmodules(){
	sepolicy_path := "/var/lib/cwalld/" // path to write the se policy files to be installed

	err := os.MkdirAll(sepolicy_path, 0755)
	utils.CheckErr(err)
	test := "hellooo"
	
	err = os.WriteFile(sepolicy_path + "test", []byte(test), 0640)
	utils.CheckErr(err)

	println("-- SEmodules Successfully Added -- ")
}

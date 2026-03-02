package senv 

import (
	"cwalld/internal/utils"
	"fmt"
	"os"
	"os/exec"
)

func Setup(DIR string) { // make sure audit is configured
	// TODO: make sure user is using sudo

	setupAuditd(DIR)

	scanCOI(DIR)

	setupSEmodules()
}

func setupAuditd(DIR string){
	rule_path := "/etc/audit/rules.d/audit.rules"
	rule := fmt.Sprintf("-D\n-w %s -p rwa -k cwalld", DIR)
	
	err := os.WriteFile(rule_path, []byte(rule), 0640)
	utils.CheckErr(err, true)

	cmd := exec.Command("augenrules", "--load") // auditd rules must be refreshed so new daemons follow them

	err = cmd.Run()
	utils.CheckErr(err, true)

	println("-- Audit Rule Successfully Added --")
}

func scanCOI(DIR string) {

}

func setupSEmodules(){
	sepolicy_path := "/var/lib/cwalld/"

	err := os.MkdirAll(sepolicy_path, 0755)
	utils.CheckErr(err, true)
	test := "hellooo"
	
	err = os.WriteFile(sepolicy_path + "test", []byte(test), 0640)
	utils.CheckErr(err, true)

	println("-- SEmodules Successfully Added -- ")
}

func activateSEmodules(){

}

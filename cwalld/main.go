package main

import (
	"cwalld/internal/senv"
	"cwalld/internal/sleuth"
	"cwalld/internal/utils"
)

var (
	subjects = []utils.Subject{}
	audits = []utils.Audit{}
)

func main() {
	DIR := "/home/cal/testgrounds/static_wall" // TODO: accept this from cli when I make it

	arg := utils.GetArg()

	if arg == "init" {
		initialize(DIR)
	} else 
	if arg == "enforce" { 
		enforce(DIR)
	} else {
		println("Unsupported argument, try 'init' or 'enforce'")
	}

	// os_type := utils.GetOS();

}

func initialize(DIR string) {
	senv.Setup(DIR)

	println("Chinese Wall Initialized")
}

func enforce(DIR string) {
	println("############## 中國長城 Online ##############")

	go sleuth.TailAuditd(DIR, &subjects, &audits) // follow auditd updates in subprocess

	<-make(chan struct{}) // infinite loop
}

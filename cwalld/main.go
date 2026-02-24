package main

import (
	"cwalld/senv"
	"cwalld/sleuth"
	"cwalld/utils"
)

var (
	subjects = []utils.Subject{}
	audits = []utils.Audit{}
)

func main() {
	println("############## 中國長城 Online ##############")

	DIR := "/home/cal/testgrounds/static_wall" // TODO: accept this from cli when I make it

	senv.Setup(DIR)

	go sleuth.TailAuditd(DIR, &subjects, &audits) // follow auditd updates in subprocess

	<-make(chan struct{}) // infinite loop
}

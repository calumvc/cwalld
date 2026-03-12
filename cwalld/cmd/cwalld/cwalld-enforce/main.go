package main

import "cwalld/internal/sleuth"

func main() {
	DIR := "/home/testgrounds/"

	go sleuth.TailAuditd(DIR) // follow auditd updates in subprocess
	<-make(chan struct{}) // infinite loop
}

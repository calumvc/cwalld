package main

import (
	"cwalld/internal/senv"
	"cwalld/internal/sleuth"
	"cwalld/internal/utils"
)

func main() {
	DIR := "/home/testgrounds/"

	args := utils.GetArgs()

	if args[0] == "init" {
		// initialize(args[1]) // uncomment when i want to accept from cli
		initialize(DIR)
	} else 
	if args[0] == "enforce" { 
		// enforce(args[1])
		enforce(DIR)
	} else {
		println("Unsupported argument, try 'init' or 'enforce'")
	}
}

func initialize(DIR string) {
	senv.Setup(DIR)

	println("Chinese Wall Initialized")
}

func enforce(DIR string) {
	println("############## 中國長城 Online ##############")

	go sleuth.TailAuditd(DIR) // follow auditd updates in subprocess

	<-make(chan struct{}) // infinite loop
}

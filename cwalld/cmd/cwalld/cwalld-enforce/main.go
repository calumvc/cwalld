package main

import (
	"cwalld/internal/decorator"
	"cwalld/internal/sleuth"
)

func main() {
	DIR := "/home/testgrounds/"

	err := sleuth.TailAuditd(DIR) // follow auditd updates in subprocess

	if err != nil {
		decorator.DecorateAndLog(err.Error(), decorator.Error)
	}
}

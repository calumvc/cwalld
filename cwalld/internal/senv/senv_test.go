package senv

import (
	"os/exec"
	"testing"
)

func TestSetup(t *testing.T) {
	cmd := exec.Command("sudo", "setenforce", "1")

	err := cmd.Run()

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	DIR := "/home/testgrounds/"
	err = Setup(DIR);

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

// func TestSetupSetenforceZero(t *testing.T) {
// 	cmd := exec.Command("sudo", "setenforce", "0")
//
// 	err := cmd.Run()
//
// 	if err != nil {
// 		t.Errorf("Error: %s", err.Error())
// 	}
//
// 	DIR := "/home/testgrounds/"
// 	err = Setup(DIR);
//
// 	if err.Error() != "Error: selinux not enforcing" {
// 		t.Errorf("Error handled wrong")
// 	}
// }

func TestSetupBadPath(t *testing.T) {
	cmd := exec.Command("sudo", "setenforce", "1")

	err := cmd.Run()

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	DIR := "/akjhsfasdf/asdfa"
	err = Setup(DIR);

	if err.Error() != "Error: directory /akjhsfasdf/asdfa doesn't exist" {
		t.Errorf("Error handled wrong")
	}
}

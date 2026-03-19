package senv

import (
	"testing"
)

func TestSetup(t *testing.T) {
	DIR := "/home/testgrounds/"
	err := Setup(DIR);

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

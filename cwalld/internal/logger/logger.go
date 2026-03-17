package logger

import (
	"fmt"
	"os"
)

var (
	log_file *os.File
)

func init() { // this runs once before Log ever is called
	var err error // have to make this since log_file variable already exists
	os.MkdirAll("/var/log/cwalld/", 0644)
	log_file, err = os.OpenFile("/var/log/cwalld/cwalld.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0744)

	if err != nil { 
		fmt.Println("Error opening log:", err.Error())
		os.Exit(1)
	}

	fmt.Fprintln(log_file, "Chinese Wall Initialised")
}

func Log(s string) {
	fmt.Fprintln(log_file, s)
}

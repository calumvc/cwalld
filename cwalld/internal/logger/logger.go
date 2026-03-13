package logger

import (
	"fmt"
	"os"
	"sync"
)

var (
	log_file *os.File
	mutex sync.Mutex
)

func init() { // this runs once before Log ever is called
	var err error // have to make this since log_file variable already exists
	os.MkdirAll("/var/log/cwalld/", 0644)
	log_file, err = os.OpenFile("/var/log/cwalld/cwalld.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil { 
		Log(err.Error())
		os.Exit(1)
	}

	mutex.Lock()
	fmt.Fprintln(log_file, "Chinese Wall Initialised")
	mutex.Unlock()
}

func Log(s string) {
	println(s)
	mutex.Lock() // lock and unlock to prevent race conditions
	defer mutex.Unlock()
	fmt.Fprintln(log_file, s)
}

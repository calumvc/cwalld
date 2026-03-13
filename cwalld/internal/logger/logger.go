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
	var err error // have to make this since log_file already exists
	log_file, err = os.OpenFile("/var/lib/cwalld/cwalld.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil { 
		Log(err.Error())
		os.Exit(1)
	}
}

func Log(s string) {
	mutex.Lock() // lock and unlock to prevent race conditions
	defer mutex.Unlock()
	fmt.Fprintln(log_file, s)
}

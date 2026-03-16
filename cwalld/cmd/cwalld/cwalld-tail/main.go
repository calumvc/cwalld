package main

import (
	"fmt"
	"os"

	"github.com/nxadm/tail"
)

func main(){
	t, err := tail.TailFile("/var/log/cwalld/cwalld.log", tail.Config{})

	if err != nil { 
		fmt.Print(err.Error())
		os.Exit(1)
	}

	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

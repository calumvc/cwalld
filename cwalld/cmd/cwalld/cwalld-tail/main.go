package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nxadm/tail"
)

func main() {
	cwalldtail()
}

func cwalldtail() {
	t, err := tail.TailFile("/var/log/cwall/cwall.log", tail.Config{
		Follow: true,
		Location: &tail.SeekInfo{ Offset: 0, Whence: io.SeekStart },
	})

	if err != nil { 
		fmt.Print(err.Error())
		os.Exit(1)
	}

	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

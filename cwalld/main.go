package main

import (
	"log"

	"github.com/hpcloud/tail"
	"k8s.io/utils/inotify"
)

func main() {
	println("############## 中國長城 Online ##############")

	DIR := "/home/cal/testgrounds/static_wall"

	go watch(DIR) // watch directory for changes

	go tail_auditd() // follow auditd updates in subprocess

	<-make(chan struct{})
}

func watch(DIR string) {
	watcher, err := inotify.NewWatcher()
	
	if err != nil { log.Fatal(err) }
	defer watcher.Close()

	go func() {
		for { 
			select { 
			case event := <-watcher.Event:
				log.Println("event:",event)
			}
		}
	}()

	err = watcher.Watch(DIR)
	if err != nil { log.Fatal(err) }

	<-make(chan struct{})
}


func tail_auditd() {
		t, err := tail.TailFile("/var/log/audit/audit.log", tail.Config{Follow: true})
		if err != nil { log.Fatal(err) }

		for line := range t.Lines {
			log.Println(line.Text)
		}
}

package main

import (
	// "fmt"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"

	// "time"

	"strings"

	"github.com/hpcloud/tail"
	"k8s.io/utils/inotify"
)

type subject struct {
	pid string 
	name string
}

type audit struct {
	subject subject
}

var (
	subjects = []subject{}
	audits = []audit{}
)

func main() {
	println("############## 中國長城 Online ##############")

	DIR := "/home/cal/testgrounds/static_wall" // TODO: accept this from cmd when I make it

	setup_environment(DIR)

	go tail_auditd(DIR) // follow auditd updates in subprocess

	go watch_directory(DIR) // watch directory for changes

	// go func() { 
	// 	time.Sleep(50 * time.Second) 
	// }()

	<-make(chan struct{})
}

func watch_directory(DIR string) {
	watcher, err := inotify.NewWatcher()
	check_err(err, true)
	defer watcher.Close()

	println("-- watching --")

	go func() {
		for event := range watcher.Event { 
			log.Println("event:",event)
			fmt.Println()
		}
	}()

	err = watcher.Watch(DIR)
	check_err(err, true)

	<-make(chan struct{})
}

func tail_auditd(DIR string) {
	t, err := tail.TailFile("/var/log/audit/audit.log", tail.Config{ 
		Follow: true,
		Location: &tail.SeekInfo{ Offset: 0, Whence: io.SeekEnd },}) // we only wanna know what happens after we start running the daemon

		println("-- tailing --")

	check_err(err, true)
	
	go func() {
		for line := range t.Lines { // auditd has 2 parts, the syscall and path, we are going to combine them into a struct

			log.Println(line.Text)
			log.Println()

			if strings.Contains(line.Text, "cwalld"){ // this is the syscall part, containing pid and subject name

				regex := regexp.MustCompile(`\bpid=(\d+)`) // regex to catch pid
				pid := regex.FindStringSubmatch(line.Text) // pid[0] = "pid=..." pid[1] = "..."
				
				if len(pid) == 0 { log.Fatal("NO PID FOUND - AUDITD LOG UNPREDICTABLE") }

				regex = regexp.MustCompile(`\bcomm="([^"]+)"`) // regex to catch subject name
				subject_name := regex.FindStringSubmatch(line.Text)

				if len(subject_name) == 0 { log.Fatal("NO SUBJECT NAME FOUND - AUDITD LOG UNPREDICTABLE") }

				already_tracking := false // check we arent already aware of the subject
				for i := range subjects {
					if subjects[i].pid == pid[1] {
						already_tracking = true
					}
				}

				if !already_tracking { // add it to the global list of subjects
					new_subject := subject{ pid: pid[1], name: subject_name[1] }
					new_subject.subject_to_string()
					subjects = append(subjects, new_subject)
				}
			}

			if strings.Contains(line.Text, DIR){ // this is the path part, containing the evidence of the subjects actions

			}
		}
	}()
	
	<-make(chan struct{})
}

func setup_environment(DIR string) { // make sure audit is configured, then install selinux modules
	// TODO: make sure user is using sudo
	cmd := exec.Command("sudo", "auditctl", "-w", DIR, "-p", "rwa", "-k", "cwalld") // add a rule to auditd to watch all reads and writes and operations and give them a label
	
	err := cmd.Run()

	check_err(err, false) // false to say its not worth crashing out over

	cmd = exec.Command("sudo", "systemctl", "daemon-reload") // daemons must be reloaded after rule is added
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	check_err(err, true)

	println("-- Audit Rule Successfully Added --")
}

func (s subject) subject_to_string() {
	fmt.Printf("SUBJECT PID: %s - RUNNING AS %s \n", s.pid, s.name)
}

func check_err(err error, important bool) {
	log_it := true
	if err != nil {
		if important == true {
			log.Fatal(err)

		} else {
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitError.ExitCode() == 255 { // this is for the auditd rule already existing, not a problem and should only happen in development / if the user reboots the daemon without rebooting the system
					log_it = false
				}
			}

			if log_it { 
				log.Println(err)
			}
		}
	}
}

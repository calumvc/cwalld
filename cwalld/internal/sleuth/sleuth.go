package sleuth

import (
	"cwalld/internal/utils"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hpcloud/tail"
	"golang.org/x/sys/unix"
)


func TailAuditd(DIR string) {

	subjects := &[]utils.Subject{}
	audits := &[]utils.Audit{}

	t, err := tail.TailFile("/var/log/audit/audit.log", tail.Config{ 
		Follow: true,
		Location: &tail.SeekInfo{ Offset: 0, Whence: io.SeekEnd },}) // we only wanna know what happens after we start running the daemon

	println("-- tailing --\n")

	utils.CheckErr(err, true)
	
	go func() { // run this part concurrently
		for line := range t.Lines { // auditd has 2 parts, the syscall and path, we are going to combine them into a struct

			if strings.Contains(line.Text, "cwalld"){ // this is the syscall part, containing pid, operation and subject name
				spot_pid(subjects, audits, line.Text)
			}

			if strings.Contains(line.Text, DIR){ // this is the path line, containing the affected object path
				spot_dir(audits, line.Text)
			}

			if strings.Contains(line.Text, "denied") {
				println("Denial INC")
				spot_denial(line.Text)
			}
		}
	}()
	
	<-make(chan struct{})
}

func spot_pid(subjects *[]utils.Subject, audits *[]utils.Audit, line string){
	regex := regexp.MustCompile(`\bpid=(\d+)`) // regex to catch pid
	pid := regex.FindStringSubmatch(line)[1] // pid[0] = "pid=..." pid[1] = "..."
	
	if len(pid) == 0 { log.Fatal("Audit log probably disappeared") }

	regex = regexp.MustCompile(`\bcomm="([^"]+)"`) // regex to catch subject name
	subject_name := regex.FindStringSubmatch(line)[1]

	regex = regexp.MustCompile(`\bmsg=audit\(([^)]+)`) // regex to catch audit id to combine with other line
	audit_id := regex.FindStringSubmatch(line)[1]

	regex = regexp.MustCompile(`\ba2=(\d+)`)
	operation := regex.FindStringSubmatch(line)[1]

	flags, err := strconv.ParseInt(operation, 16, 64) // convert the string that is hexadecimal, into straight binary, which is read as an int64 but actually is just straight flags 
	utils.CheckErr(err, true)

	var op utils.Operation

	if flags & unix.O_RDWR != 0 { // and mask with the O_RDWR flag for both x86 and ARM architecture
		op = utils.ReadWrite
	} else 
	if flags & unix.O_PATH != 0 { // O_PATH is the only operation that declares read but doesnt actually read anything
		op = utils.Metadata
	} else 
	if flags & unix.O_APPEND != 0 || flags & unix.O_TRUNC != 0 || flags & unix.O_CREAT != 0 || flags & unix.O_WRONLY != 0 { // all of the operations that involve writes
		op = utils.Write
	} else { // if the write flag isnt on then it must be a read
		op = utils.Read
	}

	var subject *utils.Subject

	for i := range *subjects { // if subject is already accounted for
		if (*subjects)[i].Pid == pid {
			subject = &(*subjects)[i]
			break
		}
	}

	if subject == nil { // add it to the global list of subjects if not
		subject = &utils.Subject{ Pid: pid, Name: subject_name }
		subject.ToString()
		*subjects = append(*subjects, *subject)
	}

	audit := &utils.Audit{ Id: audit_id, Subject: subject, Directory: "", Operation: op } // create new audit
	*audits = append(*audits, *audit)
}

func spot_dir(audits *[]utils.Audit, line string){
	regex := regexp.MustCompile(`\bmsg=audit\((([^)]+))`)
	audit_id := regex.FindStringSubmatch(line)[1]

	regex = regexp.MustCompile(`\bname="([^"]+)"`)
	directory := regex.FindStringSubmatch(line)[1]

	for _, a := range *audits {
		if a.Id == audit_id {
			a.Directory = directory
			a.ToString()
			break
		}
	}
}

func spot_denial(line string){
	log.Println("DENIAL!!")
	log.Println(line)
	regex := regexp.MustCompile(`\{ ([^ }]+)`)

	operation := regex.FindStringSubmatch(line)[1]
	log.Println(operation)

	regex = regexp.MustCompile(`\bpid=(\d+)`)
	pid := regex.FindStringSubmatch(line)[1]
	log.Println(pid)

	regex = regexp.MustCompile(`\bname="([^"]+)"`)
	object := regex.FindStringSubmatch(line)[1]
	log.Println(object)
	

}

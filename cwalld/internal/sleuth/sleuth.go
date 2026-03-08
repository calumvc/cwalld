package sleuth

import (
	"cwalld/internal/audit"
	"cwalld/internal/object"
	"cwalld/internal/subject"
	"cwalld/internal/utils"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hpcloud/tail"
	"github.com/opencontainers/selinux/go-selinux"
	"golang.org/x/sys/unix"
)

type State struct {
	subjects []subject.Subject
	audits []audit.Audit
}

type regexResult struct {
	pid string
	name string
	label string
	audit_id string
	success string
	operation string
}

func TailAuditd(DIR string) {
	state := State{}

	t, err := tail.TailFile("/var/log/audit/audit.log", tail.Config{ 
		Follow: true,
		Location: &tail.SeekInfo{ Offset: 0, Whence: io.SeekEnd }}) // we only wanna know what happens after we start running the daemon

	utils.CheckErr(err)

	// println("-- tailing --\n")
	
	go func() { // run this part concurrently
		for line := range t.Lines { // auditd has 3 parts, syscall, path and avc

			// log.Println(line)

			if strings.Contains(line.Text, "cwalld") && strings.Contains(line.Text, "SYSCALL") { // this is the syscall part, containing pid, operation and subject name
				state.trackSubject(line.Text)
			}

			if strings.Contains(line.Text, DIR) { // this is the path line, containing the affected object path
				state.trackObject(line.Text)
			}

			if strings.Contains(line.Text, "denied") { // this inclues avc denials
				state.trackAVC(line.Text)
			}
		}
	}()
	
	<-make(chan struct{})
}

func (state *State) trackSubject(line string) {
	// log.Println(line)
	regexes := regexer(line)
	if regexes.name == "setroubleshootd" { return }

	if regexes.success == "no" {
		log.Println("AVC Denial Succesful")
	}

	flags, err := strconv.ParseInt(regexes.operation, 16, 64) // convert the string that is hexadecimal, into straight binary, which is read as an int64 but actually is just straight flags 
	utils.CheckErr(err)

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

	var sjct *subject.Subject
	
	for _, s := range state.subjects { // if subject is already accounted for
		if s.Pid == regexes.pid {
			sjct = &s
			break
		}
	}

	if sjct == nil { // add it to the global list of subjects if not
		sjct = &subject.Subject{ Pid: regexes.pid, Name: regexes.name, Label: regexes.label }
		sjct.ToString()
		state.subjects = append(state.subjects, *sjct)
	}

	audit := audit.Audit{ Id: regexes.audit_id, Subject: sjct, Object: nil, Operation: op } // create new audit - only half complete so far
	state.audits = append(state.audits, audit)
}

func (state *State) trackObject(line string) {
	regex := regexp.MustCompile(`\bmsg=audit\((([^)]+))`)
	audit_id := regex.FindStringSubmatch(line)[1]

	for i := range state.audits {
		if state.audits[i].Id == audit_id {
			if state.audits[i].Subject.Name == "cwalld" { return } // dont track cwalld 
		}
	}

	regex = regexp.MustCompile(`\bname="([^"]+)"`)
	object_path := regex.FindStringSubmatch(line)[1]

	object_label, err := selinux.FileLabel(object_path)
	utils.CheckErr(err)

	for i := range state.audits {
		if state.audits[i].Id == audit_id {
			state.audits[i].Object = &object.Object{ Name: object_path, Label: object_label } 
			state.audits[i].ToString()
			// state.audits[i].Subject.

			break
		}
	}
}

func (state *State) trackAVC(line string) {
	regex := regexp.MustCompile(`\{ ([^ }]+)`)
	operation := regex.FindStringSubmatch(line)[1]

	regex = regexp.MustCompile(`\bpid=(\d+)`)
	pid := regex.FindStringSubmatch(line)[1]

	regex = regexp.MustCompile(`\bname="([^"]+)"`)
	object := regex.FindStringSubmatch(line)[1]
	
	for _, s := range state.subjects {
		if s.Pid == pid {
			logDenial(s.Name, operation, object)
			break
		}
	}
}

func logDenial(s string, op string, obj string) { // operation here is just text because its reprented in string form by AVC already
	fmt.Printf("<!DENIAL!>:\t%s\tattempted { %s }\ton %s\n\n", s, op, obj)
}

func regexer(line string) regexResult {
	s := regexResult{}

	regex := regexp.MustCompile(`\bpid=(\d+)`) // regex to catch pid
	regex_pid := regex.FindStringSubmatch(line) // pid[0] = "pid=..." pid[1] = "..."
	pid := utils.RegexErr(regex_pid, "pid")

	s.pid = pid

	regex = regexp.MustCompile(`\bcomm="([^"]+)"`) // regex to catch subject name
	regex_subject_name := regex.FindStringSubmatch(line)
	subject_name := utils.RegexErr(regex_subject_name, "subject name")

	s.name = subject_name

	intpid, err := strconv.Atoi(pid) // convert to int for PidLabel function
	utils.CheckErr(err)

	label, err := selinux.PidLabel(intpid)
	utils.CheckErr(err)

	regex = regexp.MustCompile(`r:([^:]+)`)
	regex_label_type := regex.FindStringSubmatch(label)
	label_type := utils.RegexErr(regex_label_type, "label type")

	s.label = label_type

	regex = regexp.MustCompile(`\bmsg=audit\(([^)]+)`) // regex to catch audit id to combine with other line
	regex_audit_id := regex.FindStringSubmatch(line)
	audit_id := utils.RegexErr(regex_audit_id, "audit id")

	s.audit_id = audit_id

	regex = regexp.MustCompile(`\bsuccess=([^ ]+)`)
	regex_success := regex.FindStringSubmatch(line)
	success := utils.RegexErr(regex_success, "success")

	s.success = success

	regex = regexp.MustCompile(`\ba2=([^ ]+)`)
	regex_operation := regex.FindStringSubmatch(line)
	operation := utils.RegexErr(regex_operation, "operation")

	s.operation = operation

	return s
}

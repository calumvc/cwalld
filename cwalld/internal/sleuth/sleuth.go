package sleuth

import (
	"cwalld/internal/audit"
	"cwalld/internal/object"
	"cwalld/internal/subject"
	"cwalld/internal/utils"
	"fmt"
	"io"
	"os"
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
	println("Chinese Wall Enforcing")
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

			if strings.Contains(line.Text, DIR) && strings.Contains(line.Text, "PATH") { // this is the path line, containing the affected object path
				state.trackObject(line.Text)
			}

			if strings.Contains(line.Text, "AVC") { //&& strings.Contains(line.Text, "path"){ // this inclues avc denials
				state.trackAVC(line.Text)
			}
		}
	}()
	
	<-make(chan struct{})
}

func (state *State) trackSubject(line string) { // we will track details about the subject from this audit, creating details for a new subject if we havent seen it before
	regexes := regexer(line)
	if regexes.name == "setroubleshootd" { return } // this guy is annoying

	var subj *subject.Subject
	
	for _, s := range state.subjects { // if subject is already registered
		if s.Pid == regexes.pid {
			subj = &s
			break
		} else
		if s.Name == regexes.name {
			subj = &s
			fmt.Printf("%s changed label to %s\n", s.Name, s.Label)
		}
	}

	entrypoint, err := os.Readlink(fmt.Sprintf("/proc/%s/exe", regexes.pid))
	utils.CheckErr(err)

	if subj == nil { // add it to the global list of subjects if not
		subj = &subject.Subject{ Pid: regexes.pid, Name: regexes.name, Label: regexes.label, Entrypoint: entrypoint }
		state.subjects = append(state.subjects, *subj)
		subj.ToString()
	}

	success := true
	if regexes.success == "no" {
		success = false
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

	audit := audit.Audit{ Id: regexes.audit_id, Subject: subj, Object: nil, Operation: op, Success: success } // create new audit - only half complete so far
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
	regex_object_path := regex.FindStringSubmatch(line)
	object_path := utils.RegexErr(regex_object_path, "object name")

	object_label, err := selinux.FileLabel(object_path)
	utils.CheckErr(err)

	regex = regexp.MustCompile(`r:([^:]+)`)
	regex_label_type := regex.FindStringSubmatch(object_label)
	label_type := utils.RegexErr(regex_label_type, "label type")

	for i := range state.audits {
		if state.audits[i].Id == audit_id {
			state.audits[i].Object = &object.Object{ Name: object_path, Label: label_type } 
			state.audits[i].ToString()

			if state.audits[i].Success == true { // if it succesfully read/wrote, then alter the label as necessary
				state.audits[i].Subject.AlterLabel(label_type, state.audits[i].Operation)
			}

			break
		}
	}
}

func (state *State) trackAVC(line string) {
	regex := regexp.MustCompile(`\{ ([^ }]+)`)
	regex_operation := regex.FindStringSubmatch(line)
	operation := utils.RegexErr(regex_operation, "Operation")
	
	var object string

	if strings.Contains(line, "name") {
		regex = regexp.MustCompile(`\bname="([^"]+)"`)
		regex_object := regex.FindStringSubmatch(line)
		object = utils.RegexErr(regex_object, "Object")
	} else
	if strings.Contains(line, "path") {
		regex = regexp.MustCompile(`\bpath="([^"]+)"`)
		regex_object := regex.FindStringSubmatch(line)
		object = utils.RegexErr(regex_object, "Object")
	}

	regex = regexp.MustCompile(`\bpid=(\d+)`)
	regex_pid := regex.FindStringSubmatch(line)
	pid := utils.RegexErr(regex_pid, "Pid")

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

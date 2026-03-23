package sleuth

import (
	"cwalld/internal/audit"
	"cwalld/internal/decorator"
	"cwalld/internal/subject"
	"cwalld/internal/utils"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/nxadm/tail"
	"github.com/opencontainers/selinux/go-selinux"
	"golang.org/x/sys/unix"
)

type State struct {
	subjects []subject.Subject
	audits []audit.Audit
}

type regexResult struct { // just used to more easily seperate regex logic from main logic
	pid string
	name string
	label string
	audit_id string
	success string
	operation string
}

func TailAuditd(DIR string) error {
	state := State{}

	t, err := tail.TailFile("/var/log/audit/audit.log", tail.Config{ 
		Follow: true, // keep reading new lines
		ReOpen: true, // follow & reopen new log rotations
		Location: &tail.SeekInfo{ Offset: 0, Whence: io.SeekEnd }}) // we only wanna know what happens after we start running the daemon

	if err != nil {
		return err
	}

	for line := range t.Lines { // auditd has 3 parts, syscall, path and avc
		if strings.Contains(line.Text, "setroubleshootd") { continue } // ignore this guy

		if strings.Contains(line.Text, "cwalld") && strings.Contains(line.Text, "SYSCALL") { // this is the syscall part, containing pid, operation and subject name
			err := state.trackSubject(line.Text)
			if err != nil {
				return err
			}
		}
		if strings.Contains(line.Text, DIR) && strings.Contains(line.Text, "PATH") { // this is the path line, containing the affected object path
			err := state.trackObject(line.Text)
			if err != nil {
				return err
			}
		}

		if strings.Contains(line.Text, "AVC") { // this inclues avc denials
			err := state.trackAVC(line.Text)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (state *State) trackSubject(line string) error { // we will track details about the subject from this audit, creating details for a new subject if we havent seen it before 
	regexes, err := regexer(line)
	if err != nil {
		if err.Error() != "Atomic process" {
			return err
		}

		return nil // we just wanna ignore atomic processes entirely
	}

	if regexes.name == "cwalld-enforce" { return nil } // if we log ourselves we will start an infinite loop

	var subj *subject.Subject
	seen := false // this is so we can see when a process comes back with a new pid
	
	for i, s := range state.subjects { // if subject is already registered
		if s.Pid == regexes.pid {
			subj = &s
			break
		} else
		if s.Name == regexes.name { // if weve seen the process before but it got restarted - likely because of a label change
			state.subjects[i].Pid = regexes.pid // update pid so we match it correctly when it comes back
			state.subjects[i].Label = regexes.label 
			seen = true
			break
		}
	}

	entrypoint, err := os.Readlink(fmt.Sprintf("/proc/%s/exe", regexes.pid)) // get the entrypoint of the subject

	if err != nil {
		return err
	}

	if subj == nil { // add it to the global list of subjects if not
		subj = &subject.Subject{ Pid: regexes.pid, Name: regexes.name, Label: regexes.label, Entrypoint: entrypoint }
		state.subjects = append(state.subjects, *subj)
		if seen != true {
			decorator.DecorateAndLog(subj.String(), decorator.Register)
		} else {
			decorator.DecorateAndLog(subj.ReString(), decorator.Reregister)
			seen = false
		}
	}

	success := true
	if regexes.success == "no" { // if the process fails its still audited so we can notify from that
		success = false
	}

	flags, err := strconv.ParseInt(regexes.operation, 16, 64) // convert the string that is hexadecimal into binary, which is read as an int64 but actually is just straight flags of syscalls

	if err != nil { 
		return err 
	}

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

	audit := audit.Audit{ Id: regexes.audit_id, Subject: subj, Object: nil, Operation: op, Success: success } // create new audit - will finish it in trackObject
	state.audits = append(state.audits, audit)

	return nil
}

func (state *State) trackObject(line string) error {
	regex := regexp.MustCompile(`\bmsg=audit\((([^)]+))`)
	audit_id := regex.FindStringSubmatch(line)[1]

	for i := range state.audits {
		if state.audits[i].Id == audit_id {
			if state.audits[i].Subject.Name == "cwalld-enforce" { return nil } // dont track cwalld 
		}
	}

	regex = regexp.MustCompile(`\bname="([^"]+)"`)
	regex_object_path := regex.FindStringSubmatch(line)
	object_path, err := utils.RegexErr(regex_object_path, "object name")
	
	if err != nil {
		return err
	}

	object_label, err := selinux.FileLabel(object_path)

	if err != nil { 
		return err
	}

	regex = regexp.MustCompile(`r:([^:]+)`)
	regex_label_type := regex.FindStringSubmatch(object_label)
	label_type, err := utils.RegexErr(regex_label_type, "label type")

	if err != nil {
		return err
	}

	for i := range state.audits {
		if state.audits[i].Id == audit_id {
			state.audits[i].Object = &utils.Object{ Name: object_path, Label: label_type } 
			decorator.DecorateAndLog(state.audits[i].String(), decorator.Audit)

			if state.audits[i].Success == true { // if it succesfully read/wrote, then alter the label as necessary
				state.audits[i].Subject.AlterLabel(label_type, state.audits[i].Operation)
			}

			break
		}
	}

	return nil
}

func (state *State) trackAVC(line string) error {
	regex := regexp.MustCompile(`\{ ([^ }]+)`)
	regex_operation := regex.FindStringSubmatch(line)
	operation, err := utils.RegexErr(regex_operation, "Operation")

	if err != nil {
		return err
	}
	
	var object string

	if strings.Contains(line, "name") {
		regex = regexp.MustCompile(`\bname="([^"]+)"`)
		regex_object := regex.FindStringSubmatch(line)
		object, err = utils.RegexErr(regex_object, "Object")

		if err != nil { 
			return err
		}
	} else
	if strings.Contains(line, "path") {
		regex = regexp.MustCompile(`\bpath="([^"]+)"`)
		regex_object := regex.FindStringSubmatch(line)
		object, err = utils.RegexErr(regex_object, "Object")

		if err != nil {
			return err
		}
	}

	regex = regexp.MustCompile(`\bpid=(\d+)`)
	regex_pid := regex.FindStringSubmatch(line)
	pid, err := utils.RegexErr(regex_pid, "Pid")

	if err != nil { 
		return err 
	}

	for _, s := range state.subjects {
		if s.Pid == pid {
			line := fmt.Sprintf("%s\tattempted { %s }\ton %s", s.Name, operation, object)
			decorator.DecorateAndLog(line, decorator.Denial)
			break
		}
	}

	return nil
}

func regexer(line string) (*regexResult, error) {
	s := regexResult{}

	regex := regexp.MustCompile(`\bpid=(\d+)`) // regex to catch pid
	regex_pid := regex.FindStringSubmatch(line) // pid[0] = "pid=..." pid[1] = "..."
	pid, err := utils.RegexErr(regex_pid, "pid")

	if err != nil {
		return nil, err
	}

	s.pid = pid

	regex = regexp.MustCompile(`\bcomm="([^"]+)"`) // regex to catch subject name
	regex_subject_name := regex.FindStringSubmatch(line)
	subject_name, err := utils.RegexErr(regex_subject_name, "subject name")

	if err != nil {
		return nil, err
	}

	s.name = subject_name

	intpid, err := strconv.Atoi(pid) // convert to int for PidLabel function

	if err != nil {
		return nil, err
	}

	label, err := selinux.PidLabel(intpid)

	if err != nil {
		return nil, err
	}

	if label == "" { // must be an atomic process (like cat), so we just ignore it
		decorator.DecorateAndLog(subject_name, decorator.Atomic)
		return nil, fmt.Errorf("Atomic process")
	}

	regex = regexp.MustCompile(`r:([^:]+)`)
	regex_label_type := regex.FindStringSubmatch(label)
	label_type, err := utils.RegexErr(regex_label_type, "label type")

	if err != nil {
		return nil, err
	}

	s.label = label_type

	regex = regexp.MustCompile(`\bmsg=audit\(([^)]+)`) // regex to catch audit id to combine with other line
	regex_audit_id := regex.FindStringSubmatch(line)
	audit_id, err := utils.RegexErr(regex_audit_id, "audit id")

	if err != nil {
		return nil, err
	}

	s.audit_id = audit_id

	regex = regexp.MustCompile(`\bsuccess=([^ ]+)`)
	regex_success := regex.FindStringSubmatch(line)
	success, err:= utils.RegexErr(regex_success, "success")

	if err != nil {
		return nil, err
	}

	s.success = success

	regex = regexp.MustCompile(`\ba2=([^ ]+)`)
	regex_operation := regex.FindStringSubmatch(line)
	operation, err := utils.RegexErr(regex_operation, "operation")

	s.operation = operation

	if err != nil {
		return nil, err
	}

	return &s, nil
}

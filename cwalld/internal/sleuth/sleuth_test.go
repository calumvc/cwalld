package sleuth

import (
	"cwalld/internal/audit"
	"cwalld/internal/subject"
	"cwalld/internal/utils"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/opencontainers/selinux/go-selinux"
)

func TestSetup(t *testing.T) {
	cmd := exec.Command("sudo", "chcon", "-t", "bin_t", "/usr/local/bin/cwalldtestd")
	err := cmd.Run()
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	cmd = exec.Command("sudo", "systemctl", "restart", "cwalldtestd.service")
	err = cmd.Run()
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
}

func TestTrackSubject(t *testing.T) { // this should be the same for any time its ran
	state := State{}

	cmd := exec.Command("pgrep", "cwalldtestd") // need to get the pid of the live process for the code to work
	response, err := cmd.CombinedOutput()

	pid := strings.TrimSpace(string(response))

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	expected_subj := subject.Subject{
		Pid: pid,
		Name: "cwalldtestd",
		Label: "unconfined_service_t",
		Entrypoint: "/usr/local/bin/cwalldtestd",
	}

	line := fmt.Sprintf("type=SYSCALL msg=audit(1773872750.325:248): arch=c000003e syscall=257 success=yes exit=3 a0=ffffff9c a1=402018 a2=0 a3=0 items=1 ppid=1 pid=%s auid=4294967295 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=(none) ses=4294967295 comm=\"cwalldtestd\" exe=\"/usr/local/bin/cwalldtestd\" subj=system_u:system_r:unconfined_service_t:s0 key=\"cwalld\"", pid)
	err = state.trackSubject(line)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	if !reflect.DeepEqual(state.subjects[0], expected_subj) {
		t.Errorf("\nexpected\n%s\ngot\n%s", expected_subj.String(), state.subjects[0].String())
	}

	expected_audit := audit.Audit{
		Id: "1773872750.325:248",
		Subject: &state.subjects[0],
		Object: nil,
		Operation: utils.Read,
		Success: true,
	}

	if !reflect.DeepEqual(state.audits[0], expected_audit) {
		t.Errorf("expected %s got %s", expected_audit.String(), state.audits[0].String())
	}
}

func TestTrackObjectLabelChange(t *testing.T) { // a test where it should change the label
	cmd := exec.Command("pgrep", "cwalldtestd")
	response, err := cmd.CombinedOutput()

	pid := strings.TrimSpace(string(response))

	intpid, err := strconv.Atoi(pid)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	old_label, err := selinux.PidLabel(intpid)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	state := State{
		subjects: []subject.Subject{
			{
				Pid: pid,
				Name: "cwalldtestd",
				Label: "unconfined_service_t",
				Entrypoint: "/usr/local/bin/cwalldtestd",
			},
		},
	}

	state.audits = append(state.audits, audit.Audit{
		Id: "1773872750.325:248",
		Subject: &state.subjects[0],
		Object: nil,
		Operation: utils.Read,
		Success: true,
	})

	expected_audit := audit.Audit{
		Id: "1773872750.325:248",
		Subject: &state.subjects[0],
		Object: &utils.Object{ Name: "/home/testgrounds/objects/alpha_logs", Label: "alpha_t" },
		Operation: utils.Read,
		Success: true,
	}

	line := "type=PATH msg=audit(1773872750.325:248): item=0 name=\"/home/testgrounds/objects/alpha_logs\" inode=58142854 dev=fd:00 mode=0100666 ouid=1000 ogid=1000 rdev=00:00 obj=unconfined_u:object_r:alpha_t:s0 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0"
	err = state.trackObject(line)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	if !reflect.DeepEqual(state.audits[0], expected_audit) {
		t.Errorf("expected %s got %s", expected_audit.String(), state.audits[0].String())
	}

	cmd = exec.Command("pgrep", "cwalldtestd")
	response, err = cmd.CombinedOutput()

	pid = strings.TrimSpace(string(response))

	intpid, err = strconv.Atoi(pid)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	new_label, err := selinux.PidLabel(intpid)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	expected_old := "system_u:system_r:unconfined_service_t:s0"
	expected_new := "system_u:system_r:alpha_rw_t:s0"
	if old_label != expected_old  && new_label != expected_new {
		t.Errorf("expected old %s new %s got old %s new %s", expected_old, expected_new, old_label, new_label)
	}
}

func TestTrackObjectNoLabelChange(t *testing.T) { // a test where it shouldnt change the label
	cmd := exec.Command("pgrep", "cwalldtestd")
	response, err := cmd.CombinedOutput()

	pid := strings.TrimSpace(string(response))

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	state := State{
		subjects: []subject.Subject{
			{
				Pid: pid,
				Name: "cwalldtestd",
				Label: "alpha_rw_t",
				Entrypoint: "/usr/local/bin/cwalldtestd",
			},
		},
	}

	state.audits = append(state.audits, audit.Audit{ // this part should already exist in it because of trackSubject is ran before it
		Id: "1773872750.325:248",
		Subject: &state.subjects[0],
		Object: nil,
		Operation: utils.Read,
		Success: true,
	})

	expected_audit := audit.Audit{
		Id: "1773872750.325:248",
		Subject: &state.subjects[0],
		Object: &utils.Object{ Name: "/home/testgrounds/objects/alpha_logs", Label: "alpha_t" },
		Operation: utils.Read,
		Success: true,
	}

	line := "type=PATH msg=audit(1773872750.325:248): item=0 name=\"/home/testgrounds/objects/alpha_logs\" inode=58142854 dev=fd:00 mode=0100666 ouid=1000 ogid=1000 rdev=00:00 obj=unconfined_u:object_r:alpha_t:s0 nametype=NORMAL cap_fp=0 cap_fi=0 cap_fe=0 cap_fver=0 cap_frootid=0"
	err = state.trackObject(line)

	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	if !reflect.DeepEqual(state.audits[0], expected_audit) {
		t.Errorf("expected %s got %s", expected_audit.String(), state.audits[0].String())
	}
}

func TestTrackAVC(t *testing.T) {

}


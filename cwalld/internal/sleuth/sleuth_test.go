package sleuth

import (
	"cwalld/internal/subject"
	"testing"
)

func TestTrackSubject(t *testing.T) {
	state := State{}

	expected := subject.Subject{
		Pid: "32220",
		Name: "alpha_civild",
		Label: "alpha_rw_t",
		Entrypoint: "/usr/local/bin/alpha_civild",
	}

	line := "type=SYSCALL msg=audit(1773707745.541:3401283): arch=c000003e syscall=257 success=yes exit=3 a0=ffffff9c a1=402018 a2=241 a3=1b6 items=2 ppid=1 pid=32220 auid=4294967295 uid=0 gid=0 euid=0 suid=0 fsuid=0 egid=0 sgid=0 fsgid=0 tty=(none) ses=4294967295 comm=\"alpha_civild\" exe=\"/usr/local/bin/alpha_civild\" subj=system_u:system_r:alpha_rw_t:s0 key=\"cwalld\" type=PROCTITLE msg=audit(1773707749.542:3401288): proctitle=\"/usr/local/bin/alpha_civild\""
	err := state.trackSubject(line)

	if err != nil {
		t.Errorf("Error")
	}

	if state.subjects[0] != expected {
		t.Errorf("expected %s got %s", expected.String(), state.subjects[0].String())
	}
}

func TestTrackObject(t *testing.T) {

}

func TestTrackAVC(t *testing.T) {

}

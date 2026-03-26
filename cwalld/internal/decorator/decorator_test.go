package decorator

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestDecorateAndLog(t *testing.T) {
	cases := []struct{
		in string
		in2 Decor
		want string
	}{
		{
			in: "audit details...",
			in2: Audit,
			want: "|Audit:\taudit details...",
		},
		{
			in: "subj details...",
			in2: Register,
			want: "|New Subject:\tsubj details...",
		},
		{
			in: "subj details...",
			in2: Reregister,
			want: "|Reregistered:\tsubj details...",
		},
		{
			in: "denial details...",
			in2: Denial,
			want: "<!DENIAL!>:\tdenial details...\t<!DENIAL!>",
		},
		{
			in: "relabel details...",
			in2: Relabel,
			want: "|Relabel:\trelabel details...",
		},
		{
			in: "dbus details...",
			in2: Dbus,
			want: "|Daemon restart successful:\tdbus details...",
		},
		{
			in: "atomic details...",
			in2: Atomic,
			want: "|Atomic process occured:\tatomic details...",
		},
		{
			in: "error details...",
			in2: Error,
			want: "ERROR:\terror details...",
		},
	}

	for _, c := range cases {
		cmd := exec.Command("sudo", "sh", "-c", "echo '' > /var/log/cwall/cwall.log")
		err := cmd.Run()
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		DecorateAndLog(c.in, c.in2)
		cmd = exec.Command("cat", "/var/log/cwall/cwall.log")
		res, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		if !bytes.Contains(res, []byte(c.want)) {
			t.Errorf("Got:\n%s wanted:\n%s", string(res), c.want)
		}

	}

}

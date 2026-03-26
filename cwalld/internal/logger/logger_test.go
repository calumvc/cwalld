package logger

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestLogger(t *testing.T) {
	cases := []struct{
		in string
		want string
	}{
		{
			in: "test",
			want: "test",
		},
		{
			in: "test2",
			want: "test2",
		},
	}

	for _, c := range cases {
		cmd := exec.Command("sudo", "sh", "-c", "echo '' > /var/log/cwall/cwall.log")
		err := cmd.Run()
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		Log(c.in)
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

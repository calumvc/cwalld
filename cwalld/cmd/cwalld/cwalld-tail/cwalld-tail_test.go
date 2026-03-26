package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	oldStoud := os.Stdout
	defer func() { os.Stdout = oldStoud }()

	cases := []struct{
		in string
		want string
	}{
		{
			in: "test",
			want: "\ntest\n",
		},
	}

	for _, c := range cases {

		r, w, _ := os.Pipe()
		os.Stdout = w

		cmd := exec.Command("sudo", "sh", "-c", "echo '' > /var/log/cwall/cwall.log")
		err := cmd.Run()
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		go cwalldtail()

		cmd = exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo %s >> /var/log/cwall/cwall.log", c.in))
		err = cmd.Run()
		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		w.Close()

		var buf bytes.Buffer
		io.Copy(&buf, r)

		out := buf.String()
		if out != c.want {
			t.Errorf("Got %s wanted %s", buf.String(), c.want)
		}
	}
}

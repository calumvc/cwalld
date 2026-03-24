package subject

import (
	"cwalld/internal/utils"
	"os/exec"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	cases := []struct{
		in Subject
		want string
	}{
		{
			in: Subject{
				Pid: "12312",
				Name: "test",
				Label: "test_t",
				Entrypoint: "subj/test",
			},
			want : "pid=12312\tcomm=test\tlabel=test_t\tentrypoint=subj/test",
		},
	}

	for _, c := range cases {
		got := c.in.String()

		if got != c.want {
			t.Errorf("got %s wanted %s", got, c.want)
		}
	}
}

func TestReString(t *testing.T) {
	cases := []struct{
		in Subject
		want string
	}{
		{
			in: Subject{
				Pid: "12312",
				Name: "test",
				Label: "test_t",
				Entrypoint: "subj/test",
			},
			want : "test under label test_t and new pid 12312",
		},
	}

	for _, c := range cases {
		got := c.in.ReString()

		if got != c.want {
			t.Errorf("got %s wanted %s", got, c.want)
		}
	}
}

func TestAlterLabelLayer2(t *testing.T) {
	cases := []struct{
		in Subject
		in2 string
		in3 string
		want bool
	}{
		{ 
			in: Subject{
					Pid: "1",
					Name: "cwalldtestd",
					Label: "unconfined_service_t",
					Entrypoint: "/usr/local/bin/cwalldtestd",
			},
			in2: "alpha_t",
			in3: "alpha_rw_t",
			want: true,
		},
		{ 
			in: Subject{
					Pid: "1",
					Name: "cwalldtestd",
					Label: "unconfined_service_t",
					Entrypoint: "/usr/local/bin/cwalldtestd",
			},
			in2: "beta_t",
			in3: "beta_rw_t",
			want: true,
		},
		{ 
			in: Subject{
					Pid: "1",
					Name: "cwalldtestd",
					Label: "unconfined_service_t",
					Entrypoint: "/usr/local/bin/cwalldtestd",
			},
			in2: "gamma_t",
			in3: "gamma_rw_t",
			want: true,
		},
	}

	for _, c := range cases {

		cmd := exec.Command("sudo", "chcon", "-t", "bin_t", "/usr/local/bin/cwalldtestd")
		res, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		cmd = exec.Command("sudo", "systemctl", "start", "cwalldtestd")
		res, err = cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		cmd = exec.Command("sudo", "ps", "-efZ")
		res, err = cmd.CombinedOutput()

		correct := false
		for line := range strings.SplitSeq(string(res), "\n") {
			if strings.Contains(line, "cwalldtestd") {
				if strings.Contains(line, "unconfined_service_t") {
					t.Log(line)
					correct = true
				}
			}
		}

		if !correct {
			t.Error("Error setting test up")
		}

		c.in.AlterLabel(c.in2, utils.Read)

		cmd = exec.Command("sudo", "ps", "-efZ")
		res, err = cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		correct = false
		for line := range strings.SplitSeq(string(res), "\n") {
			if strings.Contains(line, "cwalldtestd"){
				if strings.Contains(line, c.in3) {
					correct = true
				}
			}
		} 

		if correct == false {
			t.Errorf("Label change failed")
		}

		cmd = exec.Command("sudo", "systemctl", "stop", "cwalldtestd")
		res, err = cmd.CombinedOutput()

		if correct != c.want {
			t.Errorf("Error: %s", err.Error())
		}
	}
}

func TestAlterLabelLayer3(t *testing.T) {
	cases := []struct{
		in Subject
		in2 string // object label read
		in3 string // label to end on
		in4 string // entrypoint label at start
		in5 string // label to run as after start
		want bool
	}{
		{ 
			in: Subject{
					Pid: "1",
					Name: "cwalldtestd",
					Label: "alpha_rw_t",
					Entrypoint: "/usr/local/bin/cwalldtestd",
			},
			in2: "gamma_t",
			in3: "alpha_gamma_r_t",
			in4: "alpha_rw_exec_t",
			in5: "alpha_rw_t",
			want: true,
		},
		{ 
			in: Subject{
					Pid: "1",
					Name: "cwalldtestd",
					Label: "beta_rw_t",
					Entrypoint: "/usr/local/bin/cwalldtestd",
			},
			in2: "gamma_t",
			in3: "beta_gamma_r_t",
			in4: "beta_rw_exec_t",
			in5: "beta_rw_t",
			want: true,
		},
		{ 
			in: Subject{
					Pid: "1",
					Name: "cwalldtestd",
					Label: "gamma_rw_t",
					Entrypoint: "/usr/local/bin/cwalldtestd",
			},
			in2: "alpha_t",
			in3: "alpha_gamma_r_t",
			in4: "gamma_rw_exec_t",
			in5: "gamma_rw_t",
			want: true,
		},
	}

	for _, c := range cases {

		cmd := exec.Command("sudo", "chcon", "-t", c.in4, "/usr/local/bin/cwalldtestd")
		res, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		cmd = exec.Command("sudo", "systemctl", "start", "cwalldtestd")
		res, err = cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		cmd = exec.Command("sudo", "ps", "-efZ")
		res, err = cmd.CombinedOutput()

		correct := false
		for line := range strings.SplitSeq(string(res), "\n") {
			if strings.Contains(line, "cwalldtestd") {
				if strings.Contains(line, c.in5) {
					t.Log(line)
					correct = true
				}
			}
		}

		if !correct {
			t.Error("Error setting test up")
		}

		c.in.AlterLabel(c.in2, utils.Read)

		cmd = exec.Command("sudo", "ps", "-efZ")
		res, err = cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Error: %s", err.Error())
		}

		correct = false
		for line := range strings.SplitSeq(string(res), "\n") {
			if strings.Contains(line, "cwalldtestd"){
				t.Log(line)
				if strings.Contains(line, c.in3) {
					correct = true
				}
			}
		} 

		if correct == false {
			t.Errorf("Label change failed")
		}

		cmd = exec.Command("sudo", "systemctl", "stop", "cwalldtestd")
		res, err = cmd.CombinedOutput()

		if correct != c.want {
			t.Errorf("Error: %s", err.Error())
		}
	}
}

func TestRestartSubject(t *testing.T) {}

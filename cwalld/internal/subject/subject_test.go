package subject

import (
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
			want : "test under label test_t",
		},
	}

	for _, c := range cases {
		got := c.in.ReString()

		if got != c.want {
			t.Errorf("got %s wanted %s", got, c.want)
		}
	}
}

func TestAlterLabel(t *testing.T) {}

func TestRestartSubject(t *testing.T) {}

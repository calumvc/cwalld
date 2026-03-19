package audit

import (
	"cwalld/internal/subject"
	"cwalld/internal/utils"
	"testing"
)

func TestString(t *testing.T) {
	cases := []struct{
		in Audit
		want string
	}{
		{
			in: Audit{ 
				Id: "1773872750.325:248", 
				Subject: &subject.Subject{
					Pid: "12312",
					Name: "test",
					Label: "test_t",
					Entrypoint: "subj/test",
				}, 
				Object: &utils.Object{
					Name: "obj/test",
					Label: "obj_t",
				},
				Operation: utils.Read,
				Success: true,
			},
			want: "subject=test : test_t\toperation=Read : true\tobject=obj/test : obj_t",
		},
	}

	for _, c := range cases {
		got := c.in.String()

		if got != c.want {
			t.Errorf("got %s wanted %s", got, c.want)
		}
	}
}

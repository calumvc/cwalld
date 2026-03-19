package utils

import (
	"errors"
	"testing"
)

func TestString(t *testing.T) {
	cases := []struct{
		in Operation
		want string
	}{
		{ Unknown, "Unknown" },
		{ Read, "Read" },
		{ Write, "Write" },
		{ ReadWrite, "ReadWrite" },
		{ Metadata, "Metadata" },
	}

	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("got %s expected %s", got, c.want)
		}
	}
}

func TestRegexErr(t *testing.T) {
	type input struct {
		s []string
		msg string
	}

	type expected struct {
			res string
			err error
	}

	cases := []struct{
		in input
		want expected
	}{
		{
			in: input{ s: []string{ "test=works", "works" }, msg: "doesit?" }, 
			want: expected{ res: "works", err: nil },
		}, 
		{ 
			in: input{ s: nil, msg: "doesntwork" }, 
			want: expected{ res: "", err: errors.New("Regex failed on doesntwork") },
		},
	}

	for _, c := range cases {
		res, err := RegexErr(c.in.s, c.in.msg)

		if res != c.want.res {
			t.Errorf("got %s expected %s", res, c.want.res)
		}

		if (err != nil && c.want.err == nil) || (c.want.err != nil && err == nil) {
			t.Errorf("got %s expected %s", err, c.want.err)
		}
	}
}

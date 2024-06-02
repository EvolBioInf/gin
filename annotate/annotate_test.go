package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestAnnotate(t *testing.T) {
	var tests []*exec.Cmd
	g := "../data/toy.gff"
	i := "../data/toyIv.txt"
	test := exec.Command("./annotate", g, i)
	tests = append(tests, test)
	test = exec.Command("./annotate", "-s", g, i)
	tests = append(tests, test)
	test = exec.Command("./annotate", "-t", g, i)
	tests = append(tests, test)
	test = exec.Command("./annotate", "-c", g, i)
	tests = append(tests, test)
	test = exec.Command("./annotate", "-l", "0", g, i)
	tests = append(tests, test)
	test = exec.Command("./annotate", "-l", "1", g, i)
	tests = append(tests, test)
	test = exec.Command("./annotate", "-l", "2", g, i)
	tests = append(tests, test)
	test = exec.Command("./annotate", "-l", "3", g, i)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err.Error())
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("can't read %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s", get, want)
		}
	}
}

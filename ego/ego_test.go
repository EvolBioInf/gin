package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strconv"
	"testing"
)

func TestEgo(t *testing.T) {
	var tests []*exec.Cmd
	ge := "../data/gene2go"
	ob := "../data/obsId.txt"
	ra := "../data/expId.txt"
	test := exec.Command("./ego", "-o", ge, ob, ra)
	tests = append(tests, test)
	test = exec.Command("./ego", "-o", ge, "-m", "5", ob, ra)
	tests = append(tests, test)
	test = exec.Command("./ego", "-o", ge, "-r", ob, ra)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Errorf("can't run %q", test)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := ioutil.ReadFile(f)
		if err != nil {
			t.Errorf("can't open %q", f)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}

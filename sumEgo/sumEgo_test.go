package main

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestSumEgo(t *testing.T) {
	tests := make([]*exec.Cmd, 0)
	i1 := "../data/ego1.txt"
	i2 := "../data/ego2.txt"
	i3 := "../data/ego3.txt"
	test := exec.Command("./sumEgo", i1)
	tests = append(tests, test)
	test = exec.Command("./sumEgo", i1, i2)
	tests = append(tests, test)
	test = exec.Command("./sumEgo", i1, i2, i3)
	tests = append(tests, test)
	for i, test := range tests {
		get, err := test.Output()
		if err != nil {
			t.Error(err)
		}
		f := "r" + strconv.Itoa(i+1) + ".txt"
		want, err := os.ReadFile(f)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(get, want) {
			t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
		}
	}
}

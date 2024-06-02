package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestShuffle(t *testing.T) {
	g := "../data/toy.gff"
	i := "../data/toyIv.txt"
	test := exec.Command("./shuffle", "-s", "3",
		"-n", "20", g, i)
	get, err := test.Output()
	if err != nil {
		t.Errorf("can't run %q", test)
	}
	f := "r.txt"
	want, err := ioutil.ReadFile(f)
	if err != nil {
		t.Errorf("can't open %q", f)
	}
	if !bytes.Equal(get, want) {
		t.Errorf("get:\n%s\nwant:\n%s", get, want)
	}
}

package main

import (
	"os/exec"
	"testing"
)

func TestStartCommand(t *testing.T) {
	outputByt, err := exec.Command("sh", "./out/test1").Output()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(outputByt))
}

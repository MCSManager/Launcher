package main

import (
	"os/exec"
	"testing"
)

func TestStartCommand(t *testing.T) {
	cmder := exec.Command("/Users/xiamingjie/gopath/src/test-go/test")
	reader, err := cmder.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	err = cmder.Start()
	if err != nil {
		t.Fatal(err)
		return
	}

	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil {
			t.Log(err)
			break
		}
		t.Log(string(buf[:n]))
	}

	t.Fatal(cmder.Wait())
}

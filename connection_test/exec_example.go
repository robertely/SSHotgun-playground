package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var (
		cmdOut []byte
		err    error
	)

	cmdName := "scp"
	cmdArgs := []string{"-B",
		"-p",
		"-i ../test_fixture/testing_key.rsa",
		"-oStrictHostKeyChecking=no",
		"-oUserKnownHostsFile=/dev/null",
		"-P 2200",
		"../test_fixture/bigfile.linuxiso",
		"test_user@127.0.0.1:/tmp/bigfile.linuxiso"}

	cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, "derp'd it: ", err)
	}
	fmt.Println(string(cmdOut))
}

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

	// scp -Bp -i testing_key.rsa -o StrictHostKeyChecking=no -P 2020 linux.iso test_user@127.0.0.1:tacos
	cmdName := "scp"
	cmdArgs := []string{"-l1000", "-B", "-p", "-i/Users/rely/Projects/bevy_pg/test_fixture/testing_key.rsa", "-oStrictHostKeyChecking=no", "-P 2020", "linux.iso", "test_user@127.0.0.1:target.file"}

	for i := 0; i < 1; i++ {
		cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput()
		if err != nil {
			fmt.Fprintln(os.Stderr, "derp'd it: ", err)
		}
		fmt.Println(string(cmdOut))
	}
}

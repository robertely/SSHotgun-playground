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
	cmdName := "printenv"
	// cmdArgs := []string{""}
	if cmdOut, err = exec.Command(cmdName).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "derp'd it: ", err)
		os.Exit(1)
	}

	fmt.Println(string(cmdOut))
}

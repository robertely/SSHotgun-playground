package main

import (
	"fmt"
)

func main() {

	t2200 := NewTarget(TargetOptions{
		Username:   "test_user",
		Hostname:   "127.0.0.1",
		Port:       2200,
		SSHOptions: []string{"-itest_fixture/testing_key.rsa", "-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null"},
	})

	t2200.controlMaster.Open()
	defer t2200.controlMaster.Close()

	// read the logs
	go func() {
		for i := range t2200.logs {
			fmt.Println(i)
		}
	}()

	// polls and blocks waiting for a ready state.
	t2200.controlMaster.BReady()
	fmt.Println("--- Online ---")
	fmt.Println(t2200.GetRemoteTemp())
	t2200.SendCommand([]string{"ls -hal /tmp/"})
}

package main

import (
	"fmt"
)

func main() {

	t := TargetOptions{
		Username:   "test_user",
		Hostname:   "127.0.0.1",
		Port:       2200,
		SSHOptions: []string{"-itest_fixture/testing_key.rsa", "-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null"},
	}
	t2200 := NewTarget(t)

	t2200.controlMaster.Open()
	defer t2200.controlMaster.Close()

	// not working at all
	go func() {
		for i := range t2200.controlMaster.logs {
			fmt.Println(i)
		}
	}()

	t2200.controlMaster.BReady()
	fmt.Println("Online")

	// time.Sleep(5 * time.Second)
}

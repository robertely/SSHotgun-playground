package main

import (
	// "fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)

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
			log.Info(i)
		}
	}()

	// polls and blocks waiting for a ready state.
	t2200.controlMaster.BReady()
	// Create remote tempdir
	// _ = t2200.GetRemoteTemp()
	t2200.SendCommand([]string{"echo", "Hello", "World"})
	t2200.SendCommand([]string{"whoami"})
	time.Sleep(3 * time.Second)
	t2200.controlMaster.Close()
	time.Sleep(3 * time.Second)
}

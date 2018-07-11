package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	t2200 := NewTarget(TargetOptions{
		Username: "test_user",
		Hostname: "127.0.0.1",
		Port:     2200,
		SSHOptions: []string{"-oStrictHostKeyChecking=no",
			// "-itest_fixture/testing_key.rsa",
			"-oUserKnownHostsFile=/dev/null"},
	})
	t2200.controlMaster.usePty = true
	// x := regexp.MustCompile("a")
	// t2200.controlMaster.expecters = []*Expecter{x}
	t2200.controlMaster.Open()
	defer t2200.controlMaster.Exit()

	// read the logs
	go func() {
		for i := range t2200.logs {
			fmt.Println(i)
		}
	}()

	// polls and blocks waiting for a ready state.
	err := t2200.controlMaster.BlockingReady(10 * time.Second)
	if err != nil {
		log.Error(err)
		t2200.controlMaster.Kill()
		defer os.Exit(5)
		runtime.Goexit()
	}
	// Create remote tempdir
	// _ = t2200.GetRemoteTemp()
	t2200.SendCommand([]string{"echo", "Some Error"})
	t2200.SendCommand([]string{"sudo", "-S", "whoami"})
	time.Sleep(3 * time.Second)
}

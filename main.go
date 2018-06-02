package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type Log struct {
	Origin *Target
	Msg    string
	RxTime time.Time
	// I want to indicate what sent this log. SCP/Controlmaster/SSHCommand
	Process string
	// What's the best way to indicate one of [stderr, stdout, combined]
	pipe string
}

// Target holds tacos
type Target struct {
	Host          string
	Port          int
	User          string
	ControlMaster ControlMaster
}

func main() {
	master := NewControlMaster("192.168.1.238", "22", "sshotgun-%h-%p-%C.sock")
	master.Open()
	defer master.Close()
	master.BReady()

	scanner := bufio.NewScanner(master.ptmx)
	for scanner.Scan() {
		fmt.Fprintf(os.Stdout, "\033[0;36mRx\033[0m: %s\n", scanner.Text())
	}
}

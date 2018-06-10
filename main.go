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
	Source string
	// What's the best way to indicate one of [stderr, stdout, combined]
	Type string
}

func (l Log) String() string {
	return l.Msg
}

// Target holds tacos
type Target struct {
	Host   string
	Port   int
	User   string
	Logs   chan string
	CM     ControlMaster
	SSHOps []string
}

// NewTarget Makes new TARGET BROOOOWWNNN
func NewTarget(u, h string, p, lbuf int, sshops []string) *Target {
	return &Target{
		Host:   h,
		Port:   p,
		User:   u,
		Logs:   make(chan string, lbuf),
		CM:     NewControlMaster(u, h, p, sshops),
		SSHOps: sshops}
}

func (t *Target) Connect() {
	t.CM.Open()
	t.CM.BReady()
}

func (t *Target) SendFile(p string) {
	fmt.Println("Send", p)
}

// func LocalExec
// TODO: Soething about the env.
// fwict a "command" is a combination of:
//   path to exec (string)
//   arguments ([]string)
//   environment([]string)
func (t *Target) Execute(cmd []string, pty bool) {
	fmt.Println("Execute", cmd, pty)
}

func main() {
	ops := []string{"-itest_fixture/testing_key.rsa", "-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null"}
	t2200 := NewTarget("test_user", "127.0.0.1", 2200, 1000, ops)
	scanner := bufio.NewScanner(t2200.CM.ptmx)
	for scanner.Scan() {
		fmt.Fprintf(os.Stdout, "\033[0;36mRx\033[0m: %s\n", scanner.Text())
	}
	t2200.CM.Open()
	t2200.CM.BReady()
	fmt.Println("Online")
	time.Sleep(5 * time.Second)
	// master.Open()
	// defer master.Close()
	// master.BReady()
	//
}

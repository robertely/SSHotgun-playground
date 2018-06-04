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
	Host string
	Port int
	User string
	Logs chan string
	CM   ControlMaster
}

// NewTarget Makes new TARGET BROOOOWWNNN
func (t Target) NewTarget(h, u string, p, lbuf int) {
	t.User = u
	t.Host = h
	t.Port = p

	t.Logs = make(chan string, lbuf)
	t.CM = NewControlMaster(t.User, t.Host, t.Port)
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
	// master := NewControlMaster("192.168.1.238", "22", "sshotgun-%h-%p-%C.sock")
	// master.Open()
	// defer master.Close()
	// master.BReady()
	//
	// scanner := bufio.NewScanner(master.ptmx)
	// for scanner.Scan() {
	// 	fmt.Fprintf(os.Stdout, "\033[0;36mRx\033[0m: %s\n", scanner.Text())
	// }
}

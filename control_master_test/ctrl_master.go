package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/kr/pretty"
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

// Terms:
//   pty  - pseudo terminal
//   pts  - pty slave
//   ptmx - pty master

var eot = []byte{4} // End Of Transmission

// ControlMaster holds connection details of the master ssh connection
type ControlMaster struct {
	cmd        *exec.Cmd
	ptmx       *os.File
	ptySize    pty.Winsize
	socketPath string
	targetHost string
	targetPort string
}

// NewControlMaster - ControlMaster consturctor
func NewControlMaster(tHost, tPort, socketPath string) *ControlMaster {
	cm := new(ControlMaster)
	cm.targetHost = tHost
	cm.targetPort = tPort
	cm.ptySize = pty.Winsize{Rows: 24, Cols: 80, X: 1024, Y: 768}
	cm.cmd = exec.Command("../ipsum", "ssh", "-M", "-N", "-oControlPath="+socketPath, tHost, "-p", tPort)
	// cm.cmd = exec.Command("ssh", "-M", "-N", "-oControlPath="+socketPath, tHost, "-p", tPort)
	cm.socketPath = socketPath
	return cm
}

// Open - starts ssh with control master configuration
func (cm *ControlMaster) Open() {
	var err error // K. https://github.com/golang/go/issues/6842
	cm.ptmx, err = pty.Start(cm.cmd)
	if err != nil {
		panic(err)
	}
	// Initialize ...
	pty.Setsize(cm.ptmx, &cm.ptySize)
	terminal.MakeRaw(int(cm.ptmx.Fd()))
}

func (cm *ControlMaster) Close() {
	cm.ptmx.Close()
}

func (cm ControlMaster) Send(s string) {
	cm.ptmx.Write([]byte(s))
}

// Check - ssh ctl_cmd
// (check that the master process is running)
func (cm ControlMaster) Check() {
	fmt.Println("Check")
}

// Exit - ssh ctl_cmd
// (request the master to	exit)
func (cm ControlMaster) Exit() {
	fmt.Println("Exit")
}

// Stop - ssh ctl_cmd
// (request the master to stop accepting further multiplexing requests)
func (cm ControlMaster) Stop() {
	fmt.Println("Stop")
}

// Kill - Signal sigKill to the ControlMaster process
//   sigKill ssh control master
func (cm ControlMaster) Kill() {
	cm.cmd.Process.Signal(os.Kill)
}

func main() {
	master := NewControlMaster("192.168.1.238", "22", "sshotgun-%h-%p-%C.sock")
	master.Open()
	defer master.Close()
	fmt.Printf("%# v", pretty.Formatter(master))

	master.Send("echo 'Should Only see this message once.'\n")

	scanner := bufio.NewScanner(master.ptmx)
	for scanner.Scan() {
		fmt.Fprintf(os.Stdout, "\033[0;36mRx\033[0m: %s\n", scanner.Text())
	}
}

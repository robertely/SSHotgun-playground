package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

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
	targetPort int
	logs       chan string
	fqcp       string
	sshops     []string
}

// NewControlMaster - ControlMaster consturctor
func NewControlMaster(tUser, tHost string, tPort int, sshops []string) ControlMaster {
	name := "ssh"
	fqcp := tUser + "@" + tHost
	// socketPath := "sshotgun-%h-%p.sock"
	cm := ControlMaster{}
	cm.sshops = sshops
	cm.targetHost = tHost
	cm.fqcp = fqcp
	cm.targetPort = tPort
	cm.ptySize = pty.Winsize{Rows: 24, Cols: 80, X: 1024, Y: 768}
	cm.socketPath = fmt.Sprintf("sshotgun-%s-%s-%s.sock", tHost, tUser, strconv.Itoa(cm.targetPort))
	cm.logs = make(chan string, 1000) // TODO: not this.
	args := append([]string{"-M", "-N", "-oControlPath=" + cm.socketPath, fqcp, "-p", strconv.Itoa(cm.targetPort)}, sshops...)

	cm.cmd = exec.Command(name, args...)
	// fmt.Println(name, args)

	// go func() {
	// 	outPipe, _ := action.cmd.StdoutPipe()
	// 	outScanner := bufio.NewScanner(outPipe)
	// 	for outScanner.Scan() {
	// 		action.Logs <- outScanner.Text()
	// 	}
	// }()
	return cm
}

// Open - starts ssh with control master configuration
func (cm ControlMaster) Open() {
	var err error // K. https://github.com/golang/go/issues/6842
	cm.ptmx, err = pty.Start(cm.cmd)
	if err != nil {
		panic(err)
	}
	// Initialize ...
	pty.Setsize(cm.ptmx, &cm.ptySize)
	terminal.MakeRaw(int(cm.ptmx.Fd()))
}

func (cm ControlMaster) sendCtrlCmd(ctrlcmd string) string {
	name := "ssh"
	args := append([]string{"-oControlPath=" + cm.socketPath, cm.fqcp, "-p", strconv.Itoa(cm.targetPort), "-O", ctrlcmd}, cm.sshops...)
	cmd := exec.Command(name, args...)
	// fmt.Println(name, args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Printf("%s\n", out)
		return ""
	}
	return string(out)
}

func (cm ControlMaster) Close() {
	cm.ptmx.Close()
}

func (cm ControlMaster) Send(s string) {
	cm.ptmx.Write([]byte(s))
}

// Kill - Signal sigKill to the ControlMaster process
//   sigKill ssh control master
func (cm ControlMaster) Kill() {
	cm.cmd.Process.Signal(os.Kill)
}

// Ready - ssh ctl_cmd
// (check that the master process is running)
func (cm ControlMaster) Ready() bool {
	if _, err := os.Stat(cm.socketPath); err == nil {
		stdout := cm.sendCtrlCmd("check")
		if strings.HasPrefix(string(stdout), "Master running") {
			return true
		}
		return false
	}
	fmt.Println(cm.socketPath, "Waiting for control master socket...")
	return false
}

func (cm ControlMaster) BReady() {
	for !cm.Ready() {
		time.Sleep(100 * time.Millisecond)
	}
	return
}

// Exit - ssh ctl_cmd
// (request the master to	exit)
func (cm ControlMaster) Exit() bool {
	stdout := cm.sendCtrlCmd("exit")
	if strings.HasPrefix(string(stdout), "Exit request sent.") {
		return true
	}
	return false
}

// Stop - ssh ctl_cmd
// (request the master to stop accepting further multiplexing requests)
func (cm ControlMaster) Stop() bool {
	stdout := cm.sendCtrlCmd("stop")
	fmt.Println(stdout)
	if strings.HasPrefix(stdout, "Stop listening request sent.") {
		return true
	}
	return false
}

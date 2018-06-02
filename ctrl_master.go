package main

import (
	"fmt"
	"os"
	"os/exec"
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
	targetPort string
}

// NewControlMaster - ControlMaster consturctor
func NewControlMaster(tHost, tPort, socketPath string) *ControlMaster {
	name := "ssh"
	args := []string{"-M", "-N", "-oControlPath=" + socketPath, tHost, "-p", tPort} // "-oControlPersist=yes"

	cm := new(ControlMaster)
	cm.targetHost = tHost
	cm.targetPort = tPort
	cm.ptySize = pty.Winsize{Rows: 24, Cols: 80, X: 1024, Y: 768}
	cm.cmd = exec.Command(name, args...)
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

func (cm *ControlMaster) sendCtrlCmd(ctrlcmd string) string {
	name := "ssh"
	args := []string{"-oControlPath=" + cm.socketPath, cm.targetHost, "-p", cm.targetPort, "-O", ctrlcmd}
	cmd := exec.Command(name, args...)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(stdout)
}

func (cm *ControlMaster) Close() {
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
	stdout := cm.sendCtrlCmd("check")
	if strings.HasPrefix(string(stdout), "Master running") {
		return true
	}
	return false
}

func (cm ControlMaster) BReady() {
	for !cm.Ready() {
		time.Sleep(50 * time.Millisecond)
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

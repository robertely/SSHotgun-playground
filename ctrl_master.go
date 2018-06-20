package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	// "strconv"
	"strings"
	"time"
	// "github.com/kr/pty"
)

// Terms:
//   pty  - pseudo terminal
//   pts  - pty slave
//   ptmx - pty master

var eot = []byte{4} // End Of Transmission

// ControlMaster holds connection details of the master ssh connection
type ControlMaster struct {
	cmd *exec.Cmd
	// ptmx   *os.File
	target *Target
	// ptySize    pty.Winsize
	socketPath string
}

// NewControlMaster - ControlMaster consturctor
func NewControlMaster(t *Target) *ControlMaster {
	cm := ControlMaster{}
	cm.target = t
	cm.socketPath = fmt.Sprintf("sshotgun-%%h-%%p-%%r.%s.sock", t.sessionID)
	return &cm
}

// Open - starts ssh with control master configuration
func (cm *ControlMaster) Open() {
	name := cm.target.sshcmd
	args := append(cm.target.CmdBuilder(true), "-M", "-N", "-S", cm.socketPath)
	log.Debug(name, args)
	cm.cmd = exec.Command(name, args...)
	cmdOut, _ := cm.cmd.StdoutPipe()
	go func() {
		outScanner := bufio.NewScanner(cmdOut)
		for outScanner.Scan() {
			cm.target.logs <- Log{
				Origin: cm.target,
				Msg:    outScanner.Text(),
				RxTime: time.Now(),
				Source: "ControlMaster",
				Type:   "stdout"}
		}
	}()

	cmdErr, _ := cm.cmd.StderrPipe()
	go func() {
		errScanner := bufio.NewScanner(cmdErr)
		cm.target.logs <- Log{
			Origin: cm.target,
			Msg:    errScanner.Text(),
			RxTime: time.Now(),
			Source: "ControlMaster",
			Type:   "stderr"}
	}()
	err := cm.cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
	}
	go func() {
		err = cm.cmd.Wait()
		if err != nil {
			log.Error("Control master exited unexpectidly:", err)
		}
	}()
}

func (cm ControlMaster) sendCtrlCmd(ctrlcmd string) string {
	name := cm.target.sshcmd
	args := append(cm.target.CmdBuilder(true), "-O", ctrlcmd)
	log.Debug(name, args)
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Warn(err.Error())
		log.Warn("%s\n", out)
		return ""
	}
	return string(out)
}

func (cm ControlMaster) Close() {
	cm.Exit()
}

// func (cm ControlMaster) Send(s string) {
// 	cm.ptmx.Write([]byte(s))
// }

// Kill - Signal sigKill to the ControlMaster process
//   sigKill ssh control master
func (cm ControlMaster) Kill() {
	cm.cmd.Process.Signal(os.Kill)
}

// Ready - ssh ctl_cmd
// (check that the master process is running)
func (cm ControlMaster) Ready() bool {
	files, _ := filepath.Glob(fmt.Sprintf("*.%s.sock", cm.target.sessionID))
	if len(files) == 0 {
		log.Info("ControlMaster Socket does not yet exist")
		return false
	}
	stdout := cm.sendCtrlCmd("check")
	// fmt.Println(stdout)
	if strings.HasPrefix(string(stdout), "Master running") {
		return true
	}
	return false
}

func (cm ControlMaster) BReady() {
	for !cm.Ready() {
		time.Sleep(250 * time.Millisecond)
	}
	return
}

// Exit - ssh ctl_cmd
// (request the master to exit)
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
	if strings.HasPrefix(stdout, "Stop listening request sent.") {
		return true
	}
	return false
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
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
	stdin  io.Writer
	// ptySize    pty.Winsize
	socketPath string
	running    bool
	expectExit bool
}

// NewControlMaster - ControlMaster consturctor
func NewControlMaster(t *Target) *ControlMaster {
	cm := ControlMaster{}
	cm.target = t
	cm.socketPath = fmt.Sprintf("bevy-%%h-%%p-%%r.%s.sock", t.sessionID)
	return &cm
}

// Open - starts ssh with control master configuration
func (cm *ControlMaster) Open() {
	name := cm.target.sshcmd
	args := append(cm.target.CmdBuilder(true), "-M", "-N")
	log.Debug(name, args)
	cm.cmd = exec.Command(name, args...)
	cmdOut, _ := cm.cmd.StdoutPipe()
	go func() {
		outScanner := bufio.NewScanner(cmdOut)
		for outScanner.Scan() {
			cm.target.logs <- Log{
				Origin:  cm.target,
				Msg:     outScanner.Text(),
				RxTime:  time.Now(),
				Source:  "Master",
				Context: strings.Join(cm.cmd.Args, " "),
				Stream:  "stdout"}
		}
	}()

	cmdErr, _ := cm.cmd.StderrPipe()
	go func() {
		errScanner := bufio.NewScanner(cmdErr)
		for errScanner.Scan() {
			cm.target.logs <- Log{
				Origin:  cm.target,
				Msg:     errScanner.Text(),
				RxTime:  time.Now(),
				Source:  "Master",
				Context: strings.Join(cm.cmd.Args, " "),
				Stream:  "stderr"}
		}
	}()
	cm.stdin, _ = cm.cmd.StdinPipe()

	cm.running = true
	err := cm.cmd.Start()
	if err != nil {
		log.Error("Could not start ControlMaster connection: ", err)
	}
	go func() {
		err = cm.cmd.Wait()
		cm.running = false
		if err != nil && !cm.expectExit {
			log.Error("Control master exited unexpectidly: ", err)
		}
	}()
}

func (cm *ControlMaster) sendCtrlCmd(ctrlcmd string) string {
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

// Send - string on stdin to ControlMaster process
func (cm ControlMaster) Send(s string) {
	io.WriteString(cm.stdin, s)
}

// Kill - Signal sigKill to the ControlMaster process
//   sigKill ssh control master
func (cm *ControlMaster) Kill() {
	cm.expectExit = true
	cm.cmd.Process.Signal(os.Kill)
}

// Ready - ssh ctl_cmd
// check that the master process is running and prepared to accept connections
func (cm *ControlMaster) Ready() bool {
	if !cm.running {
		return false
	}
	files, _ := filepath.Glob(fmt.Sprintf("*.%s.sock", cm.target.sessionID))
	if len(files) == 0 {
		log.Debug("ControlMaster Socket does not yet exist")
		return false
	}
	stdout := cm.sendCtrlCmd("check")
	// fmt.Println(stdout)
	if strings.HasPrefix(string(stdout), "Master running") {
		return true
	}
	return false
}

// BlockingReady - Blocks and polls waiting for control master to come up.
// time out specifies a time to wait. (imperfect but near enogh)
func (cm *ControlMaster) BlockingReady(timeout time.Duration) error {
	log.Info("Waiting for control master...")
	start := time.Now()
	// should cm.Ready return error type? should that bubble up ?
	for !cm.Ready() {
		time.Sleep(250 * time.Millisecond)
		if time.Now().After(start.Add(timeout)) {
			return errors.New("Exceded timeout waiting for ControlMaster Ready")
		}
	}
	return nil
}

// Exit - ssh ctl_cmd
// (request the master to exit)
func (cm *ControlMaster) Exit() error {
	if !cm.running {
		log.Warn("ControlMaster already exited")
		return nil
	}
	cm.expectExit = true
	stdout := cm.sendCtrlCmd("exit")
	if strings.HasPrefix(string(stdout), "Exit request sent.") {
		log.Debug("ControlMaster accepted exit request")
		return nil
	}
	return errors.New("ControlMaster Exit request failed")
}

// Stop - ssh ctl_cmd
// (request the master to stop accepting further multiplexing requests)
func (cm *ControlMaster) Stop() bool {
	stdout := cm.sendCtrlCmd("stop")
	if strings.HasPrefix(stdout, "Stop listening request sent.") {
		return true
	}
	return false
}

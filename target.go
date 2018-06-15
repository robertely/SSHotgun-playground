package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Target struct {
	username      string
	hostname      string
	port          int
	logs          chan Log
	sshOptions    []string
	controlMaster *ControlMaster
	sessionID     string
}

type TargetOptions struct {
	Username      string
	Hostname      string
	Port          int
	Logs          chan Log
	LogLen        int
	SSHOptions    []string
	ControlMaster ControlMaster
}

// Create a new Target, configured with the given TargetOptions
func NewTarget(o TargetOptions) *Target {
	// explicitly copy options, one by one, providing sane defaults where possible

	target := Target{
		username:   o.Username,
		hostname:   o.Hostname,
		logs:       make(chan Log, o.LogLen),
		sshOptions: o.SSHOptions,
	}

	// append port if specificed
	if o.Port != 0 {
		target.sshOptions = append(target.sshOptions, "-p", strconv.Itoa(o.Port))
	}
	// append hostname if specificed
	if o.Username == "" {
		target.sshOptions = append(target.sshOptions, o.Hostname)
	} else {
		target.sshOptions = append(target.sshOptions, o.Username+"@"+o.Hostname)
	}
	// Initialize control master
	// if o.ControlMaster != (ControlMaster{}) {
	// 	// use the ControlMaster that we were given
	// 	target.controlMaster = o.ControlMaster
	// } else {
	// create a control master, bound to this target
	target.controlMaster = NewControlMaster(&target)
	// }
	target.sessionID = target.makeSessionId()
	return &target
}

func (t *Target) makeSessionId() string {
	s := md5.New()
	str := t.username + t.hostname + strconv.Itoa(t.port) + strconv.Itoa(int(time.Now().UnixNano()))
	s.Write([]byte(str))
	return fmt.Sprintf("%X", s.Sum(nil)[:])
}

func (t *Target) SendCommand(s []string) {
	name := "ssh"
	args := t.sshOptions
	if t.controlMaster.Ready() {
		args = append(args, "-oControlPath="+t.controlMaster.socketPath)
	}
	args = append(args, s...) // oh man is this wrong
	fmt.Println(name, args)
	cmd := exec.Command(name, args...)

	cmdOut, _ := cmd.StdoutPipe()
	go func() {
		outScanner := bufio.NewScanner(cmdOut)
		for outScanner.Scan() {
			t.logs <- Log{
				Origin: t,
				Msg:    outScanner.Text(),
				RxTime: time.Now(),
				Source: "Command",
				Type:   "stdout"}
		}
	}()

	cmdErr, _ := cmd.StderrPipe()
	go func() {
		errScanner := bufio.NewScanner(cmdErr)
		t.logs <- Log{
			Origin: t,
			Msg:    errScanner.Text(),
			RxTime: time.Now(),
			Source: "Command",
			Type:   "stderr"}
	}()
	err := cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
	}
	_ = cmd.Wait()
}

func (t *Target) GetRemoteTemp() string {
	name := "ssh"
	args := t.sshOptions
	if t.controlMaster.Ready() {
		args = append(args, "-oControlPath="+t.controlMaster.socketPath)
	}
	args = append(args, "mktemp -d -t .Bevy.XXXX."+t.sessionID)
	fmt.Println(name, args)
	cmd := exec.Command(name, args...)
	cmdOut, _ := cmd.CombinedOutput()
	return strings.TrimRight(string(cmdOut[:]), "\n")
}

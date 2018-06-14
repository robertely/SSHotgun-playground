package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type Target struct {
	username      string
	hostname      string
	port          int
	logs          chan string
	sshOptions    []string
	controlMaster *ControlMaster
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
		logs:       make(chan string, o.LogLen),
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

	return &target
}

func (t *Target) SendCommand(s []string) {
	name := "ssh"
	args := t.sshOptions
	if t.controlMaster.Ready() {
		args = append(args, "-oControlPath="+t.controlMaster.socketPath)
	}
	args = append(args, s[0])
	fmt.Println(name, args)
	cmd := exec.Command(name, args...)

	cmdOut, _ := cmd.StdoutPipe()
	go func() {
		outScanner := bufio.NewScanner(cmdOut)
		for outScanner.Scan() {
			t.logs <- outScanner.Text()
		}
	}()

	cmdErr, _ := cmd.StderrPipe()
	go func() {
		errScanner := bufio.NewScanner(cmdErr)
		for errScanner.Scan() {
			t.logs <- errScanner.Text()
		}
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
	args = append(args, "mktemp -d -t .Bevy.$(date +%s).XXXXX")
	fmt.Println(name, args)
	cmd := exec.Command(name, args...)
	cmdOut, _ := cmd.CombinedOutput()
	return string(cmdOut[:])
}

// func (cm ControlMaster) sendCtrlCmd(ctrlcmd string) string {
// 	name := "ssh"
// 	args := append([]string{"-O", ctrlcmd}, cm.target.sshOptions...)
// 	fmt.Println(name, args)
// 	cmd := exec.Command(name, args...)
// 	// fmt.Println(name, args)
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		fmt.Printf("%s\n", out)
// 		return ""
// 	}
// 	return string(out)
// }

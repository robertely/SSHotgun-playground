package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Target struct {
	username      string
	hostname      string
	password      string
	sudopassword  string
	sshcmd        string
	port          int
	logs          chan Log
	sshOptions    []string
	controlMaster *ControlMaster
	sessionID     string
}

type TargetOptions struct {
	Username      string
	Hostname      string
	Password      string
	SudoPassword  string
	SSHCMD        string
	Port          int
	Logs          chan Log
	LogLen        int
	SSHOptions    []string
	ControlMaster *ControlMaster
}

// NewTarget Creates a new Target, configured with the given TargetOptions
func NewTarget(o TargetOptions) *Target {
	// explicitly copy options, one by one, providing sane defaults where possible

	target := Target{
		username:     o.Username,
		hostname:     o.Hostname,
		port:         o.Port,
		logs:         make(chan Log, o.LogLen),
		sshOptions:   o.SSHOptions,
		password:     o.Password,
		sudopassword: o.SudoPassword,
	}
	target.sessionID = target.makeSessionId()

	if o.SSHCMD == "" {
		target.sshcmd = "ssh"
	} else {
		target.sshcmd = o.SSHCMD
	}
	// Initialize control master
	if o.ControlMaster != nil {
		target.controlMaster = o.ControlMaster
	} else {
		target.controlMaster = NewControlMaster(&target)
	}
	return &target
}

func (t *Target) CmdBuilder(useCM bool) []string {
	result := t.sshOptions
	// append port if specificed
	if t.port != 0 {
		result = append(result, "-p", strconv.Itoa(t.port))
	}
	// Add username if specificed
	if t.username == "" {
		result = append(result, t.hostname)
	} else {
		result = append(result, t.username+"@"+t.hostname)
	}
	// Use controlmaster if ready
	if useCM {
		result = append(result, "-oControlPath="+t.controlMaster.socketPath)
	}

	return result
}

func (t *Target) makeSessionId() string {
	s := md5.New()
	str := t.username + t.hostname + strconv.Itoa(t.port) + strconv.Itoa(int(time.Now().UnixNano()))
	s.Write([]byte(str))
	return fmt.Sprintf("%X", s.Sum(nil)[:])
}

func (t *Target) SendCommand(s []string) {
	name := t.sshcmd
	args := t.CmdBuilder(t.controlMaster.Ready())
	args = append(args, s...) // oh man is this wrong
	log.Debug(name, args)
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
	name := t.sshcmd
	args := t.CmdBuilder(t.controlMaster.Ready())

	args = append(args, "mktemp -d -t .Bevy.XXXX."+t.sessionID)
	log.Debug(name, args)
	cmd := exec.Command(name, args...)
	cmdOut, _ := cmd.CombinedOutput()
	return strings.TrimRight(string(cmdOut[:]), "\n")
}

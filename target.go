package main

import (
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
	Logs          chan string
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

package main

import (
	"fmt"
	"time"
)

type Log struct {
	Origin  *Target
	Msg     string
	RxTime  time.Time
	Source  string
	Context string
	Stream  string
}

func (l Log) String() string {
	var s string
	if l.Stream == "stderr" {
		s += "[ERR]"
	} else if l.Stream == "stdout" {
		s += "[OUT]"
	} else {
		s += fmt.Sprintf("[%s]", l.Stream)
	}
	s += fmt.Sprintf("[%s]", l.Origin.sessionID)
	s += fmt.Sprintf("[%s]", l.Origin.hostname)
	s += fmt.Sprintf("[%s]", l.Source)
	s += fmt.Sprintf(": %s", l.Msg)
	return s
}

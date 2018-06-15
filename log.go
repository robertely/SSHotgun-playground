package main

import (
	"fmt"
	"time"
)

type Log struct {
	Origin *Target
	Msg    string
	RxTime time.Time
	Source string
	Type   string
}

func (l Log) String() string {
	s := fmt.Sprintf("%s\n", l.Origin.sessionID)
	s += fmt.Sprintf("  source: '%s'\n", l.Source)
	s += fmt.Sprintf("  type: '%s'\n", l.Type)
	s += fmt.Sprintf("  RxTime: '%v'\n", l.RxTime.Unix())
	s += fmt.Sprintf("  msg: '%s'\n", l.Msg)
	return s
}

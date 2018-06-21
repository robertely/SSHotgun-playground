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
	s := fmt.Sprintf("Id: %s, ", l.Origin.sessionID)
	s += fmt.Sprintf("source: '%s', ", l.Source)
	s += fmt.Sprintf("stream: '%s', ", l.Stream)
	s += fmt.Sprintf("context: '%s', ", l.Context)
	s += fmt.Sprintf("RxTime: '%v', ", l.RxTime.Unix())
	s += fmt.Sprintf("msg: '%s'}", l.Msg)
	return s
}

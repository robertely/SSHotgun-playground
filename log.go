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
	s := fmt.Sprintf("Id: %s, ", l.Origin.sessionID)
	s += fmt.Sprintf("source: '%s', ", l.Source)
	s += fmt.Sprintf("type: '%s', ", l.Type)
	s += fmt.Sprintf("RxTime: '%v', ", l.RxTime.Unix())
	s += fmt.Sprintf("msg: '%s'}", l.Msg)
	return s
}

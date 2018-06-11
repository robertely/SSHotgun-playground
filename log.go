package main

import "time"

type Log struct {
	Origin *Target
	Msg    string
	RxTime time.Time
	Source string
	Type   string
}

func (l Log) String() string {
	return l.Msg
}

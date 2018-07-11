package main

import "regexp"

type Expecter struct {
	Expression *regexp.Regexp
	Action     interface{} // a function, or a string
	Uses       int
}

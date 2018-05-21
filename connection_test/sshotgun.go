package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

var eot = []byte{4}

func main() {
	c0 := exec.Command("ssh", "192.168.1.238")
	f0, _ := pty.Start(c0)
	// Make sure to close the pty at the end.
	defer func() { _ = f0.Close() }() // Best effort.

	ws := pty.Winsize{Rows: 24, Cols: 80, X: 1024, Y: 768}
	pty.Setsize(f0, &ws)
	terminal.MakeRaw(int(f0.Fd()))

	f0.Write([]byte("stty -echo\n"))
	f0.Write([]byte("export PS1=''\n\n\n\n"))
	f0.Write([]byte("uptime\n"))
	f0.Write([]byte("echo 'Should Only see this message once.'\n"))
	f0.Write(eot)

	scanner := bufio.NewScanner(f0)
	for scanner.Scan() {
		fmt.Fprintf(os.Stdout, "\033[0;36mRx\033[0m: %s\n", scanner.Text())
	}
}

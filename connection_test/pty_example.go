package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

var eot = []byte{4}

func main() {
	// c0 := exec.Command("ssh", "-q", "-i/Users/rely/Projects/bevy_pg/test_fixture/testing_key.rsa", "-oStrictHostKeyChecking=no", "-p2020", "test_user@127.0.0.1")
	cmdName := "scp"
	cmdArgs := []string{"-l 1000", "-B", "-p", "-i/Users/rely/Projects/bevy_pg/test_fixture/testing_key.rsa", "-oStrictHostKeyChecking=no", "-P 2020", "linux.iso", "test_user@127.0.0.1:target.file"}
	c0 := exec.Command(cmdName, cmdArgs...)
	f0, _ := pty.Start(c0)
	// Make sure to close the pty at the end.
	defer func() { _ = f0.Close() }() // Best effort.

	ws := pty.Winsize{Rows: 24, Cols: 80, X: 1024, Y: 768}
	pty.Setsize(f0, &ws)
	terminal.MakeRaw(int(f0.Fd()))

	go func() {
		scanner := bufio.NewScanner(f0)
		for scanner.Scan() {
			fmt.Fprintf(os.Stdout, "\033[0;36mRx\033[0m: %s\n", scanner.Text())
		}
	}()

	// f0.Write([]byte("stty -echo\n"))
	// f0.Write([]byte("export PS1=''\n\n\n\n"))
	// f0.Write([]byte("uptime\n"))
	// f0.Write([]byte("echo 'Should Only see this message once.'\n"))
	// f0.Write([]byte("printenv\n"))
	// f0.Write(eot)
	fmt.Println("Sent EOT")
	time.Sleep(200 * time.Second)
}

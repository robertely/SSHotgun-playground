package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/kr/pty"
)

var eot = []byte{4}

func main() {
	cmdName := "ssh"
	cmdArgs := []string{"-oStrictHostKeyChecking=no", "-oUserKnownHostsFile=/dev/null", "-S", "bevy-%h-%p-%r.6F3E41E23B21.sock", "-p", "2200", "test_user@127.0.0.1", "-M", "-N"}
	cmd := exec.Command(cmdName, cmdArgs...)
	fmt.Println(cmdArgs)
	ptmx, _ := pty.Start(cmd)

	// Make sure to close the pty at the end.
	// defer func() { _ = f0.Close() }() // Best effort.
	// ws := pty.Winsize{Rows: 24, Cols: 80, X: 1024, Y: 768}
	// pty.Setsize(f0, &ws)
	// terminal.MakeRaw(int(ptmx.Fd()))

	// go func() {
	// 	_, _ = io.Copy(ptmx, os.Stdin)
	// }()

	go func() {
		scanner := bufio.NewScanner(ptmx)
		for scanner.Scan() {
			fmt.Fprintf(os.Stdout, "\033[0;36mRx\033[0m: %s\n", scanner.Text())
		}
	}()

	time.Sleep(time.Second * 3)
	ptmx.Write([]byte("\n"))
	time.Sleep(20 * time.Second)
}

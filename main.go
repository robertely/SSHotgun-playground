package main

import (
	"fmt"
	"time"
)

func main() {
	master := NewControlMaster("192.168.1.238", "22", "sshotgun-%h-%p-%C.sock")
	master.Open()
	defer master.Close()
	defer master.Kill()

	time.Sleep(3 * time.Second)
	if master.Check() {
		fmt.Println("READY")
	} else {
		fmt.Println("NOT READY")
	}

	// scanner := bufio.NewScanner(master.ptmx)
	// for scanner.Scan() {
	// 	fmt.Fprintf(os.Stdout, "\033[0;36mRx\033[0m: %s\n", scanner.Text())
	// }
}
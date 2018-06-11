### Multiplexing:
* First openssh release supporting multiplexing [_OpenSSH 3.9/3.9p1 (2004-08-18)_](https://www.openssh.com/txt/release-3.9)
  * Doesn't support `ControlMaster=auto` but we don't need it if we manage the sockets our self.
  * Centos newer than 5.x needed
  * We should manage the sockets our self just incase.
* Limited to 10 simultaneous connections by default as of:
  * [ OpenSSH 5.1/5.1p1 (2008-07-22) ] (https://www.openssh.com/txt/release-5.1)

### get to know pty:
* /proc/sys/kernel/pty/{max, nr, reserve} (osx?)
* ~https://github.com/google/goexpect~
* https://github.com/kr/pty
* https://unix.stackexchange.com/questions/406911/why-does-ssh-utility-considered-a-pty
* https://unix.stackexchange.com/questions/117981/what-are-the-responsibilities-of-each-pseudo-terminal-pty-component-software
* http://www.linusakesson.net/programming/tty/
* useful commands for monitoring:
  * `watch -n 1 -d "ls -1 *.sock"`
  * `watch -n 1 lsof +D /dev/pts/`
  * `watch -n 1 pstree -lA '\`pidof control_master_test\`'`
* FDs are going to be a limit. A dangerous one.
* Prometheus seems to have some kind of awareness.
  * `level=info ts=2018-06-05T18:49:29.503230907Z caller=main.go:223 fd_limits="(soft=7168, hard=9223372036854775807)"``
    * Never mind. It reports, but that's all. It doesn't do shit with it:
      * https://github.com/prometheus/prometheus/search?utf8=%E2%9C%93&q=FdLimits&type=

### Preflight script
* stick to posix tools `#!/bin/sh`
* create temproary directory `mktemp`
* make sure it's writeable
* report full path to temporary directory back to SSHotgun
* smth smth...

```
#!/bin/sh

stty rows 24
stty cols 80

RUNNING_USER=$(id -u -n)
WORK_DIR=$(mktemp -d -t sshotgun.$RUNNING_USER.XXXXXX)
cd $WORK_DIR
echo $WORK_DIR
```

### File transfer
* Fuck man i may use scp
* https://unix.stackexchange.com/questions/8707/whats-the-difference-between-sftp-scp-and-fish-protocols/116691#116691
* https://rsync.samba.org/how-rsync-works.html

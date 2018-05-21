### Supported ssh versions: 
[_OpenSSH 3.9/3.9p1 (2004-08-18)_](https://www.openssh.com/txt/release-3.9)
* First release supporting multiplexing
* Functionally this kills compatability for old versions of centos 5
* Doesn't support `ControlMaster=auto` but we don't need it if we manage the sockets our self.
* We should manage the sockets our self just incase.
* There are a number of bugs and fixes to this mode after this release, but this is I think the bare minimum.

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


### Preflight script
* stick to posix tools `#!/bin/sh`
* create temproary directory `mktemp`
* make sure it's writeable
* report full path to temporary directory back to SSHotgun

### File transfer
* Fuck man i may use scp
* https://unix.stackexchange.com/questions/8707/whats-the-difference-between-sftp-scp-and-fish-protocols/116691#116691
* https://rsync.samba.org/how-rsync-works.html

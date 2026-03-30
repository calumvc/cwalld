
# Chinese Wall Security Daemon for SELinux

A security daemon capable of enforcing the Chinese Wall (Brewer-Nash) Model by utilising Security Enhanced Linux's mandatory access control system

## Requirements
- SELinux system (Red Hat / Fedora / ...)
- Go 1.25.7+
- policycoreutils-devel package

## Installation

### 0. ssh reverse tunnel (optional)
from host
```sudo systemctl start sshd```
from VM (10.0.2.2 is standard for VM to host ip)
```ssh -R 2222:localhost:22 host_username@10.0.2.2``` 
from host
```ssh -p 2222 vm_username@localhost```

### 1. Unzip cwalld.zip 
Run `unzip cwalld.zip`

### 2. Install Policies and Integration test daemons compilation
The policy install script is currently tied to a testgrounds directory. 
This also means that discretionary access controls don't interfere.
```
sudo mkdir /home/testgrounds
sudo chmod 755 /home/testgrounds
sudo cp -r testgrounds /home/
```

Finally, run the policy install script (takes a minute)
from `/home/testgrounds/`
```
sudo chmod +x policyinstall.sh
sudo ./policyinstall.sh

sudo chmod +x daemoninstall.sh
sudo ./daemoninstall.sh
```

### 3. Build Executables
from `cwalld/`
```
go build ./cmd/cwalld/cwalld-enforce

sudo cp cwalld-enforce /usr/local/sbin
sudo cp cwalld-enforce.service /usr/lib/systemd/system
sudo systemctl daemon-reload
```

### 4. Prepare Types

from `home/testgrounds/`
```
sudo chmod +x labelinstall.sh
sudo ./labelinstall.sh
```

### 5. Run & Tail
We need SELinux to be on `Enforcing`, so run `sudo getenforce`

If `Enforcing`, we're ready to go. 

If `Permissive`, run `sudo setenforce 1` and check again

from `cwalld/`
```
sudo go run ./cmd/cwalld/cwalld-init
```
to initialise the auditd rule and change the AVC cache for denial logs
running cwalld-init is persistent and should change audit rules forever on the system

```
sudo systemctl start cwalld-enforce
```

Run any desired subject with `sudo systemctl start exampled`

`cwalld-enforce` writes to `/var/log/cwall/cwall.log`, so we can tail it and view updates from `cwalld/` with 
```
sudo ./cmd/cwalld/cwalld-tail
```

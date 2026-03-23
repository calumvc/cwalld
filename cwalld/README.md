
# Chinese Wall Security Daemon for SELinux

A security daemon capable of enforcing the Chinese Wall (Brewer-Nash) Model by utilising Security Enhanced Linux's mandatory access control system

## Requirements
- SELinux system (Red Hat / Fedora / ...)
- Go 1.25.7+

## Installation

### 1. Unzip cwalld.zip 
Run `unzip cwalld.zip`

### 2. Install Policies
```
cd subject_types/ ... # for all folders
sudo ./example.sh
cd object_types/ ... # for all folders
sudo make -f /usr/share/selinux/devel/Makefile ___.te
sudo semanage -i ___.pp
```

### 3. Build Executables
from `cwalld/`
```
go build ./cmd/cwalld/cwalld-tail
go build ./cmd/cwalld/cwalld-enforce

sudo cp cwalld-enforce /usr/local/sbin
sudo cp cwalld-enforce.service /usr/lib/systemd/system
sudo systemctl daemon-reload
```
For all daemons you want for examples
```
cd subjects/ ... # for all folders
gcc -o exampled exampled.c
sudo cp example /usr/local/bin
sudo cp example.service /usr/lib/systemd/system

sudo systemctl daemon-reload
```

### 4. Prepare Types
```
# change all daemon labels
sudo chcon -t bin_t /usr/local/bin/*
# fcontext and restorecon survive reboots, for cwalld-enforce
sudo semanage fcontext -a -t cwalld_exec_t /usr/local/sbin/cwalld-enforce
sudo restorecon -v /usr/local/sbin/cwalld-enforce
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
running cwalld-init is persistant and should change audit rules forever on the system

```
sudo systemctl start cwalld-enforce
```

Run any desired subject with `sudo systemctl start exampled`

`cwalld-enforce` writes to `/var/log/cwalld/cwalld.log`, so we can tail it and view updates from `cwalld/` with 
```
sudo ./cwalld-tail
```

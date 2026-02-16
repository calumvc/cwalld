
Obj1
Obj2
Obj3

Daemon1
Daemon2
Daemon3

Obj1 and Obj2 in conflict of interest class

Daemon1 can read and write to Obj1, cannot read or write to Obj2. Can also read from Obj3 but not write
Daemon2 can read and write to Obj2, cannot read or write to Obj1. Can also read from Obj3 but not write
Daemon3 can read from Obj2, and Obj3, but not write to any. Cannot read from Obj1


Daemon3 denied write to obj3
Daemon1 denied read on obj2
Daemon1 denied open on obj2
Daemon3 denied write to obj2
Daemon1 denied write to obj2
Daemon3 denied write to obj1
Daemon3 denied open on obj1
Daemon2 denied write to obj1
Daemon1 denied write to obj3
Daemon2 denied write to obj3

Daemon1 denied read on obj2
Daemon1 denied write to obj2
Daemon1 denied write to obj3

Daemon2 denied write to obj1
Daemon2 denied write to obj3

Daemon3 denied write to obj3
Daemon3 denied write to obj2
Daemon3 denied write to obj1
Daemon3 denied open on obj1


To add a label to a daemon 
Daemon in /bin
Service file wherever it needs to be
makefile used to compile to .pp
sudo semodule -i to install the policy
might have to sudo chcon -t daemon_2_exec_t ./daemon_2.sh but probs not
sudo semanage fcontext -a -t daemon_2_exec_t "usr/local/bin/daemon_2" 
./daemon_2 (edited bottom to have path to /bin/daemon_2)
sudo systemctl restart daemon_2

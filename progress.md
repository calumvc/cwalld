

1. Installed Red Hat via Virt Manager (KVM)

2. Setup basic redhat - 
$sestatus 
-> SELinux already enabled

3. Downloaded setroubleshoot & setroubleshoot-server

4. Creating foo and bar users and giving them roles
```bash
$ useradd foo
$ su -  # login as root
$# passwd -d foo # remove password for foo

$ su -
$# cd ../home
$# mkdir testgrounds # make a folder to have files to share between users
$# touch file1
$# ls -Z # unconfined_u:object_r:user_home_t:s0 file1
$# chcon -t httpd_sys_content_t file1 # change context of file1 -> label is now unconfined_u:object_r:httpd_sys_content_t:s0 
# using chcon to change it to a fake label doesnt work and SELinux alerts that there may be in issue
$# cd ..
$# chcon -R -t httpd_sys_content_t testgrounds # change directory label
# new files under this directory will have the same label
```

- Downloaded policycoreutils
- Downloaded policycoreutils-devel
- Downloaded setools-console

Following red hat docs chapter 8, made a daemon just to open and read a file in var/log/messages, ran it unconfined. Then confined it to selinux with a mydaemon_t label and it
got denied because it didnt have rights to edit things with var_log_t label. Sealert gave me steps to update its policy to work by editing the mydaemon.te file


```bash
$# semanage user -a -R "user_r" -L s0 -r s0 foo_u 
$# semanage user -a -R "user_r" -L s0 -r s0 bar_u
# created foo and bar selinux users with user_r permissions
# seinfo -u for info on selinux users

# link the selinux user to the 
$# semanage user login -a -s "foo_u" foo 
$# semanage user login -a -s "bar_u" bar

# Cannot just switch with su to login with selinux user
# Need to 
$# passwd foo # cant ssh without a password for some reason
$ ssh foo@localhost 
foo~$ ip -Z 
foo~$ foo_u:user_r:user_t:s0
# foo now has the foo selinux role
# foo is allowed in the home/testgrounds folder but unable to write anything, only read
$ ls - qdl testgrounds # results in 
drwxr-xr-x. 3 root root... # meaning that the creator of the folder (root) can rwx, users in the same group can r-x and other users can r-x. (read, write, execute), so that explains why foo could read it
# AKA CHMOD 755 ^

$ ls -dl foo 
drwx------. 4 foo foo # meaning only foo (except root) can read write or execute in this directory

# I want a directory where all users can read write and execute so I can edit the specific users, so foo can make a file that bar cannot read
$ sudo chmod 777 testgrounds

foo~$ touch my_deepest_secrets.txt
foo~$ echo "seecreeettss" >> my_deepest_secrets.txt
foo~$ ls -Z says foo_u:object_r:user_home_t:s0 ...
foo~$ ls -l deepest_darkest_secrets.txt
-rw-r--r-- # so other users can read but not write

bar~$ vim my_deepest_secrets.txt # bar can view the file
# now i need to prevent bar from reading it as in this example they are in different COI groups 

#created policy my_secret_access.te
-----------
policy_module(my_secret_access, 1.0);

require {
    type foo_u;
    type bar_u;
}

type foo_secret_t;
files_type(foo_secret_t);

allow foo_u foo_secret_t:file { create read write getattr unlink setattr };

neverallow bar_u foo_secret_t:file { create read write getattr unlink setattr };
----------------

# checkmodule would not work at all to compile the file so 
# figure out i had to use 
$# sepolicy generate --init /home/testgrounds/secret_access # new folder secret_access
$# make -f /usr/share/selinux/devel/include/Makefile my_secret_access.pp
# Policy is working!

# turns out you cant have users in definitions like that, only processes, the neverallow and allow notation is for things of _t and not _u

# Ill need to instead use categories to limit the users
chcon -t secret_file_t -l s0:c100 deepest_darkest_secrests.txt
gave it the secret_type_t label and changed category to s0:c100

# right now the label isnt doing anythig but hopefully for now bar cannot do anything with the file!

# Now I cannot ssh to foo@localhost because of the category range
# After messing further with categories it still doesnt work as intended
# Im going to scrap the category approach because it isnt using policies

```

Made a simple policy to add a secret_type_t
by default no user can read it but i can make a policy to allow all user_t (user processes) to read it
Made a simple policy to allow user_t level processes (like cat) to read it
now all users can read it

If I could change foos
foo_u:user_r:user_t:s0
to ... foo_t:s0
then it would be simple, but selinux doesnt allow you to edit the user process level


From what ive learnt so far:

I can change users foo and bar to SElinux users foo_u and bar_u which are based off of user_u permission levels
I could pretty simply make it so foo_u can read something and bar_u cant by making foo_u based off of staff_u or unconfined_u but that defeats the purpose of the chinese wall
SELinux policies dont allow you to restrict things based solely on users like the first policy i tried to make, you can only restrict things like that based on 
process types, such as user_t, but I cannot just assign foo_t and bar_t to those respective users.

Seems like I would have to find a way to fight against what SELinux is doing to deny specific users without completely changing their roles.

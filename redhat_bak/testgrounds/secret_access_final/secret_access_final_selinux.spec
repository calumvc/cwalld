# vim: sw=4:ts=4:et


%define relabel_files() \
restorecon -R /home/testgrounds/secret_access_final; \

%define selinux_policyver 42.1.7-1

Name:   secret_access_final_selinux
Version:	1.0
Release:	1%{?dist}
Summary:	SELinux policy module for secret_access_final

Group:	System Environment/Base
License:	GPLv2+
# This is an example. You will need to change it.
# For a complete guide on packaging your policy
# see https://fedoraproject.org/wiki/SELinux/IndependentPolicy
URL:		http://HOSTNAME
Source0:	secret_access_final.pp
Source1:	secret_access_final.if
Source2:	secret_access_final_selinux.8


Requires: policycoreutils-python-utils, libselinux-utils
Requires(post): selinux-policy-base >= %{selinux_policyver}, policycoreutils-python-utils
Requires(postun): policycoreutils-python-utils
BuildArch: noarch

%description
This package installs and sets up the  SELinux policy security module for secret_access_final.

%install
install -d %{buildroot}%{_datadir}/selinux/packages
install -m 644 %{SOURCE0} %{buildroot}%{_datadir}/selinux/packages
install -d %{buildroot}%{_datadir}/selinux/devel/include/contrib
install -m 644 %{SOURCE1} %{buildroot}%{_datadir}/selinux/devel/include/contrib/
install -d %{buildroot}%{_mandir}/man8/
install -m 644 %{SOURCE2} %{buildroot}%{_mandir}/man8/secret_access_final_selinux.8
install -d %{buildroot}/etc/selinux/targeted/contexts/users/


%post
semodule -n -i %{_datadir}/selinux/packages/secret_access_final.pp

if [ $1 -eq 1 ]; then

fi
if /usr/sbin/selinuxenabled ; then
    /usr/sbin/load_policy
    %relabel_files
fi;
exit 0

%postun
if [ $1 -eq 0 ]; then

    semodule -n -r secret_access_final
    if /usr/sbin/selinuxenabled ; then
       /usr/sbin/load_policy
       %relabel_files
    fi;
fi;
exit 0

%files
%attr(0600,root,root) %{_datadir}/selinux/packages/secret_access_final.pp
%{_datadir}/selinux/devel/include/contrib/secret_access_final.if
%{_mandir}/man8/secret_access_final_selinux.8.*


%changelog
* Sat Nov 29 2025 YOUR NAME <YOUR@EMAILADDRESS> 1.0-1
- Initial version


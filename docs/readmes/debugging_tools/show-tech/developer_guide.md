## Context:


This feature enables an operator to capture the essential state of the gateway (AGW/FEG) package it in intention to upload it to some user configured location.

* Feature should be developed with extensibility point of view. We should be able to easily add new state which needs to be snapshotted. We shouldn’t have to rebuild or recompile in case we want to add additional state which needs to be collected. It should be straightforward for operators to modify this.
* The tool should be reliable. It shouldn’t crash. We should be able to recover from any failure and still gather as much information as possible.


## How to use it?

Show-tech feature is built based on ansible playbook which is located on magma repo.
We currently provide 2 playbooks:
1) install_prereq.yml - to be run once on host to upgrade ansible version to > 2.9.0
2) show-tech.yml - the report status feature playbook - it responsible to package all files collected and all commend processed to single .tar.gz file.

instructions: https://github.com/magma/magma/blob/master/show-tech/README.md


## Developer guide:


- Each playbook is including all the roles and tasks in ordered way as they are mentioned.
- Constants file is **constants.yml** and it includes all files we collect and all commands we run and gather their outputs.

```
# PREREQ section
ansible_apt_repo: "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main"
ubuntu_apt_keyserver: "keyserver.ubuntu.com"
ubuntu_apt_id: "93C4A3FD7BB9C367"
apt_prereq_packages:
  - git
  - ansible

# SHOW_TECH section
general_commands:
  - "uname -a"
  - "top -b -n 1"
  - "df -kh"
debian_commands:
  - "dpkg -l | grep magma"
  - "ovs-vsctl show"
  - "apt show magma"
  - "service magma@* status"
  - "show_gateway_info.py"
  - "checkin_cli.py"
  - "mobility_cli.py get_subscriber_table"
  - "pipelined_cli.py debug display_flows"
  - "enodebd_cli.py get_all_status"
  - "ip addr"
  - "ping google.com -I eth0 -c 5"
  - "journalctl -u magma@*"
  - "timeout 60 sudo tcpdump -i any sctp -w {{report_dir_output.stdout}}/sctp.pcap"


# files to collect from src to relative destination in tar.gz package
paths_to_collect:
  - "/var/opt/magma/configs/*"
  - "/etc/magma/configs/*"
  - "/etc/systemd/system/*"
  - "/tmp/*core*"
  - "/usr/local/bin/mme"
  - "/var/log/syslog"
  - "/var/log/MME.*"
  - "/var/log/enode.*"

```

The role that reads all constants is load_vars
as can seen in show-tech playbook:

```
- name: Collect GW data
  hosts: localhost
  pre_tasks:
    - name: Verify Ansible meets version requirements.
      assert:
        that: "ansible_version.full is version_compare('2.9', '>=')"
  become: no
  gather_facts: yes
  roles:
    - role: load_vars
    - role: install_ansible_modules
    - role: destdir
    - role: files
    - role: commands
    - role: package

```

Or pre-req playbook example:

```
- name: install pre-requisites for show-tech
  hosts: localhost
  become: no
  gather_facts: yes
  roles:
    - role: load_vars
    - role: upgrade_ansible

```

### How does role look like?

Role includes tasks which can be builtin ansible modules or community ones.
In this example we will look into the role that collect files as you set them in variable called “*paths_to_collect*” and uses builtin modules to work with files. (https://docs.ansible.com/ansible/2.8/modules/list_of_files_modules.html)

*NOTE*: *report_dir_output.stdout* is the variable containing the  dynamic report dir to be used across all roles and is being created on role *destdir.*

```
# Magma files collection based on map paths_to_collect
- name: Create a dest directory for magma files under report directory
  file:
    path: "{{report_dir_output.stdout}}/{{item | dirname}}"
    state: directory
    mode: "0700"
  with_fileglob: "{{paths_to_collect}}"



- name: Copy magma files to destination
  copy:
    src: "{{item}}"
    dest: "{{report_dir_output.stdout }}/{{item | dirname}}"
  with_fileglob: "{{paths_to_collect}}"

```
 ### How do i test my changes:


On agw host as user root perform the following after made some git PR.

Clone your PR to host and upgrade ansible:


```
$ git clone https://github.com/<your_github_user>/magma.git -b <your_branch>
$ cd magma/show-tech/
$ ansible-playbook install_prereq.yml


[WARNING]: provided hosts list is empty, only localhost is available


PLAY [install pre-requisites for show-tech] ************************************

TASK [setup] *******************************************************************
ok: [localhost]

TASK [load_vars : include the vars for deployment from constants.yml] **********
ok: [localhost]

TASK [upgrade_ansible : Validate OS is debian] *********************************
ok: [localhost] => {
    "changed": false,
    "msg": "All assertions passed"
}

TASK [upgrade_ansible : Updating repos with ansible for debian repo] ***********
changed: [localhost]

TASK [upgrade_ansible : Import the ubuntu apt key] *****************************
ok: [localhost]

TASK [upgrade_ansible : install packages] **************************************
changed: [localhost] => (item=[u'git', u'ansible'])

PLAY RECAP *********************************************************************
localhost                  : ok=6    changed=2    unreachable=0    failed=0

```

Run the show-tech playbook:

```

# ansible-playbook show-tech.yml
[WARNING]: provided hosts list is empty, only localhost is available. Note that the implicit localhost does not match 'all'

PLAY [Collect GW data] ********************************************************************************************************************************************************************************************************************************************************************

TASK [Gathering Facts] ********************************************************************************************************************************************************************************************************************************************************************
ok: [localhost]

TASK [Verify Ansible meets version requirements.] *****************************************************************************************************************************************************************************************************************************************
ok: [localhost] => {
    "changed": false,
    "msg": "All assertions passed"
}

TASK [load_vars : include the vars for deployment from constants.yml] *********************************************************************************************************************************************************************************************************************
ok: [localhost]

TASK [destdir : Get timestamp] ************************************************************************************************************************************************************************************************************************************************************
changed: [localhost]

TASK [destdir : Get report destination directory] *****************************************************************************************************************************************************************************************************************************************
changed: [localhost]

TASK [destdir : debug] ********************************************************************************************************************************************************************************************************************************************************************
ok: [localhost] => {
    "report_dir_output.stdout_lines": [
        "/tmp/magma_reports/report.magma.2020-11-29T14:34:49Z"
    ]
}

TASK [destdir : Create a directory for report] ********************************************************************************************************************************************************************************************************************************************
changed: [localhost]

TASK [Create a dest directory for magma files under report directory] *********************************************************************************************************************************************************************************************************************

TASK [Copy magma files to destination] ****************************************************************************************************************************************************************************************************************************************************
[WARNING]: Unable to find '/etc/magma/configs' in expected paths (use -vvvvv to see paths)
changed: [localhost] => (item=/var/opt/magma/configs/gateway.mconfig)
changed: [localhost] => (item=/etc/systemd/system/magma@pipelined.service.dpkg-old)
changed: [localhost] => (item=/etc/systemd/system/magma@mme.service)
changed: [localhost] => (item=/etc/systemd/system/magma@lighttpd.service.dpkg-old)
changed: [localhost] => (item=/etc/systemd/system/magma@redirectd.service.dpkg-old)
[...]
changed: [localhost] => (item=/usr/local/bin/mme)
changed: [localhost] => (item=/var/log/syslog)
changed: [localhost] => (item=/var/log/MME.magma.root.log.INFO.20201027-190839.7640)
changed: [localhost] => (item=/var/log/MME.INFO)


TASK [Get output of commands for all distributions] ***************************************************************************************************************************************************************************************************************************************
changed: [localhost] => (item=uname -a)
changed: [localhost] => (item=top -b -n 1)
changed: [localhost] => (item=df -kh)

TASK [Run debian commands] ****************************************************************************************************************************************************************************************************************************************************************
included: /root/magma/show-tech/roles/commands/tasks/debian.yml for localhost
TASK [Get output of debian commands] ******************************************************************************************************************************************************************************************************************************************************
changed: [localhost] => (item=dpkg -l | grep magma)
changed: [localhost] => (item=ovs-vsctl show)
changed: [localhost] => (item=apt show magma)
changed: [localhost] => (item=service magma@* status)
changed: [localhost] => (item=show_gateway_info.py)
changed: [localhost] => (item=checkin_cli.py)
changed: [localhost] => (item=mobility_cli.py get_subscriber_table)
changed: [localhost] => (item=pipelined_cli.py debug display_flows)
changed: [localhost] => (item=enodebd_cli.py get_all_status)
changed: [localhost] => (item=ip addr)
changed: [localhost] => (item=ping google.com -I eth0 -c 5)
changed: [localhost] => (item=journalctl -u magma@*)
changed: [localhost] => (item=timeout 60 sudo tcpdump -i any sctp -w /tmp/magma_reports/report.magma.2020-11-29T14:34:49Z/sctp.pcap)

TASK [package : Compress directory report_dir into /tmp/magma_reports/report.<HOSTNAME>.<DATE>.tgz] ***************************************************************************************************************************************************************************************
changed: [localhost]

TASK [package : Recursively remove report directory /tmp/magma_reports/report.<HOSTNAME>.<DATE>] ******************************************************************************************************************************************************************************************
changed: [localhost]

PLAY RECAP ********************************************************************************************************************************************************************************************************************************************************************************
localhost                  : ok=16   changed=10   unreachable=0    failed=0    skipped=0    rescued=0    ignored=0

```

Validate the content of the package.

commands_output.log - will include all commands and their stdout and stderr events.
sctp.pcap - will include pcap file.

```
# tar tvf /tmp/magma_reports/report.magma.2020-11-29T14\:21\:42Z.tgz
drwx------ root/root 0 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/usr/
drwx------ root/root 0 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/var/
drwx------ root/root 0 2020-11-29 14:21 report.magma.2020-11-29T14:21:42Z/etc/
-rw-r--r-- root/root 15996134 2020-11-29 14:23 report.magma.2020-11-29T14:21:42Z/commands_output.log
-rw-r--r-- root/root 2336 2020-11-29 14:24 report.magma.2020-11-29T14:21:42Z/sctp.pcap
drwx------ root/root 0 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/usr/local/
drwx------ root/root 0 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/usr/local/bin/
-rw-r--r-- root/root 69056616 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/usr/local/bin/mme
drwx------ root/root 0 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/var/log/
drwx------ root/root 0 2020-11-29 14:21 report.magma.2020-11-29T14:21:42Z/var/opt/
-rw-r--r-- root/root 948 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/var/log/MME.magma.root.log.INFO.20201027-190839.7640
-rw-r--r-- root/root 948 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/var/log/MME.magma.root.log.INFO.20201027-190905.8105
-rw-r--r-- root/root 948 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/var/log/MME.magma.root.log.INFO.20201027-190813.7191
-rw-r--r-- root/root 12471697 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/var/log/MME.INFO
-rw-r--r-- root/root 50620666 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/var/log/syslog
drwx------ root/root 0 2020-11-29 14:21 report.magma.2020-11-29T14:21:42Z/var/opt/magma/
drwx------ root/root 0 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/var/opt/magma/configs/
-rw-r--r-- root/root 2795 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/var/opt/magma/configs/gateway.mconfig
drwx------ root/root 0 2020-11-29 14:21 report.magma.2020-11-29T14:21:42Z/etc/systemd/
drwx------ root/root 0 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/etc/systemd/system/
-rw-r--r-- root/root 363 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/etc/systemd/system/magma@pipelined.service.dpkg-old
-rw-r--r-- root/root 1258 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/etc/systemd/system/magma@mme.service
-rw-r--r-- root/root 446 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/etc/systemd/system/magma@lighttpd.service.dpkg-old
-rw-r--r-- root/root 313 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/etc/systemd/system/magma@redirectd.service.dpkg-old
-rw-r--r-- root/root 470 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/etc/systemd/system/magma@redis.service.dpkg-old
-rw-r--r-- root/root 1010 2020-11-29 14:22 report.magma.2020-11-29T14:21:42Z/etc/systemd/system/magma@lighttpd.service

```


Check the commands output:
```

$ grep "Command:" -A 3 report.magma.2020-11-29T14:21:42Z/commands_output.log

Command: uname -a
Linux magma 4.9.0-8-amd64 #1 SMP Debian 4.9.110-3+deb9u6 (2018-10-08) x86_64 GNU/Linux
Command: top -b -n 1
top - 14:22:59 up 108 days, 21:53, 1 user, load average: 1.31, 0.73, 0.46
Tasks: 133 total, 2 running, 130 sleeping, 0 stopped, 1 zombie
%Cpu(s): 3.8 us, 0.6 sy, 0.0 ni, 92.9 id, 0.0 wa, 0.0 hi, 2.6 si, 0.0 st
--
Command: df -kh
Filesystem Size Used Avail Use% Mounted on
udev 3.9G 0 3.9G 0% /dev
tmpfs 787M 81M 707M 11% /run
--
Command: ovs-vsctl show
de6cf0e5-454c-4cef-9a65-d93ba083ab8e
Bridge "gtp_br0"
Controller "tcp:127.0.0.1:6654"
--
Command: apt show magma
Package: magma
Version: 1.3.1-1605904866-7933c828
Priority: extra
--
Command: service magma@* status
● magma@td-agent-bit.service - TD Agent Bit
Loaded: loaded (/etc/systemd/system/magma@td-agent-bit.service; disabled; vendor preset: enabled)
Active: active (running) since Thu 2020-11-26 11:03:49 UTC; 3 days ago
--
Command: show_gateway_info.py
Hardware ID:
e677ec40-19d3-461e-bc75-b526432a0f33
--
Command: checkin_cli.py
1. -- Testing TCP connection to controller.magma.etagecom.io:443 (http://controller.magma.etagecom.io:443/) --
2. -- Testing Certificate --
3. -- Testing SSL --
--
Command: mobility_cli.py get_subscriber_table
SID IP APN
IMSI311980000039039 192.168.128.150 oai.ipv4
IMSI311980000039058 192.168.128.216 oai.ipv4
--
Command: pipelined_cli.py debug display_flows
cookie=0x0, duration=361944.238s, table=mme(main_table), n_packets=12686776, n_bytes=1255948751, idle_age=2, hard_age=65534, priority=10,tun_id=0x51,in_port=1 actions=mod_dl_src:02:00:00:00:00:01,mod_dl_dst:ff:ff:ff:ff:ff:ff,load:0x8ddf408a29279→OXM_OF_METADATA[],resubmit(,ingress(main_table))
cookie=0x0, duration=361944.063s, table=mme(main_table), n_packets=9554644, n_bytes=1226610964, idle_age=45, hard_age=65534, priority=10,tun_id=0x39,in_port=1 actions=mod_dl_src:02:00:00:00:00:01,mod_dl_dst:ff:ff:ff:ff:ff:ff,load:0x8ddf408a292b1→OXM_OF_METADATA[],resubmit(,ingress(main_table))
cookie=0x0, duration=361943.145s, table=mme(main_table), n_packets=1283, n_bytes=97588, idle_age=53, hard_age=65534, priority=10,tun_id=0x52,in_port=1 actions=mod_dl_src:02:00:00:00:00:01,mod_dl_dst:ff:ff:ff:ff:ff:ff,load:0x8ddf408a283f9→OXM_OF_METADATA[],resubmit(,ingress(main_table))
--
Command: enodebd_cli.py get_all_status
--- eNodeB Serial: 120200004917CNJ0028 ---
IP Address..................10.0.2.243
eNodeB Connected via TR-069............ON
--
Command: ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1
link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
inet 127.0.0.1/8 scope host lo
--

Command: ping google.com -I eth0 -c 5
PING google.com (172.217.11.238) from 192.168.149.2 eth0: 56(84) bytes of data.
64 bytes from den02s01-in-f14.1e100.net (172.217.11.238): icmp_seq=1 ttl=117 time=10.4 ms
64 bytes from den02s01-in-f14.1e100.net (172.217.11.238): icmp_seq=2 ttl=117 time=10.2 ms
--
Command: journalctl -u magma@*
-- Logs begin at Sun 2020-11-29 11:16:52 UTC, end at Sun 2020-11-29 14:23:19 UTC. --
Nov 29 11:16:52 magma sessiond[11821]: I1129 11:16:52.015013 11821 LocalEnforcer.cpp:255] IMSI311980000039313-78566 used 40 tx bytes and 52 rx bytes for rule allowlist_sid-IMSI311980000039313-oai.ipv4
Nov 29 11:16:52 magma sessiond[11821]: I1129 11:16:52.015242 11821 LocalEnforcer.cpp:255] IMSI311980000021494-624955 used 62 tx bytes and 78 rx bytes for rule allowlist_sid-IMSI311980000021494-oai.ipv4
--
Command: timeout 60 sudo tcpdump -i any sctp -w /tmp/magma_reports/report.magma.2020-11-29T14:21:42Z/sctp.pcap
```



## References:

https://docs.ansible.com/ansible/2.8/modules/modules_by_category.html - all ansible modules by category.
https://docs.ansible.com/ansible/2.8/user_guide/playbooks.html - ansible user guide on playbooks.

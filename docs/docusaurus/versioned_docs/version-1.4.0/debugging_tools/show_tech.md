---
id: version-1.4.0-show_tech
title: Show Tech
hide_title: true
original_id: show_tech
---
# Show-Tech

## Overview:

This feature enables an operator to capture the essential state of the gateway (Currently only supported on AGW) and dumps it in a destination directory. The collected data can then be shared with support teams to help identify and resolve issues quickly.


### Usage

```
# On your GW host, run the following command as user root:
# If you have git repo checked out already
$ cd ${MAGMA_ROOT}/show-tech
$ ansible-playbook show-tech.yml

# In case you want to download and process latest version of this playbook from Magma's master:
$ ansible-pull -U https://github.com/magma/magma.git show-tech/show-tech.yml -d /tmp/show-tech --purge
```

The captured output is dumped in `/tmp/magma_reports/report.magma.<date>.tgz`

### What is collected upon running show-tech ?

Currently, show-tech tool collects the following information and packages them into a destination directory.
*Disclaimer : Note this tool collects information necessary for debugging issues on the AGW. These could contain customer sensitive data.*

**Outputs from running various commands**

* System information  (kernel-name, node-name, kernel-version, machine, processor, hardware platform, OS)
* Top command output
* Disk space usage
* Magma packages installed
* Ovs switch database content
* Status of magma services
* Gateway information (hardwareID and challenge key of the gateway)
* Orc8r check in information
* Subscriber table on the gateway
* Current set of flows managed by pipelined
* IP address of interfaces on the gateway
* Logs of all magma services
* Sctp packet capture for 60 seconds
* MME core stack traces

**Following Files on the gateway**

* Gateway configs: mconfig and overriding yaml files
* Syslog, MME, enodebd service logs
* MME binary
* Core files
* Systemd files

## Developer guide

### Design principles

* **Extensibility**: Easy to add new state collectors without the need for rebuilding software. This allows operators to extend the functionality at runtime.
* **Reliability**: The tool is used to debug issues and needs to reliably collect all needed state. In the event of any failure, tool should still gather as much information as possible.

The show-tech tool is built using ansible playbooks that are located on the magma github repo. We currently provide the following playbooks:

1. install_prereq.yml - to be run once on host to upgrade ansible version to > 2.9.0
2. show-tech.yml - the state collection playbook. It is responsible for collecting all state and packaging the all files into one .tar.gz file.
3. constants.yml - is the constants file and includes all files we collect, as well as all commands we run and gather their outputs. The role that reads all constants is `load_vars` as can be seen in show-tech playbook.

### How to add a new file or set of files to be included in the destination ?

This can be added to the `paths_to_collect` in `constants.yml` file in “show-tech” directory.
For example, if we want to include core files to the list of files we capture using `show-tech` tool, add the line as indicated below to the `constants.yml` file.

```
# files to collect from src to relative destination in tar.gz package
paths_to_collect:
  - "/var/opt/magma/configs/*"
  - "/etc/magma/configs/*"
  - "/etc/systemd/system/*"
  - "/usr/local/bin/mme"
  - "/var/log/syslog"
  - "/var/log/MME.*"
  - "/var/log/enode.*"
 **** **- "/var/core/*"** **** << Add the new file/path here
```



### How to add a new command output to be included in the destination ?

This can be added to the `debian_commands` in `constants.yml` file in “show-tech” directory.
For example, if we want to add current running processes to the output we capture using `show tech` tool, add the command as indicated below to the `constants.yml` file

```
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
  - "timeout 60 sudo tcpdump -i any port 48080 -w {{report_dir_output.stdout}}/any-48080.pcap"
  - "ps aux" <<< Add the new command here
```

### How to add a new role ?

Roles include tasks that can be built-in ansible modules or community ones. In this example we will look into the role that collects files already set in the variable called “*paths_to_collect*”. This uses a built-in module to work with the specified files. (https://docs.ansible.com/ansible/2.8/modules/list_of_files_modules.html)
*NOTE*: *report_dir_output.stdout* is the variable containing the dynamic report dir to be used across all roles and is being created on role *destdir.*

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

### How to test changes:

On the AGW host, clone your PR and upgrade ansible. Run the following as root

```
$ git clone https://github.com/<your_github_user>/magma.git -b <your_branch>
$ cd magma/show-tech/
$ ansible-playbook install_prereq.yml
...
PLAY [install pre-requisites for show-tech] ************************************
...
```

Run the show-tech playbook:

```
$ ansible-playbook show-tech.yml
```

Validate the content of the package.
commands_output.log includes all commands and their stdout/stderr

```
$ tar tvf /tmp/magma_reports/report.magma.2020-11-29T14\:21\:42Z.tgz
```

```
$ grep "Command:" -A 3 report.magma.2020-11-29T14:21:42Z/commands_output.log
```

## References:

https://docs.ansible.com/ansible/2.8/modules/modules_by_category.html - all ansible modules by category. https://docs.ansible.com/ansible/2.8/user_guide/playbooks.html - ansible user guide on playbooks.

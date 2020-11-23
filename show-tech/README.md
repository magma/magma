## Prereq:

On your GW host make sure you have the following packages:

```
$ echo "deb http://ppa.launchpad.net/ansible/ansible/ubuntu trusty main" >> /etc/apt/sources.list
$ apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 93C4A3FD7BB9C367
$ sudo apt update
$ apt-get install git ansible software-properties-common dirmngr -y
$ ansible-galaxy collection install community.general
```

## Execution:

```
# On your GW host, run the following command as user root:

$ ansible-playbook show-tech.yml

# In case you want to download and process latest version of this playbook from Magma's master:

$ ansible-pull -U https://github.com/magma/magma.git show-tech/show-tech.yml -d /tmp/show-tech --purge
```

## Output:

Report will be saved under /tmp/magma_reports/report.<HOSTNAME>.<DATE>.tgz

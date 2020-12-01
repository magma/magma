## Prereq:

On your GW host make sure you have the following packages:
```
$ apt-get install git ansible -y
```

and cloned magma repo:
```
$ git clone https://github.com/magma/magma.git
```


Then run the pre-req playbook which will upgrade ansible, need to run one time.

```
$ cd magma/show-tech
$ ansible-playbook install_prereq.yml
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

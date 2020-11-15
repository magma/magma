## Prereq:

On your GW host make sure yo have the following:
1) ansible binary
2) ansible collection 'community.general' should be installed by:
```
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

# Ansible Upgrade

This set of roles and playbook should be use and expanded to automate SUT operations

Supported operations:

- upgrade SUT with latest sw version

## How to

use ansible-runner from command line <ens_magma> location

`ansible-runner -p upgrade.yaml --limit <limit host to upgrade> run Magma_Ansible`

this module is already use in python automation as well.

For more info please [read](https://ansible-runner.readthedocs.io/en/stable/intro.html)

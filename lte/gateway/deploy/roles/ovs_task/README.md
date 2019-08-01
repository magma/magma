## Cleanup ovs build with Ansible.

Goal: Build Openvswitch with custom patches allowing gtp support

In `roles/ovs_gtp`, you can find the main task in order to build the associate patches.

#### Requirement
---

```
ansible <= 2
python
```

#### How to launch a build
---

First of all fill the hosts file, you can either do

- Local build

```[local]
127.0.0.1   ansible_connection=local
```
- Vagrant

```[localvagrant]
127.0.0.1   ansible_ssh_user=vagrant    ansible_ssh_port=2222   ansible_ssh_private_key_file=~/.vagrant.d/insecure_private_key
```
- Using a server/vm

```
[my_pool]
Ip_address ansible_port=22 ansible_user=ubuntu
```

When your hostfile is completed, complete the playbook with :

##### hosts
- `hosts: localvagrant` String      The host  or pool of hosts you want to build on.

##### vars
- `ovs_version`			    String      The ovs version (github version or tag)
- `ovs_version_short`   String      The ovs version short
- `WORK_DIR`            String      The Working Directory
- `patches`             AnsibleList All patches you want to apply

To run the playbook just do:

```ansible-playbook ovs_gtp.yml -i hosts```

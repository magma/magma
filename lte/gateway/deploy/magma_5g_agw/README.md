# 5G AGW installation

## Prerequisites

To setup a Magma Access Gateway, you will need a machine that satisfies the following requirements:

**AGW\_HOST**: 64bit-X86 machine, bare metal or virtual machine installed Ubuntu Server 20.04 LTS. You will need two ethernet ports. One port is for the SGi interface (default: enp1s0) and one for the S1 interface (default: enp2s0). Note that the Ansible scripts will rename the `enp1s0` interface to `eth0` and the enp2s0 to eth1. Below specified hardware specification is the minimum for AGW installation. 

Physical: 

| Description |  Resources |
| -------- | ----------- |
| Platform | Bare Metal |
|CPU | Intel(R) i3 7100u |
| NIC | 2X 1G ethernet interfaces |
| RAM | Minimum 4GB RAM preferably 8GB DDR3 1066Mhz |
| Storage | Minimum 32GB SSD preferably 128GB-256GB SSD |
| Operating System | Ubuntu Focal 20.04 (LTS) |

Virtual:

| Description |  Resources |
| -------- | ----------- |
| Platform | Private Cloud/DC |
|CPU | 4X vCPU / AMD64 dual-core processor around 2GHz clock speed or faster  |
| NIC | 2X 1G ethernet interfaces |
| RAM | Minimum 4GB RAM preferably 8GB DDR3 1066Mhz |
| Storage | Minimum 32GB SSD preferably 128GB-256GB SSD |
| Operating System | Ubuntu Focal 20.04 (LTS) |


To deploy AGW on the above machine you need a workstation. 

**Workstation**: It could be a VM or bare metal which should have Ansible installed and key based SSH authentication to the AWG host.

## Deployments

To deploy AGW on the server first we need to clone the Ansible repo on our workstation and then create a host file for the Ansible with the information of the two IP address of your interfaces. After cloning the Ansible repo you also have to modify the ansible.cfg file to modify the private key path; using which Ansible is going to SSH to your box.

### Clone the Ansible repo and go inside the Ansible path

```
git clone https://github.com/magma/magma.git
cd magma/lte/gateway/deploy/magma_5g_agw
```

### Modify the inventory file for Ansible:

```
# vi agw_ansible_hosts
---
hosts:
  all:
    - < IP of eth0 >
```

### Change the private Key file path and user for SSH
```
# vi ansible.cfg
Change line # 7 to correct the ssh user if it is not ubuntu 
Change line # 8 to add the correct path to your private key 
```
### Overwrite the attirbute values
If you want to overwrite any default attribute value then modify the file `ansible_vars.yaml`
```
vi ansible_vars.yaml
```

### Now we are all set to do the deployment by executing the below command:

```
ansible-playbook agw_deploy.yaml -e "@ansible_vars.yaml"
```

Once the execution completed successfully you can make sure magma is running by executing the below command on your AGW box:

```
/etc/update-motd.d/99-magma
```
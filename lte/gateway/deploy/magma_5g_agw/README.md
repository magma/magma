# AGW installation

## Prerequisites

To setup a Magma Access Gateway, user must need a machine that satisfies the following requirements:

**AGW\_HOST**: 
    -64bit-X86 machine
    -Bare metal or virtual machine installed with Ubuntu Server 20.04 LTS
    -Two ethernet ports (One port is for the SGi interface (default: enp1s0) and the  other port is for the S1 interface (default: enp2s0)

**Note**: The Ansible scripts renames the `enp1s0` interface to `eth0` and the enp2s0 to eth1. 

Below are the minimum hardware specifications that required for AGW installation. 

**Physical**: 

| Description |  Resources |
| -------- | ----------- |
| Platform | Bare Metal |
|CPU | Intel(R) i3 7100u |
| NIC | 2X 1G ethernet interfaces |
| RAM | Minimum 4GB RAM preferably 8GB DDR3 1066Mhz |
| Storage | Minimum 32GB SSD preferably 128GB-256GB SSD |
| Operating System | Ubuntu Focal 20.04 (LTS) |

**Virtual**:

| Description |  Resources |
| -------- | ----------- |
| Platform | Private Cloud/DC |
|CPU | 4X vCPU / AMD64 dual-core processor around 2GHz clock speed or faster  |
| NIC | 2X 1G ethernet interfaces |
| RAM | Minimum 4GB RAM preferably 8GB DDR3 1066Mhz |
| Storage | Minimum 32GB SSD preferably 128GB-256GB SSD |
| Operating System | Ubuntu Focal 20.04 (LTS) |


To deploy AGW on the above machine a workstation is required. 

**Workstation**: It could be a VM or bare metal which must have Ansible installed in it and key based SSH authentication to the AGW host.

## Deployments

To deploy AGW on the server, first user must clone the Ansible repo on their workstation and create a host file for it with two IP address of their interfaces. After cloning the Ansible repo, user must modify the private key path in the **ansible.cfg file**, by which the Ansible goes to SSH of their AGW box.

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
Change line # 8 to correct the ssh user if it is not ubuntu 
Change line # 9 to add the correct path to your private key 
```
### Overwrite the attirbute values
To overwrite any default attribute value, modify the `ansible_vars.yaml` file.
```
vi ansible_vars.yaml
```
In the `ansible_vars.yaml` file, if you want to use the static IP, then make **configure_dhcp = false** and update the values for the attributes listed below.
```
eth0: ""
eth1: ""
eth0_netmask: ""
eth1_netmask: ""
eth0_gateway: ""
eth1_gateway: ""
```
**Note:** If configure_dhcp = **true**, then it is not required to update the values for the attributes listed above.
#### To install 4G version of the AGW 
If you want to install 1.6.0 or 1.6.1 version of AGW, then modify the below attribute in the `ansible_vars.yaml` file.

```
magma_5g_upgrade = false
# If the AGW version is 1.6.1, then change the below attribute value to focal-1.6.1 and if the version is 1.6.0, then it should be focal-1.6.0.
magma_pkgrepo_dist: "focal-1.6.0"
```

#### To install 5G version of the AGW
If you want to install 1.7.0 version of AGW, then modify the below attribute in the `ansible_vars.yaml` file.

```
magma_5g_upgrade = true
# Please specify the version number in below attribute, to install that particular version (1.7.0)
magma5gVersion = "1.7.0"
```

#### To configure the AGW 
If you want to configure your AGW, then modify the below attribute in the `ansible_vars.yaml` file.

```
magma_control_proxy_path: ""
magma_rootCA_path: ""
```

### All set to do the deployment by executing the below command:
```
ansible-playbook agw_deploy.yaml -e "@ansible_vars.yaml"
```

Once the execution is completed successfully, user can make sure that magma is running by executing the below command on their AGW box.

```
/etc/update-motd.d/99-magma
```
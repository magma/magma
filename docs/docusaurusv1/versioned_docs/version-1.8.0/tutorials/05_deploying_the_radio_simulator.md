---
id: version-1.8.0-05_deploying_the_radio_simulator
title: 5. Deploying the radio simulator
hide_title: true
original_id: 05_deploying_the_radio_simulator
---

# 5. Deploying the radio simulator

## Create an instance on AWS

Create an AWS EC2 instance running Ubuntu 20.04:

```console
aws ec2 run-instances \
  --security-group-ids <your security group> \
  --image-id ami-0568936c8d2b91c4e \
  --count 1 \
  --instance-type t2.xlarge \
  --key-name <your ssh key name> \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=srsran}]' \
  --block-device-mapping "[ { \"DeviceName\": \"/dev/sda1\", \"Ebs\": { \"VolumeSize\": 50 } } ]"
```

Replace the security group ID with one that allows SSH access and note the instance ID.

## Attach a secondary network interface to the instance

Using `SubnetId` of the **S1** subnet that was created during step 1, create a new network interface:

```console
aws ec2 create-network-interface --subnet-id <your subnet ID> --group <your security group>
```

Attach the network interface to the EC2 instance:

```console
aws ec2 attach-network-interface --network-interface-id <your network interface ID> --instance-id <your instance ID> --device-index 1
```

## Add the machine to Juju

Wait for the instance to boot up and be accessible via SSH, then add it as a Juju machine:

```console
juju add-machine --private-key=<path to your private key> ssh:ubuntu@<EC2 instance IP address>
```

## Configure Netplan to use the secondary network interface

SSH into the machine:

```console
juju ssh <Your instance ID>
```

Retrieve the mac address used by `eth1`:

```console
ip a show eth1
```

Create a file named `99-srsran.yaml` that contains the following content and move it over
to `/etc/netplan/`:

```yaml title="99-srsran.yaml"
network:
  version: 2
  ethernets:
    eth1:
      dhcp4: true
      dhcp4-overrides:
        use-routes: false
      dhcp6: false
      match:
        macaddress: <eth1 interface mac address>
      set-name: eth1
```

Apply the netplan configuration:

```console
netplan apply
```

## Deploy the srsRAN radio simulator

Deploy srsRAN to the machine:

```console
juju deploy srs-enb-ue --channel=edge --config bind-interface="eth1" --to <Machine ID>
```

## Integrate the radio simulator with Magma Access Gateway

```console
juju relate srs-enb-ue:lte-core magma-access-gateway-operator:lte-core
```

---
id: deploy_config_agw_ha
title: Configure AGW for HA
hide_title: true
---

# Configure Access Gateway for High-Availability

This document outlines the necessary steps to deploy and configure a
Magma access gateway on AWS. This document also outlines configuring the AWS
gateway to serve as a secondary to a primary gateway running at an edge site.

## Deployment

### Build AGW AMI

Steps:

1. Download packer onto your host machine at https://www.packer.io/downloads.html
2. Run the following

```
[~] cd magma/orc8r/tools/packer
[~/magma/orc8r/tools/packer] packer build -force \
    -var "aws_access_key=YOUR_ACCESS_KEY" \
    -var "aws_secret_key=YOUR_SECRET_KEY" \
    -var "subnet=YOUR_SUBNET" \
    -var "vpc=YOUR_VPC" \
    debian-stretch-aws.json
```

YOUR_SUBNET and YOUR_VPC should specify an existing subnet and vpc on your AWS
region. The choice of subnet and vpc won't affect the final box. These are the
subnet/vpc which the box is launched into while building.

The result should show

```
==> Builds finished. The artifacts of successful builds are:
--> amazon-ebs: AMIs were created:
us-west-1: ami-0f1c9db5a767a0296
```

### Deploy AGW AMI

On AWS:

1. Navigate to the EC2 Service
2. Select `Launch Instance`
3. Select the AMI that was built in the previous step. This AMI will exist
under `My AMIs` section.
4. On page `Choose an Instance Type`, select a c4.xlarge instance type. Proceed
to `Configure Instance Details`.
5. On page `Configure Instance Details`, use the default settings. Proceed to
`Add Storage`.
6. On page `Add Storage`, use default of 8gb. Proceed to `Add Tags`.
7. On page `Add Tags`, optionally add tags (e.g. `Magma Secondary Gateway`)
to identify this as a secondary.
Magma AGW. Proceed to `Configure Security Group`.
8. On page “Configure Security Group”, create a new security group with the
rules listed below. It is advised to limit the source IPs to the subnet that i
the primary gateway resides in for all rules other than SSH. Proceed to
`Review and Launch`.

|Type	|Protocol	|Port Range	|Source	|Description	|
|---	|---	|---	|---	|---	|
|SSH	|TCP	|22	|0.0.0.0/0	|-	|
|SCTP (132)	|SCTP (132)	|All	|0.0.0.0/0	|-	|
|Custom TCP	|TCP	|3386	|0.0.0.0/0	|-	|
|All UDP	|UDP	|0 - 65535	|0.0.0.0/0	|	|
|All ICMP - IPv4	|ICMP	|All	|0.0.0.0/0	|-	|

1. Review that the selected settings are as described here. Then proceed to
`Launch`.
2. Select `Create a new key pair`, then save the key pair created to your host
machine. This pair will be used to access the gateway, so ensure the pair is
saved in a safe and durable location.
3. Finish by selecting `Launch Instances`.

### ENI Configuration

Before installing Magma, we will add a second interface to gateway by creating
an ENI and attaching it to the EC2 instance.

1. In the EC2 service on AWS, navigate to the `Network Interfaces` section
under the `Network and Security` tab on the side panel.
2. Select `Create network interface` in the upper right corner.
3. On the `Create network interface` configuration page, select the subnet for
the ENI. To work properly, this subnet cannot be the same subnet that the
EC2 instance was deployed with. These subnets must be in the same availability
zone though.
4. Select the same subnet that was used to deploy the EC2 instance.
5. Once configured, select `Create network interface`.
6. Navigate to the EC2 instances page.
7. Find the recently deployed EC2 instance on the left hand side. Then select
`Actions` → `Networking` → `Attach network interface`.
8. On page `Attach network interface`, select the recently created ENI and then
click `Attach`.

### Install Magma

1. Find the public IP for the gateway instance by navigating to `Instances` on
the AWS EC2 service. Select the instance and copy the `Public IPv4 Address` in
the instance summary.
2. Add the AWS gateway key that was created when the instance was launched:
`ssh-add ~/.ssh/aws_key.pem`
3. SSH to EC2 instance using the public IP from step 1:
`ssh admin@<instance_public_ip>`
4. Now install Magma

```
[admin@<public_ip>~/] sudo su
[root@<public_ip>:/home/admin] wget https://raw.githubusercontent.com/facebookincubator/magma/v1.4/lte/gateway/deploy/agw_install.sh
[root@<public_ip>:/home/admin] bash agw_install cloud
`
```

When  you see "AGW installation is done." It means that your AGW installation
is done, you can make sure magma is running by executing:

```
service magma@* status
```

### Access Gateway Configuration

1. Follow the [configuration steps](https://docs.magmacore.org/docs/lte/deploy_config_agw) to register the new gateway.
2. To configure the gateway to serve as a secondary use the Orc8r API (NMS does
not currently support this functionality).
    1. Use the POST request endpoint `/lte/{network_id}/gateway_pools` to
    create a new gateway pool.
    2. Add the primary gateway(s) to the pool via endpoint
    `/lte/{network_id}/gateways/{gateway_id}/cellular/pooling`.
        1. MME code should differ for each gateway in the pool.
        2. MME relative capacity should be set to 255 for each primary
    3. Add the secondary (AWS) gateway to the pool via endpoint
    `/lte/{network_id}/gateways/{gateway_id}/cellular/pooling`.
        1. MME code should differ for each gateway in the pool.
        2. MME relative capacity should be set to 1 for the secondary
3. To enable secondary AGW to retrieve the connection state of the primary
instances, the default value of `use_ha: false` should be changed to
`use_ha: true` in `/etc/magma/mme.yml`. This configuration is mainly for
Active-Standby configuration and should not be used if an Active-Active
configuration is desired. When set as true, secondary AGW starts offloading UEs
camped on it back to the primary instances when the primary instances come back
up and start syncing up the states of connected eNBs to the orc8r.
4. If the secondary AGW is in a different network with its eth1 interface
configured with a private IP address, S1-U IP address needs to be configured
with the public IP address of the interface separately as by default it will be
configured with the eth1 IP address that is private.
    1. add "ipv4_sgw_s1u_addr": **** "IP_ADDRESS_STRING" via the endpoint
    `/lte/{network_id}/gateways/{gateway_id}/cellular/epc`, where
    IP_ADDRESS_STRING is a CIDR formatted IPv4 address, e.g., 203.0.113.25/32.
5. If eNB is behind a different NAT than the AGW instance, its S1-U IP address
communicated (with AGW instance) over the S1-MME interface is a private IP
address. Then, eNB will not be reachable in the user plane (i.e., GTP-U traffic
will not be routable back to eNB). To remedy this situation, assuming that the
eNB uses the same routable IP address for S1-MME connection and S1-U
 connection, it is possible to force MME overwrite the S1-U private IP address
 with the public one during bearer context set up by changing the
 `enable_gtpu_private_ip_correction: false` to
 `enable_gtpu_private_ip_correction: true` in `/etc/magma/mme.yml` after
 ssh-ing into the AGW instance.

Note: The current functionality supports multiple primaries using the same
secondary gateway. However the ENBs configured for the primaries must not
overlap.

### Enodeb Configuration

Any enodebs that will be used in the HA pool should be added to both the
primary and secondary gateway via the NMS.

Make sure that your eNB supports MME pooling also known as S1-Flex as Magma HA
feature relies on this capability. eNBs must be configured with MME pool using
the management interface for the eNB vendor. The primary and secondary AGW’s
routable ip addresses assigned for eth1 must be used in this configuration.
Make sure that eNB simultaneously connects to each MME ip address in its pool
and there are sctp heartbeat requests and responses on each AGW.


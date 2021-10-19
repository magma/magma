# Containerized AGW

## Build and push images

### Build
```
git clone https://github.com/magma/magma
cd lte/gateway/docker
docker-compose build
```

### Push
Login to the container registry and push
```
docker login example.repo/magma/
cd lte/gateway/docker
bash publish.sh example.repo/magma/
```

## Deploy Containerized AGW

### Using cloudstrapper to create an AGW instance on AWS CF
Deploy an instance using the [agw-provision playbook](https://github.com/magma/magma/blob/master/experimental/cloudstrapper/playbooks/agw-provision.yaml)

```
cd experimental/cloudstrapper/playbooks/
ansible-playbook agw-provision.yaml --tags createGw -e '@~/cluster.yaml'
```
Example cluster.yaml
```
---
dirLocalInventory: ~/cloudstrapper
awsAgwAmi: ami-0d058fe428540cd89
buildUbuntuAmi: ami-04fade045b8da506f
awsCloudstrapperAmi: ami-0fad3311309aca4c9
awsAgwRegion: ap-southeast-1
keyHost: rmelero
idSite: MenloPark
idGw: devops01-agw-deploy-test
```

This is how your `dirLocalInventory` should look

```
cat ~/cloudstrapper/secrets.yaml
---
awsAccessKey:
awsSecretKey:
```

After deploy, it might be necessary to resize ebs volume to 100 GB to accommodate the installation of docker and the build space for images.

Resize the disk partitions

```
growpart /dev/xvda 1; resize2fs /dev/xvda1
```

### Using a generic server with 2 NICs

Create/deploy an ubuntu 20.04 (latest version) server/instance with 2 interfaces (S1, Sgi) and the appropriate resources for your use case.

Your interfaces should be named eth0 and eth1. eth0 should be your Sgi interface, and eth1 should be S1 interface.

If your interfaces are not named eth0 and eth1, you will have to use netplan to rename your interfaces and reboot
```
cat /etc/netplan/50-cloud-init.yaml
network:
    ethernets:
        eth0:
            dhcp4: true
            match:
                macaddress: 08:00:27:87:21:53
            set-name: eth0
        eth1:
            dhcp4: true
            match:
                macaddress: 08:00:27:a5:27:4b
            set-name: eth1
```

Get the agw install script

```
wget https://github.com/magma/magma/raw/master/lte/gateway/deploy/agw_install_docker.sh
```

Add your rootCA.pem generated from your orc8r deployment
```
mkdir -p /var/opt/magma/certs
cp rootCA.pem /var/opt/magma/certs/rootCA.pem
```

Run install script to install docker, docker-compose, and bootstrap host

```
bash agw_install_docker.sh
```

After this, you can install your snowflake and gateway certs if you have them, or for new gateways, let them be generated.

```
cp snowflake /etc/snowflake
cp gateway.crt /var/opt/magma/certs
cp gateway.key /var/opt/magma/certs
```

Edit .env in /var/opt/magma/docker to have your docker registry values, and S1 and SGi interface IPs.

Pull images and start containers
```
cd /var/opt/magma/docker
./agw_upgrade.sh
```

You can check your connection status with
```
docker exec -it magmad /usr/local/bin/checkin_cli.py
1. -- Testing TCP connection to controller.orc8r-deployment.dev:8443 --
2. -- Testing Certificate --
3. -- Testing SSL --
4. -- Creating direct cloud checkin --
5. -- Creating proxy cloud checkin --

Success!
```

# PLMN changes may be required for testing

We tested with 00101 PLMN. Your testing might require changes to the GUMMEI_LIST, TAI_LIST, and TAC_LIST configuration items in lte/gateway/docker/mme/configs/mme.conf that must be edited before your build of the images.

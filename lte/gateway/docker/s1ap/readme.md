# Build and publish s1aptester images

1. Clone magma and move into the AGW docker directory and run build script.

```bash
git clone https://github.com/magma/magma
cd magma/lte/gateway/docker
s1ap/build-s1ap.sh
```

2. Publish images to your registry

```bash
s1ap/publish-s1ap.sh yourregistry.com/yourrepo/
```

# Deploy containerized S1APTester

## Create your environment

### AWS Environment

- VPC (vpc-s1ap) with a CIDR of 192.168.0.0/16
- Internet Gateway (igw1) attached to VPC (vpc-s1ap)
- Route table (rt1) with a default gateway to the Internet Gateway (igw1)
- 3 Subnets with an associated Route Table (rt1)
    - SGi: 192.168.60.0/24 (subnet-sgi)
    - S1: 192.168.128.0/23 (subnet-s1)
    - S1AP Management: 192.168.59.0/24 (subnet-management)
- Import/Create Key Pair (kp1)

### AGW instance

- Ubuntu 20.04 x86
- t2.large
- Security group
    - Allow SSH from My IP, or private address of your bastion host. **Important to not allow Public SSH access**
    - Allow all traffic from 192.168.60.0/24
    - Allow all traffic from 192.168.128.0/23
- 3 Interfaces
    - eth0: 192.168.129.1 (subnet-s1)
        - Add Public IP
    - eth1: 192.168.60.142 (subnet-sgi)
    - eth2: Auto-assign (subnet-management)
- 50 GB Disk
- Key Pair (kp1)

### s1aptester instance

- Ubuntu 20.04 x86
- t2.large
- Security group
    - Allow SSH from My IP, or private address of your bastion host. **Important to not allow Public SSH access**
    - Allow all traffic from 192.168.60.0/24
    - Allow all traffic from 192.168.128.0/23
- 3 Interfaces
    - eth0: Auto-assign (subnet-management)
        - Add Public IP
    - eth1: 192.168.60.141 (subnet-sgi)
    - eth2: 192.168.128.11 (subnet-s1)
- 50 GB Disk
- Key Pair (kp1)

### Traffic Generator instance

- Ubuntu 20.04 x86
- t2.large
- Security group
    - Allow SSH from My IP, or private address of your bastion host. **Important to not allow Public SSH access**
    - Allow all traffic from 192.168.60.0/24
    - Allow all traffic from 192.168.128.0/23
- 3 Interfaces
    - eth0: Auto-assign (subnet-management)
        - Add Public IP
    - eth1: 192.168.60.144 (subnet-sgi)
    - eth2: 192.168.129.42 (subnet-s1)
- 50 GB Disk
- Key Pair (kp1)

## Add your ssh key pair (kp1) used when creating the instances to your ssh-agent

```bash
ssh-add .ssh/id_rsa
ssh-add -l
3072 SHA256:KWaiSLPiTYy5tRGXSKpkuwJcYDGy92NTMi0v3khLjJY (RSA)
```

## Run S1APTester playbook

### ssh to the agw instance and deploy as a containerized AGW

```bash
# ssh to AGW instance using the elastic IP address from AWS. Notice the '-A' argument. You must load in the key pair that you created in AWS for it to forward your ssh agent key.
ssh -A ubuntu@${YOUR_AWS_IP}

# Download containerized AGW bootstrap script
wget https://raw.githubusercontent.com/magma/magma/master/lte/gateway/deploy/agw_install_docker.sh

# Make bootstrap script executable
chmod +x agw_install_docker.sh

# Create certs directory and populate with your orc8r rootCA. Note that the certificate contents here have been omitted.
sudo mkdir -p /var/opt/magma/certs
sudo su root -c 'sudo echo "-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----" > /var/opt/magma/certs/rootCA.pem'

# Run bootstrap script and reboot
sudo agw_install_docker.sh
sudo reboot
```

### After reboot, add the s1ap and trfgen ssh keys to known_hosts of the AGW Instance

```bash
ssh -A ubuntu@${YOUR_AWS_IP}
ssh-keyscan -H 192.168.60.144 >> ~/.ssh/known_hosts
ssh-keyscan -H 192.168.60.141 >> ~/.ssh/known_hosts
```

### move to deploy folder and test your connection to s1ap and trfgen instances

```bash
cd /opt/magma/lte/gateway/deploy

ansible -m ping all -i agw_hosts -b
127.0.0.1 | SUCCESS => {
192.168.60.141 | SUCCESS => {
192.168.60.144 | SUCCESS => {
```

### Deploy s1aptester

Depending on whether you already have your AGW and s1aptest images in your registry, or want to build them, your can supply the `docker_registry` environment variable to the playbook. If you don't supply the variable, it will build the images and start the services. If you supply it, it will pull down the images from the registry and start the services.

#### Without registry

```bash
ansible-playbook -i agw_hosts magma_docker_s1ap_setup.yml -b
```bash

#### With registry

```bash
ansible-playbook -i agw_hosts magma_docker_s1ap_setup.yml -b -e "docker_registry=public.ecr.aws/yourrepo/"
```

### Running tests Commands

#### ssh to the s1aptester instance and attach to the s1aptester container to run tests

```bash
# ssh from AGW to s1aptester instance
ssh 192.168.60.141
sudo docker attach s1aptester
make integ_test TESTS=s1aptests/test_attach_detach.py
```

#### You could also exec the test running command from the AGW instance

```bash
# Run this from the AGW
ssh 192.168.60.141 sudo docker exec -t s1aptester make integ_test TESTS=s1aptests/test_attach_detach.py
```

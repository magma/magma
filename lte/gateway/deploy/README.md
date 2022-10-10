# I. Deploy AGW as Docker Containers

## 1. Pre-requisites

### a. rootCa.pem
Have already copied the `rootCA.pem` certificate from the orc8r to `/var/opt/magma/certs`.
To do that, please proceed as follow:
```bash
mkdir -p /var/opt/magma/certs/
scp -p <ORC8R>:/root/magma-galaxy/secrets/rootCA.pem /var/opt/magma/certs/
```

**NOTE**: Replace <ORC8R> with the appropriate username and IP address.

### b. Export Variables
```bash
export MAGMA_ROOT=/opt/magma
```

## 2. AGW Configuration
### a. Verify the sanity of `rootCA.pem` file
```bash
openssl x509 -text -noout -in /var/opt/magma/certs/rootCA.pem
```
### b. Download script `agw_install_docker.sh`
```bash
cd  
wget  https://github.com/magma/magma/raw/master/lte/gateway/deploy/agw_install_docker.sh
```
### c. Update the script `agw_install_docker.sh`
```bash
sed -i "113 a\sed -i 's/debian/debian-test/' /opt/magma/lte/gateway/deploy/roles/magma_deploy/vars/all.yaml\n\
sed -i 's/focal-1.7.0/focal-ci/' /opt/magma/lte/gateway/deploy/roles/magma_deploy/vars/all.yaml\n\
       " "agw_install_docker.sh"
```
### d. Run the script `agw_install_docker.sh`
```bash
chmod +x agw_install_docker.sh
./agw_install_docker.sh
```
### e. Reboot the VM
```bash
reboot
```
### f. Update `/var/opt/magma/configs/control_proxy.yml` with your orc8r/controller details
```bash
sudo mkdir -p /var/opt/magma/configs/
sudo vim /var/opt/magma/configs/control_proxy.yml

cloud_address: controller.tchitchir.com
cloud_port: 443
bootstrap_address: bootstrapper-controller.tchitchir.com
bootstrap_port: 443
fluentd_address: fluentd.tchitchir.com
fluentd_port: 443
rootca_cert: /var/opt/magma/certs/rootCA.pem

```
### g. Update .env file
```bash
vim $MAGMA_ROOT/lte/gateway/docker/.env

COMPOSE_PROJECT_NAME=agw
DOCKER_USERNAME=
DOCKER_PASSWORD=
DOCKER_REGISTRY=
IMAGE_VERSION=latest
BUILD_CONTEXT=../../..
ROOTCA_PATH=/var/opt/magma/certs/rootCA.pem
CONTROL_PROXY_PATH=/var/opt/magma/configs/control_proxy.yml 
SNOWFLAKE_PATH=/etc/snowflake
CERTS_VOLUME=/var/opt/magma/certs
CONFIGS_OVERRIDE_VOLUME=/var/opt/magma/configs
CONFIGS_OVERRIDE_TMP_VOLUME=/var/opt/magma/tmp
CONFIGS_DEFAULT_VOLUME=../configs
CONFIGS_TEMPLATES_PATH=/etc/magma/templates
LOG_DRIVER=journald
```
### h. Build AGW docker images 
#### h.1 For ARM architecture
```bash
cd ${MAGMA_ROOT}/lte/gateway/docker/
sudo docker-compose build --build-arg CPU_ARCH=aarch64 --build-arg DEB_PORT=arm64 --parallel
```
#### h.2 For X86 archtecture
```bash
cd ${MAGMA_ROOT}/lte/gateway/docker/
sudo docker-compose build --parallel
```

### i. Start the services
```bash
cd ${MAGMA_ROOT}/lte/gateway/docker/
docker-compose --env-file .env -f docker-compose.yaml -f docker-compose.dev.yaml up -d --force-recreate
```

**Note**: In case you want to use the docker images from the artifactory:
You have to change directory to `/var/opt/magma/docker`

```bash
cd /var/opt/magma/docker
vim .env

COMPOSE_PROJECT_NAME=agw
DOCKER_REGISTRY=docker.artifactory.magmacore.org/
DOCKER_USERNAME=
DOCKER_PASSWORD=
IMAGE_VERSION=latest
OPTIONAL_ARCH_POSTFIX=_arm
ROOTCA_PATH=/var/opt/magma/certs/rootCA.pem
CONTROL_PROXY_PATH=/var/opt/magma/configs/control_proxy.yml
CONFIGS_TEMPLATES_PATH=/etc/magma/templates
CERTS_VOLUME=/var/opt/magma/certs
CONFIGS_OVERRIDE_VOLUME=/var/opt/magma/configs
CONFIGS_OVERRIDE_TMP_VOLUME=/var/opt/magma/tmp
CONFIGS_DEFAULT_VOLUME=../configs
SECRETS_VOLUME=/var/opt/magma/secrets
HOST_DOCKER_INTERNAL=192.168.129.1
LOG_DRIVER=journald
```

After that you run the script `agw_upgrade.sh` or the compose file:
```bash
./agw_upgrade.sh
```

(or)

```bash
docker-compose down
docker-compose --env-file .env -f docker-compose.yaml -f docker-compose.dev.yaml up -d --force-recreate
```

# II. Install S1apTester
## 1. Pre-requisites
### a. Clone the magma project
```bash
ubuntu@magma-trfserver:~$ git clone `https://github.com/magma/magma.git`
```

## 2. Run the Ansible Playbook
To install the S1apTester, we need to run the Ansible Playbook `s1aptester.yml` with the inventory `hosts`
```bash
cd $MAGMA_ROOT/lte/gateway/deploy
ansible-playbook -i hosts s1aptester.yml -K --key-file "~/.ssh/arm-key-pair-s1ap.pem"
```

## 3. Run Tests
### a. Export variables
```bash
export MAGMA_ROOT=$HOME/magma
export c_build=/home/ubuntu/build/c
export oai_build=$c_build/core/oai
export S1AP_TESTER_ROOT=~/s1ap-tester
export S1AP_TESTER_SRC=~/S1APTester
export PIP_CACHE_HOME=~/.pipcache
export PYTHON_BUILD=/build/python
export PIP_CACHE_HOME=~/.pipcache
export CODEGEN_ROOT=/var/tmp/codegen
export SWAGGER_CODEGEN_DIR=${CODEGEN_ROOT}/modules/swagger-codegen-cli/target
export SWAGGER_CODEGEN_JAR=${SWAGGER_CODEGEN_DIR}/swagger-codegen-cli.jar
```
## b. Install some packages  
```bash
sudo pip install virtualenv
```
## c. Create the folder `/build/python`
On the s1aptester instance, create the path /build/python
```bash
sudo mkdir -p /build/python
sudo chmod +w /build/python
```

## d. Run Makefile
in the s1aptester instance, you have to build the tests using Makefiles
```bash
cd $MAGMA_ROOT/lte/gateway/python && make
cd integ_tests && make
```
## e. Tests
On the s1aptester instance, run either individual tests or the full suite of tests.
A safe, non-flaky test to run is `s1aptests/test_attach_detach.py`.
#### Individual test(s)
 ```bash
make integ_test TESTS=<test(s)_to_run>
```
#### All Sanity tests
```bash
make integ_test
```
#### All Non-Sanity tests
```bash
make nonsanity
```
#### Minimal set of tests to be executed before committing changes to magma repository
```bash
make precommit
```
#### Run with -i flag to enable continuous test runs (ignoring the failing test(s), if any)
```bash
make -i precommit or make -i integ_test
```
#### Set enable-flaky-retry=true to re-run the failing test(s) to identify flaky behavior
```bash
make precommit enable-flaky-retry=true or make integ_test enable-flaky-retry=true
```
## f. Running uplink/downlink traffic tests
    On the agw instance, run, disable-tcp-checksumming
    On the s1aptester instance, disable-tcp-checksumming
    On the trfgen instance, run disable-tcp-checksumming && trfgen-server
Running make integ_test in magma_test VM should succeed now.

# III. Install TrfServer
In this section, you will learn how to deploy and use the TrfServer.
## 1. Pre-requisites
Please make sure that the following prerequisites are satisfied.
### a. Clone the magma project
Please clone the magma project to the path `/home/{{ ansible_user }}/`
**Note**: This path is the one used in the playbooks for the variable magma_root.
```bash
ubuntu@magma-trfserver:~$ git clone https://github.com/magma/magma.git
```
### b. Export MAGMA_ROOT variable
```bash
export MAGMA_ROOT=$HOME/magma
```
## 2. Run the Ansible Playbook
To install the TrfServer, we need to run the Ansible Playbook `iperfserver.yml` with the inventory `hosts`
```bash
cd $MAGMA_ROOT/lte/gateway/deploy
ansible-playbook -i hosts iperfserver.yml -K --key-file "~/.ssh/arm-key-pair-s1ap.pem"
```
## 3. Running uplink/downlink traffic tests
To run the uplink/downlink traffic tests, you need to disable the TCP checksum, then run the iperf on the TrfServer VM.
```bash
disable-tcp-checksumming && trfgen-server
```

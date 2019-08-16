#!/bin/bash
#
# Provided by IRSOLS Inc
# Setup Docker CE instead of vanilla centos repos
# version 0.4
# last modified 07/21/2019

#TODO: Check OS version and change installer accordingly
# Check Distro
export DISTRO=`cat /etc/*release| grep centos | head -1 | sed s/ID=//g | sed s/\"//g`

if [[ $DISTRO = centos ]] ; then
    echo "distro is:  $DISTRO"
 yum install -y open-vm-tools
 yum remove docker \
                   docker-client \
                   docker-client-latest \
                   docker-common \
                   docker-latest \
                   docker-latest-logrotate \
                   docker-logrotate \
                   docker-engine
 # Setup Docker repository
 yum install -y yum-utils \
   device-mapper-persistent-data \
   lvm2

 # Setup stable repo
 yum-config-manager \
     --add-repo \
     https://download.docker.com/linux/centos/docker-ce.repo

 # Install the latest version of Docker CE and containerd
 yum -y install docker-ce docker-ce-cli containerd.io

 # Start and enable Docker Daemon
 systemctl start docker
 systemctl enable docker

 # Install and setup  Docker Compose
 sudo curl -L "https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
 sudo chmod +x /usr/local/bin/docker-compose
 docker-compose --version


 # Install and setup docker machine
 base=https://github.com/docker/machine/releases/download/v0.16.0 &&
   curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/tmp/docker-machine &&
   sudo install /tmp/docker-machine /usr/local/bin/docker-machine
 docker-machine version

else
    echo "distro is: not Centos"
    echo "could be ubuntu, so we'll try apt "
#### If the Distro is Ubuntu do this following ...
# Uninstall Older / Vanilla versions of Docker
sudo apt-get remove docker docker-engine docker.io containerd runc
sudo apt-get -y update
sudo apt-get -y install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
# Install packages to setup custom repos and set up Docker Official repos
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get -y update
sudo apt-get -y install docker-ce docker-ce-cli containerd.io


fi
exit



#!/bin/bash
# " This multi-stage installer assumes following are already"
# " setup/configured: "
# " -  Fresh Centos 7 min base installation "
# " -  Nested virtualization support to spinup vagrant VM's "
# "     if installing in a VM"
# " N.B. Some of the steps need reboot and manual intervention"
# " Proceed with caution ......"
echo " Welcome to IRSOLS Installer for Facebook Magma PacketCore"
echo " This multi-stage installer will setup the necessary required "
echo " packages,OS and other configurations necessary for a "
echo " successful installation of Magma-Core. " 
echo
echo " For support requests please use github issues "
echo " Copyright IRSOLS Inc . info@irsols.com "

# Check CPU, Mem & Virtualization Capabilities. Abort if insufficient
# Comment out if you want to bypass these checks at your own risk

./check-host-reqs.sh

###############################################################
# Stage 1 : Sys prep OS and install virtualbox pre-requisites
###############################################################


echo " Sys Prep Centos 7 base OS .." 
yum update -y
yum groupinstall 'Development Tools'
yum -y install gcc make patch  dkms qt libgomp
yum -y install kernel-headers kernel-devel fontforge binutils glibc-headers glibc-devel epel-release
cd /etc/yum.repos.d
# Disable firewalld so it doesnt interfere on intra-vm traffic
systemctl disable firewalld --now

echo " Manual Steps " 
echo " Follow these steps MANUALLY, once VirtualBox is installed and "
echo " VM is rebooted , continue with setting up Vagrant below "
echo " Starting VirtualBox setup... "
wget http://download.virtualbox.org/virtualbox/rpm/rhel/virtualbox.repo

echo " Exiting installer from Stage 1," 
echo " Please Reboot your host and continue with Stage 2 below:"
exit 
#reboot
###############################################################
# Stage 2 : Install VirtualBox AFTER reboot, install Vagrant, 
#           install Docker CE/Compose and Setup Magma Pre-Reqs
###############################################################

export KERN_DIR=/usr/src/kernels/$(uname -r)
export KERN_VER=$(uname -r)
yum install VirtualBox-5.2 -y
/sbin/rcvboxdrv setup

echo " Starting Vagrant setup... "
# Setup Vagrant
wget https://releases.hashicorp.com/vagrant/2.2.2/vagrant_2.2.2_x86_64.rpm
yum –y localinstall vagrant_2.2.2_x86_64.rpm
vagrant --version

# Setup Docker-CE, Compose and Machine properly
./docker-clean-setup.sh

echo "  Setup Magma Pre-reqs "
# First need to setup the Software Collections
yum install centos-release-scl
yum install rh-python36

echo "Exiting installer from Stage 2 "
echo "Continue manually with Stage 3 below"
exit 

###############################################################
# Stage 3 : Configure following manually . once you 
#           do 'scl enable rh-python36 bash' 
#           rest needs to be executed within that shell. 
############################################################### 

scl enable rh-python36 bash
echo " Verify Python version and ensure we're uing python3.6x "
python --version
pip3.6 install --upgrade pip
pip3.6 install ansible fabric3 requests PyYAML
vagrant plugin install vagrant-vbguest

# Replace pip3 environment with pip3.6 env in the ocr8/cloud/docker/build.py



#!/bin/sh
# Setting up env variable, user and project path
MAGMA_USER="magma"
AGW_INSTALL_CONFIG="/etc/systemd/system/multi-user.target.wants/agw_installation.service"
AGW_SCRIPT_PATH="/root/agw_install.sh"
DEPLOY_PATH="/home/$MAGMA_USER/magma/lte/gateway/deploy"
SUCCESS_MESSAGE="ok"
NEED_REBOOT=0
WHOAMI=$(whoami)
KVERS=$(uname -r)

echo "Checking if the script has been executed by root user"
if [ "$WHOAMI" != "root" ]; then
  echo "You're executing the script as $WHOAMI instead of root.. exiting"
  exit 1
fi

echo "Checking if Debian is installed"
if ! grep -q 'Debian' /etc/issue; then
  echo "Debian is not installed"
  exit 1
fi

echo "Making sure $MAGMA_USER user is sudoers"
if ! grep -q "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
  adduser $MAGMA_USER sudo
  echo "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
fi

echo "Need to check if both interfaces are named eth0 and eth1"
INTERFACES=$(ip -br a)
if [[ ! $INTERFACES == *'eth0'*  ]] || [[ ! $INTERFACES == *'eth1'* ]] || ! grep -q 'GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"' /etc/default/grub; then
  # changing intefaces name
  sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"/g' /etc/default/grub
  # changing interface name
  grub-mkconfig -o /boot/grub/grub.cfg
  sed -i 's/enp1s0/eth0/g' /etc/network/interfaces
  # configuring eth1
  echo "auto eth1
  iface eth1 inet static
  address 10.0.2.1
  netmask 255.255.255.0" > /etc/network/interfaces.d/eth1
  # Setting REBOOT flag to 1 because we need to reload new interface and network services.
  NEED_REBOOT=1
else
  echo "Interfaces name are correct, let's check if network and DNS are up"
  while ! ping -c 1 -W 1 -I eth0 google.com; do
    echo "Network not ready yet"
    sleep 1
  done
fi

echo "Checking if the righ kernel version is installed (4.9.0-9-amd64)"
if [ "$KVERS" != "4.9.0-9-amd64" ]; then
  # Adding the snapshot to retrieve 4.9.0-9-amd64
  if ! grep -q "deb http://snapshot.debian.org/archive/debian/20190801T025637Z" /etc/apt/sources.list; then
    echo "deb http://snapshot.debian.org/archive/debian/20190801T025637Z stretch main non-free contrib" >> /etc/apt/sources.list
  fi
  apt update
  # Installing prerequesites, Kvers, headers
  apt install -y sudo python-minimal aptitude linux-image-4.9.0-9-amd64 linux-headers-4.9.0-9-amd64
  # Removing dev repository snapshot from source.list
  sed -i '/20190801T025637Z/d' /etc/apt/sources.list
  # Removing incompatible Kernel version
  DEBIAN_FRONTEND=noninteractive apt remove -y linux-image-"$KVERS"-amd64
  # Setting REBOOT flag to 1 because we need to boot with the right Kernel version
  NEED_REBOOT=1
fi

if [ $NEED_REBOOT = 1 ]; then
  echo "Will reboot in a few seconds, loading a boot script in order to install magma"
  if [ ! -f "$AGW_SCRIPT_PATH" ]; then
    wget --no-cache -O $AGW_SCRIPT_PATH https://raw.githubusercontent.com/facebookincubator/magma/master/lte/gateway/deploy/agw_install.sh
  fi
  echo "[Unit]
Description=AGW Installation
After=network-online.target
Wants=network-online.target
[Service]
Type=oneshot
ExecStart=/bin/sh /root/agw_install.sh
User=root
Group=root
[Install]
WantedBy=multi-user.target" > $AGW_INSTALL_CONFIG
  chmod 644 $AGW_INSTALL_CONFIG
  reboot
fi

echo "Making sure eth0 is connected to internet"
PING_RESULT=$(ping -c 1 -I eth0 8.8.8.8 > /dev/null 2>&1 && echo "$SUCCESS_MESSAGE")
if [ "$PING_RESULT" != "$SUCCESS_MESSAGE" ]; then
  echo "eth0 (enp1s0) is not connected to internet, please double check your plugged wires."
  exit 1
fi
echo "Checking if magma has been installed"
MAGMA_INSTALLED=$(apt-cache show magma >  /dev/null 2>&1 echo "$SUCCESS_MESSAGE")
if [ "$MAGMA_INSTALLED" != "$SUCCESS_MESSAGE" ]; then
  echo "Magma not installed processing installation"
  apt-get update
  apt-get -y install curl make virtualenv zip rsync git software-properties-common python3-pip python-dev
  alias python=python3
  pip3 install ansible

  git clone https://github.com/facebookincubator/magma.git /home/$MAGMA_USER/magma
  cd /home/$MAGMA_USER/magma
  git checkout v1.0.1

  echo "Generating localhost hostfile for Ansible"
  echo "[ovs_build]
  127.0.0.1 ansible_connection=local
  [ovs_deploy]
  127.0.0.1 ansible_connection=local" > $DEPLOY_PATH/agw_hosts
  echo "Triggering ovs_build playbook"
  su - $MAGMA_USER -c "ansible-playbook -e \"MAGMA_ROOT='/home/$MAGMA_USER/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_hosts $DEPLOY_PATH/ovs_build.yml"
  echo "Triggering ovs_deploy playbook"
  su - $MAGMA_USER -c "ansible-playbook -e \"PACKAGE_LOCATION='/tmp'\" -i $DEPLOY_PATH/agw_hosts $DEPLOY_PATH/ovs_deploy.yml"
  echo "Deleting boot script if it exists"
  if [ -f "$AGW_INSTALL_CONFIG" ]; then
    rm -rf $AGW_INSTALL_CONFIG
    systemctl daemon-reload
  fi
  echo "Removing Ansible from the machine."
  pip3 uninstall --yes ansible
  rm -rf /home/$MAGMA_USER/build
  service magma@* status
  echo "AGW installation is done, make sure all services above are running correctly"
else
  echo "Magma already installed, skipping.."
fi

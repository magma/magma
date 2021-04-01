#!/bin/bash
# Setting up env variable, user and project path
MAGMA_USER="magma"
AGW_INSTALL_CONFIG="/etc/systemd/system/multi-user.target.wants/agw_installation.service"
AGW_SCRIPT_PATH="/root/agw_install.sh"
DEPLOY_PATH="/home/$MAGMA_USER/magma/lte/gateway/deploy"
SUCCESS_MESSAGE="ok"
NEED_REBOOT=0
WHOAMI=$(whoami)
KVERS=$(uname -r)
MAGMA_VERSION="${MAGMA_VERSION:-v1.4}"
CLOUD_INSTALL="cloud"
GIT_URL="${GIT_URL:-https://github.com/magma/magma.git}"



echo "Checking if the script has been executed by root user"
if [ "$WHOAMI" != "root" ]; then
  echo "You're executing the script as $WHOAMI instead of root.. exiting"
  exit 1
fi

wget https://raw.githubusercontent.com/magma/magma/"$MAGMA_VERSION"/lte/gateway/deploy/agw_pre_check.sh
if [[ -f ./agw_pre_check.sh ]]; then
  chmod 644 agw_pre_check.sh && bash agw_pre_check.sh
  while true; do
      read -p "Do you accept those modifications and want to proceed with magma installation?(y/n)" yn
      case $yn in
          [Yy]* ) break;;
          [Nn]* ) exit;;
          * ) echo "Please answer yes or no.";;
      esac
  done
else
  echo "agw_precheck.sh is not available in your version"
fi


echo "Checking if Debian is installed"
if ! grep -q 'Debian' /etc/issue; then
  echo "Debian is not installed"
  exit 1
fi

echo "Making sure $MAGMA_USER user is sudoers"
if ! grep -q "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" /etc/sudoers; then
  apt install -y sudo
  adduser --disabled-password --gecos "" $MAGMA_USER
  adduser $MAGMA_USER sudo
  echo "$MAGMA_USER ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
fi

echo "Need to check if both interfaces are named eth0 and eth1"
INTERFACES=$(ip -br a)
if [[ $1 != "$CLOUD_INSTALL" ]] && ( [[ ! $INTERFACES == *'eth0'*  ]] || [[ ! $INTERFACES == *'eth1'* ]] || ! grep -q 'GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"' /etc/default/grub); then
  # changing intefaces name
  sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="net.ifnames=0 biosdevname=0"/g' /etc/default/grub
  # changing interface name
  grub-mkconfig -o /boot/grub/grub.cfg
  echo "auto eth0
  iface eth0 inet dhcp" > /etc/network/interfaces.d/eth0
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

echo "Checking if the right kernel version is installed (4.9.0-9-amd64)"
if [ "$KVERS" != "4.9.0-9-amd64" ]; then
  # Adding the snapshot to retrieve 4.9.0-9-amd64
  if ! grep -q "deb http://snapshot.debian.org/archive/debian/20190801T025637Z" /etc/apt/sources.list; then
    echo "deb http://snapshot.debian.org/archive/debian/20190801T025637Z stretch main non-free contrib" >> /etc/apt/sources.list
  fi
  apt update
  # Installing prerequesites, Kvers, headers
  apt install -y python-minimal aptitude linux-image-4.9.0-9-amd64 linux-headers-4.9.0-9-amd64
  # Removing dev repository snapshot from source.list
  sed -i '/20190801T025637Z/d' /etc/apt/sources.list
  # Removing incompatible Kernel version
  DEBIAN_FRONTEND=noninteractive apt remove -y linux-image-"$KVERS"
  # Setting REBOOT flag to 1 because we need to boot with the right Kernel version
  NEED_REBOOT=1
fi

# configure environment variable defaults needed for ansible
ANSIBLE_VARS="PACKAGE_LOCATION=/tmp"
if [ -n "${REPO_HOST}" ]; then
    if [ -z "${REPO_PROTO}" ]; then
        REPO_PROTO=http
    fi
    if [ -z "${REPO_DIST}" ]; then
        REPO_DIST=stretch-stable
    fi
    if [ -z "${REPO_COMPONENT}" ]; then
        REPO_COMPONENT=main
    fi
    # configure pkgrepo location
    ANSIBLE_VARS="ovs_pkgrepo_proto=${REPO_PROTO} ovs_pkgrepo_host=${REPO_HOST} ovs_pkgrepo_path=${REPO_PATH} ${ANSIBLE_VARS}"

    # configure pkgrepo distribution
    ANSIBLE_VARS="ovs_pkgrepo_dist=${REPO_DIST} ovs_pkgrepo_component=${REPO_COMPONENT} ${ANSIBLE_VARS}"

    # configure pkgrepo gpg key
    ANSIBLE_VARS="ovs_pkgrepo_key=${REPO_KEY} ${ANSIBLE_VARS}"
    if [ -z "${REPO_KEY_FINGERPRINT}" ]; then
        ANSIBLE_VARS="ovs_pkgrepo_key_fingerprint=${REPO_KEY_FINGERPRINT} ${ANSIBLE_VARS}"
    fi
fi

if [[ "${REPO_PROTO}" == 'https' ]]; then
    echo "Ensure HTTPS apt transport method is installed"
    apt install -y apt-transport-https
fi

if [ $NEED_REBOOT = 1 ]; then
  echo "Will reboot in a few seconds, loading a boot script in order to install magma"
  if [ ! -f "$AGW_SCRIPT_PATH" ]; then
      cp "$(realpath $0)" "${AGW_SCRIPT_PATH}"
  fi
  cat <<EOF > $AGW_INSTALL_CONFIG
[Unit]
Description=AGW Installation
After=network-online.target
Wants=network-online.target
[Service]
Environment=MAGMA_VERSION=${MAGMA_VERSION}
Environment=GIT_URL=${GIT_URL}
Environment=REPO_PROTO=${REPO_PROTO}
Environment=REPO_HOST=${REPO_HOST}
Environment=REPO_DIST=${REPO_DIST}
Environment=REPO_COMPONENT=${REPO_COMPONENT}
Environment=REPO_KEY=${REPO_KEY}
Environment=REPO_KEY_FINGERPRINT=${REPO_KEY_FINGERPRINT}
Type=oneshot
ExecStart=/bin/bash ${AGW_SCRIPT_PATH}
TimeoutStartSec=3800
TimeoutSec=3600
User=root
Group=root
[Install]
WantedBy=multi-user.target
EOF
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
  echo "Magma not installed, processing installation"
  apt-get update
  apt-get -y install curl make virtualenv zip rsync git software-properties-common python3-pip python-dev
  alias python=python3
  pip3 install ansible

  git clone "${GIT_URL}" /home/$MAGMA_USER/magma
  cd /home/$MAGMA_USER/magma
  git checkout "$MAGMA_VERSION"

  echo "Generating localhost hostfile for Ansible"
  echo "[ovs_build]
  127.0.0.1 ansible_connection=local
  [ovs_deploy]
  127.0.0.1 ansible_connection=local" > $DEPLOY_PATH/agw_hosts
  if [ -n "${FORCE_OVS_BUILD}" ]; then
      echo "Triggering ovs_build playbook"
      su - $MAGMA_USER -c "ansible-playbook -e \"MAGMA_ROOT='/home/$MAGMA_USER/magma' OUTPUT_DIR='/tmp'\" -i $DEPLOY_PATH/agw_hosts $DEPLOY_PATH/ovs_build.yml"
      ANSIBLE_VARS="${ANSIBLE_VARS} ovs_use_pkgrepo=no"
  fi
  echo "Triggering ovs_deploy playbook"
  if [[ $1 == "$CLOUD_INSTALL" ]]; then
      su - $MAGMA_USER -c "ansible-playbook -e '${ANSIBLE_VARS}' -i $DEPLOY_PATH/agw_hosts $DEPLOY_PATH/ovs_deploy.yml --skip-tags \"skipfirstinstall\""
      su - $MAGMA_USER -c "ansible-playbook -e '${ANSIBLE_VARS}' -i $DEPLOY_PATH/agw_hosts $DEPLOY_PATH/ovs_deploy.yml"
      service openvswitch-switch restart
  else
      su - $MAGMA_USER -c "ansible-playbook -e '${ANSIBLE_VARS}' -i $DEPLOY_PATH/agw_hosts $DEPLOY_PATH/ovs_deploy.yml --skip-tags \"skipfirstinstall\""
  fi
  echo "Deleting boot script if it exists"
  if [ -f "$AGW_INSTALL_CONFIG" ]; then
    rm -rf $AGW_INSTALL_CONFIG
  fi
  rm -rf /home/$MAGMA_USER/build
  echo "AGW installation is done, make sure all services above are running correctly.. rebooting"
  reboot
else
  echo "Magma already installed, skipping.."
fi

#!/bin/bash
function check_success {
  ret=$?
  if [[ $ret == 0 ]]; then
    return 0
  fi
  echo  "$1 failed with return code $ret"
  exit 1
}

# Remove MME dependencies on Pipelined and Sessiond
sed '/magma@pipelined/d' -i /etc/systemd/system/magma@mme.service
sed '/magma@sessiond/d' -i /etc/systemd/system/magma@mme.service
check_success "Removing Sessiond and Pipelined dependencies in MME"

# Remove unused services from Magmad config
sed '/monitord/d' -i /etc/magma/magmad.yml
sed '/pipelined/d' -i /etc/magma/magmad.yml
sed '/policydb/d' -i /etc/magma/magmad.yml
sed '/sessiond/d' -i /etc/magma/magmad.yml
check_success "Removing services from Magmad config"

# Remove systemd files for unused services
rm /etc/systemd/system/magma@pipelined.service
rm /etc/systemd/system/magma@sessiond.service
check_success "Removing systemd files for unused Magma services"

# Reload systemd service files
systemctl daemon-reload
check_success "Reloading systemctl"

# Remove the openvswitch gtp bridge
ifdown gtp_br0

# Remove the oai-gtp package which installs the custom gtp module
apt-get -y remove oai-gtp
check_success "Removing oai-gtp package"

# Remove the custom gtp module needed by openvswitch
rmmod vport_gtp
rmmod openvswitch
rmmod gtp

# Install the default kernel gtp module
insmod /lib/modules/`uname -r`/kernel/drivers/net/gtp.ko
check_success "Installing kernel gtp module"

# Update the module symbols
depmod -a

# Enable NATing on the SGI interface, i.e. eth2
iptables -t nat -A POSTROUTING -o eth2 -j MASQUERADE
check_success "Installing NAT rule for eth2"

# Install libgtpnl
bash /home/vagrant/magma/third_party/libgtpnl/install.sh
check_success "Installing libgtpnl"

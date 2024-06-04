#!/bin/bash

is_kernel_module_loaded() {

  local -r module=$1

  if lsmod | grep -Eq "^$module\\s+" || [[ -d "/sys/module/$module" ]]; then
    echo "kernel module $module is loaded"
    return 0
  fi

  echo "kernel module $module is missing"
  return 1
}

assert_kernel_mod() {

  local -r module=$1

  echo "Checking if kernel module \"$module\" is loaded..."
  if is_kernel_module_loaded "$module"; then
    return
  fi

  echo "Attempting to load kernel module $module"
  modprobe -v "$module"
  if [[ $? -ne 0 ]]; then
    echo "WARNING: Unable to dynamically load kernel module $module."


    echo "Attempting to build and install dkms module."
    dkms autoinstall
    check_mod_version
    modprobe -v "$module"
    if [[ $? -eq 0 ]]; then
      if is_kernel_module_loaded "$module"; then
        echo "kernel module $module is loaded!"
        return
      else
        exit 1
      fi
    else
      echo "ERROR: Failed to install kernel module $module"
      exit 1
    fi
  fi
}

check_mod_version () {
  echo "Checking installed and loaded ovs modules versions"
  
  kernel_ver=$(cat /sys/module/openvswitch/srcversion)
  mod_ver=$(modinfo /lib/modules/"$(uname -r)"/updates/dkms/openvswitch.ko |grep srcversion|awk '{print $2}')

  if [[ "$kernel_ver" == "$mod_ver" ]]; then
    OVS_VER=$(dpkg -l openvswitch-datapath-dkms |grep 'ii' |awk '{print $3}')
    echo "Module update successful, installed: $OVS_VER"
    return 0
  else
    echo "FAIL: Update failed. Installed and loaded openvswitch module version mismatch!"
    return 1
fi
}

trap stopit SIGINT;
run=1;
stopit() {
  run=0;
};


firstrun=0
[ -f "/etc/openvswitch/conf.db" ] || firstrun=1

if [ $firstrun -ne 0 ]; then
  # create database
  ovsdb-tool create /etc/openvswitch/conf.db /usr/share/openvswitch/vswitch.ovsschema
fi

# Check if openvswitch kernel modules are loaded
assert_kernel_mod openvswitch
assert_kernel_mod vport_gtp

check_mod_version

# Start openvswitch daemons
if [ ! -d "/var/log/openvswitch" ]; then
  mkdir -p /var/log/openvswitch
fi

if [ ! -d "/var/run/openvswitch" ]; then
  mkdir -p /var/run/openvswitch
fi

echo "Starting service openvswitch-switch"
ovsdb-server --detach --pidfile --remote=punix:/var/run/openvswitch/db.sock \
             --log-file --verbose=off:syslog
ovs-vswitchd --detach --pidfile --log-file --verbose=off:syslog

# Activate bridge interfaces
echo "Activating bridge interfaces"
ifup --allow=ovs $(ifquery --allow ovs --list)
sleep 1
ifup --allow=ovs $(ifquery --allow gtp_br0 --list)
sleep 1
ifup --allow=ovs $(ifquery --allow uplink_br0 --list)

# echo "Bring up gtp_br0"
# ifup gtp_br0

# echo "Bring up uplink_br0"
# ifup uplink_br0

# echo "Bring up mtr0"
# ifup mtr0

# echo "Bring up ipfix0"
# ifup ipfix0

# echo "Bring up dhcp0"
# ifup dhcp0

while [ $run ]; do
	sleep .1
done ;

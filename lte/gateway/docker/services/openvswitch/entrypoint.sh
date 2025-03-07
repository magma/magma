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

load_kernel_mod() {

  local -r module=$1

  echo "Checking if kernel module \"$module\" is loaded..."
  if is_kernel_module_loaded "$module"; then
    return 0
  fi

  echo "Attempting to load kernel module $module"
  if ! modprobe -v "$module";
  then
    echo "WARNING: Unable to dynamically load kernel module $module."
    echo "Attempting to build and install dkms module."
    OVS_DATAPATH_MOD_VER=$(dkms status | grep openvswitch | sed "s/:.*//" | awk -F'[,]' '{print $2}')
    echo "Datapath module version is $OVS_DATAPATH_MOD_VER"
    if ! dkms install -m openvswitch -v "$OVS_DATAPATH_MOD_VER"; then
      exit 1
    fi
    check_mod_version
    if ! modprobe -v "$module"; then
      if is_kernel_module_loaded "$module"; then
        echo "kernel module $module is loaded!"
        return 0
      else
        return 1
      fi
    else
      echo "ERROR: Failed to install kernel module $module"
      return 1
    fi
  fi
  if is_kernel_module_loaded "$module"; then
    return 0
  else
    return 1
  fi
}

check_mod_version () {
  echo "Checking installed and loaded ovs modules versions"
  
  kernel_ver=$(cat /sys/module/openvswitch/srcversion)
  mod_ver=$(modinfo /lib/modules/"$(uname -r)"/updates/dkms/openvswitch.ko |grep srcversion|awk '{print $2}')

  if [[ "$kernel_ver" == "$mod_ver" ]]; then
    OVS_VER=$(dpkg -l openvswitch-datapath-dkms | grep 'ii' | awk '{print $3}')
    echo "Module update successful, installed version: $OVS_VER"
    return 0
  else
    echo "FAIL: Module update failed. Installed and loaded openvswitch module version mismatch!"
    return 1
fi
}

load-datapath-modules () {
  # Check if openvswitch kernel modules are loaded
  if ! load_kernel_mod vport_gtp; then
    echo "FATAL: Datapath modules not loaded. Exiting..."
    exit 1
  fi

  while [ "$run" ]; do
    sleep 30
  done ;
}

start_ovs () {
  firstrun=0
  [ -f "/etc/openvswitch/conf.db" ] || firstrun=1

  echo "Checking if datapath kernel module is loaded"
  if ! is_kernel_module_loaded vport_gtp; then
    echo "FATAL: vport-gtp not loaded"
    exit 1
  fi

  if [ $firstrun -ne 0 ]; then
    # create database
    echo "Creating the ovs service database"
    ovsdb-tool create /etc/openvswitch/conf.db /usr/share/openvswitch/vswitch.ovsschema
  fi
  
  # Start openvswitch daemons
  if [ ! -d "/var/log/openvswitch" ]; then
    mkdir -p /var/log/openvswitch
  fi

  echo "Starting service openvswitch-switch"
  /usr/share/openvswitch/scripts/ovs-ctl start --system-id=random || exit 1

  # Activate bridge interfaces
  echo "Activating bridge interfaces"
  ifup --force --allow=ovs "$(ifquery --allow ovs --list)"
  ifup --force mtr0
  
  while [ "$run" ]; do
    sleep .1
  done ;
}


trap stopit SIGINT;
run=1;
stopit() {
    run=0;
};

case $1 in
    "start-ovs-only")
      start_ovs
    ;;
    "load-modules-only")
      load-datapath-modules
    ;;
    "load-modules-and-start-ovs")
      if load-datapath-modules; then
          sleep 5
          start_ovs
      fi
      echo "FATAL: Datapath modules check failed"
      exit 1
    ;;
    *) echo "Invalid option $0 [start-ovs-only|load-modules-only|load-modules-and-start-ovs]"
esac

#!/bin/sh
DEPMOD_CONFIG_DIR=/etc/depmod.d
MODULES_LOAD_CONFIG_DIR=/etc/modules-load.d
DEPMOD_CONFIG_FILE=${DEPMOD_CONFIG_DIR}/gtp.conf
MODULES_LOAD_CONFIG_FILE=${MODULES_LOAD_CONFIG_DIR}/gtp.conf

mkdir -p ${DEPMOD_CONFIG_DIR}
echo "override gtp * extra" >> ${DEPMOD_CONFIG_FILE}
echo "override gtp.ko * weak-updates" >> ${DEPMOD_CONFIG_FILE}

mkdir -p ${MODULES_LOAD_CONFIG_DIR}
echo "ip6_udp_tunnel" >> ${MODULES_LOAD_CONFIG_FILE}
echo "gtp" >> ${MODULES_LOAD_CONFIG_FILE}

echo "Running depmod..."
depmod `uname -r`

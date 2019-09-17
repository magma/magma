#!/bin/bash
cd /usr/share/yang/models || exit

yanglint --strict \
  fbc-symphony-device.yang \
  openconfig-access-points.yang \
  openconfig-ap-manager.yang \
  openconfig-extensions.yang \
  openconfig-wifi-mac.yang \
  openconfig-wifi-phy.yang \
  openconfig-wifi-types.yang \
  openconfig-interfaces.yang \
  openconfig-if-ip.yang \
  ietf-system.yang \
  /validate.json

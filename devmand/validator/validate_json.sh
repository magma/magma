#!/bin/bash

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.
#
cd /usr/share/yang/models || exit

# The order of models matters here for some strange reason so just be aware.
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
  iana-if-type.yang \
  ietf-interfaces.yang \
  ietf-system.yang \
  /validate.json

#!/usr/bin/env bash
#
# Copyright (c) 2020-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

SRC_VERSION_FN=/sys/module/openvswitch/srcversion
if test -f "$SRC_VERSION_FN"; then
    echo "Checking if upgrade is necessary"
    LOADED_VER=$(cat /sys/module/openvswitch/srcversion)
    INSTALLED_VER=$(modinfo -F srcversion openvswitch)
    echo "Loaded Version: $LOADED_VER, Installed Version:$INSTALLED_VER"
    if [ "$LOADED_VER" != "$INSTALLED_VER" ]; then
        echo "Version mismatch, Reloading openvswitch kernel module"
        /usr/share/openvswitch/scripts/ovs-ctl force-reload-kmod
    else
        echo "Skipping reload, version match"
    fi
else
    echo "No src version file found, Reloading openvswitch kernel module"
    /usr/share/openvswitch/scripts/ovs-ctl load-kmod
fi

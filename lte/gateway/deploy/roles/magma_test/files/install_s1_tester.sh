#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
#install the package from repo
sudo apt-get install s1aptester

# create symlinks to the libs
ln -sf "/usr/local/lib/s1sim/libtfw.so" "$S1AP_ROOT/bin/libtfw.so"
ln -sf "/usr/local/lib/s1sim/libtrfgen.so" "$S1AP_ROOT/bin/libtrfgen.so"
ln -sf "/usr/local/lib/s1sim/libiperf.so" "$S1AP_ROOT/bin/libiperf.so"
ln -sf "/usr/local/lib/s1sim/libiperf.so.0" "$S1AP_ROOT/bin/libiperf.so.0"

# create symlinks to the headers
ln -sf "/usr/local/include/s1sim/fw_api_int.h" "$S1AP_ROOT/bin/fw_api_int.h"
ln -sf "/usr/local/include/s1sim/fw_api_int.x" "$S1AP_ROOT/bin/fw_api_int.x"
ln -sf "/usr/local/include/s1sim/trfgen.x" "$S1AP_ROOT/bin/trfgen.x"

#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
# Building library, generator and tests

# Unset the security flag while compiling S1SIM
sed -i "s/ -DLTE_UE_NAS_SEC//" $S1AP_SRC/TestCntlrApp/build/fw.mak

for TARGET_DIR in TestCntlrApp Trfgen; do
   BUILD_DIR=$S1AP_SRC/$TARGET_DIR/build
   cd "$BUILD_DIR" || exit
   echo "building $TARGET_DIR in $BUILD_DIR"
   make -j1 "$@"
done

# Copy the libs
cp -f "$S1AP_SRC/TestCntlrApp/lib/libtfw.so" "$S1AP_ROOT/bin"
cp -f "$S1AP_SRC/Trfgen/lib/libtrfgen.so" "$S1AP_ROOT/bin"
cp -f "$S1AP_SRC/Trfgen/lib/libiperf.so" "$S1AP_ROOT/bin"
cp -f "$S1AP_SRC/Trfgen/lib/libiperf.so.0" "$S1AP_ROOT/bin"

# Copy the headers
cp -f "$S1AP_SRC/TestCntlrApp/src/tfwApp/fw_api_int.h" "$S1AP_ROOT/bin"
cp -f "$S1AP_SRC/TestCntlrApp/src/tfwApp/fw_api_int.x" "$S1AP_ROOT/bin"
cp -f "$S1AP_SRC/Trfgen/src/trfgen.x" "$S1AP_ROOT/bin"

# Copy the configs
cp -f "$MAGMA_ROOT/lte/gateway/python/integ_tests/data/s1ap_tester_cfg/"* "$S1AP_ROOT/"

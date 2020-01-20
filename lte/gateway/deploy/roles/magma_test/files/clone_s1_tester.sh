#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#

if [ -d "$S1AP_TESTER_SRC" ]; then
    cd "$S1AP_TESTER_SRC" || exit
    echo "Syncing repo"
    git pull --rebase https://github.com/facebookexperimental/S1APTester.git
else
    # Git clone
    echo "cloning repo"
    git clone https://github.com/facebookexperimental/S1APTester.git \
        "$S1AP_TESTER_SRC"

    cd "$S1AP_TESTER_SRC" || exit

    echo "set pull to rebase instead of merge"
    git config --local pull.rebase true
fi;

# Clear out the libraries to force complete rebuild after cloning
for TARGET_DIR in TestCntlrApp Trfgen; do
   BUILD_DIR=$S1AP_TESTER_SRC/$TARGET_DIR/build
   echo "cleaning up $BUILD_DIR"
   cd $BUILD_DIR
   make clean -s
done

# We made it
echo "successfully updated S1AP-Tester sources"

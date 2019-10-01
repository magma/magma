#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
#
if ! ssh-add -L | grep 'mai-s1aptester' &> /dev/null; then
   echo "ssh-agent does not hold the 'mai-s1aptester' key."
   tput setaf 1; echo "Please fix by running \`ssh-add mai-s1aptester\` on your host, and provision again."
   exit 1
fi

if [ -d "$TIP_REPO" ]; then
    cd "$TIP_REPO" || exit
    echo "Syncing repo"
    GIT_SSL_NO_VERIFY=true git pull --rebase \
       git@github.com:Telecominfraproject/TelecomInfraSWStack.git
else
    # Git clone
    echo "cloning repo"
    ssh-keyscan github.com >> ~/.ssh/known_hosts
    GIT_SSL_NO_VERIFY=true git clone \
        git@github.com:Telecominfraproject/TelecomInfraSWStack.git \
        "$TIP_REPO"

    cd "$TIP_REPO" || exit

    echo "set pull to rebase instead of merge"
    git config --local pull.rebase true
fi;

# Clear out the libraries to force complete rebuild after cloning
for TARGET_DIR in TestCntlrApp Trfgen; do
   BUILD_DIR=$S1AP_SRC/$TARGET_DIR/build
   echo "cleaning up $BUILD_DIR"
   cd $BUILD_DIR
   make clean -s
done

# We made it
echo "successfully updated S1AP-Tester sources"

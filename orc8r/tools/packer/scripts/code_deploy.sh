#!/bin/bash

# Copyright (c) 2017-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# Don't error out if dpkg lock is held by someone else
function wait_for_lock() {
    while sudo fuser /var/lib/dpkg/lock >/dev/null 2>&1 ; do
        echo "\rWaiting for other software managers to finish...\n"
        sleep 1
    done
}

# Set up Code Deploy so we can do installs
wait_for_lock
apt-get -y update
wait_for_lock
apt-get -y install awscli libssl-dev ruby
wget https://aws-codedeploy-us-west-1.s3.amazonaws.com/latest/install
chmod +x ./install
./install auto

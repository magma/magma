#!/bin/bash -eux
################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

# Don't error out if dpkg lock is held by someone else
function wait_for_lock() {
    while sudo fuser /var/lib/dpkg/lock >/dev/null 2>&1 ; do
        echo "\rWaiting for other software managers to finish...\n"
        sleep 1
    done
}

# Install Ansible repository.
wait_for_lock
sudo apt-get -y update
wait_for_lock
sudo apt-get -y install software-properties-common
wait_for_lock
sudo apt-add-repository --yes --update ppa:ansible/ansible

# Install Ansible.
wait_for_lock
sudo apt-get -y update
wait_for_lock
sudo apt-get -y install ansible

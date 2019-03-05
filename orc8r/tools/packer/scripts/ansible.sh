#!/bin/bash -eux

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
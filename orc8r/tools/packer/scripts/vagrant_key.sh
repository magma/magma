#!/bin/bash
# Install the default insecure vagrant private key
mkdir -p /home/vagrant/.ssh
wget https://raw.githubusercontent.com/hashicorp/vagrant/master/keys/vagrant.pub -O /home/vagrant/.ssh/authorized_keys
chmod 0600 /home/vagrant/.ssh/authorized_keys
chown -R vagrant /home/vagrant/.ssh

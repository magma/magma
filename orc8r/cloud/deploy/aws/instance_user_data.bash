#!/bin/bash

# Copyright (c) 2017-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# In AWS, 'user data' is the aptly named optional script that can be run when
# an instance is booted. Things that should go in here are instance-wide
# configuration, like setting up CodeDeploy, logging agents, NewRelic, etc.
# Nothing application-dependent should go in here.  At some point, we can just
# bake an AMI that contains this stuff to make the script smaller.

# Set up Code Deploy so we can do installs
apt-get -y update
apt-get -y install awscli
apt-get -y install ruby
cd /home/ubuntu
aws s3 cp s3://aws-codedeploy-eu-west-1/latest/install . --region eu-west-1
chmod +x ./install
./install auto
# Set the number of revisions kept by code deploy agent to 1
sed -i '/max_revisions/d' /etc/codedeploy-agent/conf/codedeployagent.yml
sed -i '$ a :max_revisions: 1' /etc/codedeploy-agent/conf/codedeployagent.yml
sudo service codedeploy-agent restart
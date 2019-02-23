#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Install librdkafka and some other preinstall dependencies
wget -qO - https://packages.confluent.io/deb/5.1/archive.key | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://packages.confluent.io/deb/5.1 stable main"
sudo apt-get update -qq
sudo apt-get install -y librdkafka-dev librdkafka1 bzr parallel build-essential unzip default-jre

# Install protobuf compiler
sudo curl -Lfs https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip
sudo unzip protoc3.zip -d protoc3
sudo mv protoc3/bin/protoc /bin/protoc
sudo chmod a+rx /bin/protoc
sudo mv protoc3/include/google /usr/include/
sudo chmod -R a+Xr /usr/include/google
sudo rm -rf protoc3.zip protoc3

# chown /var/tmp to travis user (Makefile uses this dir)
sudo chown -R travis /var/tmp

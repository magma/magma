#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Install some preinstall dependencies
sudo apt-get update -qq
sudo apt-get install -y bzr parallel build-essential unzip default-jre

# Install protobuf compiler
sudo curl -Lfs https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip
sudo unzip protoc3.zip -d protoc3
sudo mv protoc3/bin/protoc /bin/protoc
sudo chmod a+rx /bin/protoc
sudo mv protoc3/include/google /usr/include/
sudo chmod -R a+Xr /usr/include/google
sudo rm -rf protoc3.zip protoc3

# chown /var/tmp to travis user (Makefile uses this dir)
# sudo chown -R travis /var/tmp

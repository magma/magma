#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Install and start dynamoDB
sudo mkdir -p /var/tmp/archives
sudo wget -O /var/tmp/archives/dynamo.zip https://s3-us-west-2.amazonaws.com/dynamodb-local/dynamodb_local_latest.zip
sudo unzip /var/tmp/archives/dynamo.zip -d /var/tmp/archives/dynamo
sudo java -Djava.library.path=/var/tmp/archives/dynamo/DynamoDBLocal_lib -jar /var/tmp/archives/dynamo/DynamoDBLocal.jar -dbPath /var/tmp/archives/dynamo -sharedDb &

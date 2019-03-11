#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

sudo cp orc8r/cloud/deploy/roles/controller/files/magma.service /etc/systemd/system/magma@.service
sudo cp orc8r/cloud/deploy/roles/controller/files/magma_certifier.service /etc/systemd/system/magma@certifier.service
sudo cp orc8r/cloud/deploy/roles/controller/files/magma_bootstrapper.service /etc/systemd/system/magma@bootstrapper.service
sudo cp orc8r/cloud/deploy/roles/controller/files/magma_obsidian.service /etc/systemd/system/magma@obsidian.service
sudo cp orc8r/cloud/deploy/roles/controller/files/magma_metricsd.service /etc/systemd/system/magma@metricsd.service
sudo systemctl daemon-reload

mkdir -p ${MAGMA_ROOT}/plugins
mkdir -p ${MAGMA_ROOT}/.cache/test_certs

sudo mkdir -p /var/opt/magma/
sudo ln -s ${MAGMA_ROOT}/plugins /var/opt/magma/plugins
sudo ln -s ${GOPATH}/bin /var/opt/magma/bin
sudo ln -s ${MAGMA_ROOT}/orc8r/cloud/deploy/files/envdir /var/opt/magma/envdir
sudo ln -s ${MAGMA_ROOT}/.cache/test_certs /var/opt/magma/certs

orc8r/cloud/scripts/create_test_certs ${MAGMA_ROOT}/.cache/test_certs
orc8r/cloud/scripts/create_test_vpn_certs ${MAGMA_ROOT}/.cache/test_certs

sudo pip install envdir
sudo ln -s /usr/local/bin/envdir /usr/bin/envdir

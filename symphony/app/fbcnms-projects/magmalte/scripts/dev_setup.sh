#!/usr/bin/env bash
#
# Copyright 2004-present Facebook. All Rights Reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

set -e

# Add the test admin user
docker-compose exec magmalte yarn migrate
docker-compose exec magmalte yarn setAdminPassword fb-test admin@magma.test password1234
docker-compose exec magmalte yarn setAdminPassword master admin@magma.test password1234

# Docker run in a Linux host doesn't resolve host.docker.internal to the host IP.
# See https://github.com/docker/for-linux/issues/264
# Add an entry for the host. This is a no-op for Mac.
docker-compose exec magmalte /bin/sh -c "ip -4 route list match 0/0 | awk '{print \$3 \" host.docker.internal\"}' >> /etc/hosts"

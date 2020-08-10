#!/usr/bin/env bash
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e

# Add the test admin user
docker-compose exec magmalte yarn migrate
docker-compose exec magmalte yarn setAdminPassword magma-test admin@magma.test password1234
docker-compose exec magmalte yarn setAdminPassword master admin@magma.test password1234

# Docker run in a Linux host doesn't resolve host.docker.internal to the host IP.
# See https://github.com/docker/for-linux/issues/264
# Add an entry for the host. This is a no-op for Mac.
docker-compose exec magmalte /bin/sh -c "ip -4 route list match 0/0 | awk '{print \$3 \" host.docker.internal\"}' >> /etc/hosts"

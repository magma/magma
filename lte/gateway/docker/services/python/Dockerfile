################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

# -----------------------------------------------------------------------------
# Builder image for Python binaries and Magma proto files
# -----------------------------------------------------------------------------
ARG CPU_ARCH=x86_64
ARG DEB_PORT=amd64
ARG OS_DIST=ubuntu
ARG OS_RELEASE=focal

FROM $OS_DIST:$OS_RELEASE AS builder
ARG CPU_ARCH
ARG DEB_PORT
ARG OS_DIST
ARG OS_RELEASE

ENV MAGMA_DEV_MODE 0
ENV TZ=Europe/Paris
ENV MAGMA_ROOT=/magma
ENV PIP_CACHE_HOME="~/.pipcache"
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
  docker.io \
  git \
  lsb-release \
  libsystemd-dev \
  pkg-config \
  python3-dev \
  python3-pip \
  sudo \
  wget \
  && rm -rf /var/lib/apt/lists/*

# Download Bazel
RUN  wget -P /usr/sbin --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-"${DEB_PORT}" && \
  chmod +x /usr/sbin/bazelisk-linux-"${DEB_PORT}" && \
  ln -s /usr/sbin/bazelisk-linux-"${DEB_PORT}" /usr/sbin/bazel

COPY ./lte/gateway/python/integ_tests/s1aptests/ovs $MAGMA_ROOT/lte/gateway/python/integ_tests/s1aptests/ovs
COPY ./lte/gateway/python/load_tests $MAGMA_ROOT/lte/gateway/python/load_tests
COPY ./lte/gateway/python/scripts $MAGMA_ROOT/lte/gateway/python/scripts
COPY ./lte/gateway/python/dhcp_helper_cli $MAGMA_ROOT/lte/gateway/python/dhcp_helper_cli
COPY ./lte/gateway/python/magma $MAGMA_ROOT/lte/gateway/python/magma
COPY ./orc8r/gateway/python $MAGMA_ROOT/orc8r/gateway/python

COPY ./lte/gateway/configs $MAGMA_ROOT/lte/gateway/configs
COPY ./orc8r/gateway/configs/templates $MAGMA_ROOT/orc8r/gateway/configs/templates

COPY ./protos $MAGMA_ROOT/protos
COPY ./cwf/protos $MAGMA_ROOT/cwf/protos
COPY ./lte/protos $MAGMA_ROOT/lte/protos
COPY ./orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY ./feg/protos $MAGMA_ROOT/feg/protos
COPY ./dp/protos $MAGMA_ROOT/dp/protos

COPY ./cwf/swagger $MAGMA_ROOT/cwf/swagger
COPY ./feg/swagger $MAGMA_ROOT/feg/swagger
COPY ./lte/swagger $MAGMA_ROOT/lte/swagger
COPY ./orc8r/swagger $MAGMA_ROOT/orc8r/swagger

COPY ./lte/gateway/deploy/roles/magma/files/patches $MAGMA_ROOT/lte/gateway/deploy/roles/magma/files/patches

# Copy Bazel files
COPY WORKSPACE.bazel BUILD.bazel .bazelignore .bazelrc .bazelversion ${MAGMA_ROOT}/
COPY bazel/ ${MAGMA_ROOT}/bazel
COPY ./lte/gateway/python/BUILD.bazel $MAGMA_ROOT/lte/gateway/python/BUILD.bazel
COPY ./lte/gateway/release/BUILD.bazel $MAGMA_ROOT/lte/gateway/release/BUILD.bazel
COPY ./lte/gateway/release/deb_dependencies.bzl $MAGMA_ROOT/lte/gateway/release/deb_dependencies.bzl

WORKDIR /magma
RUN bazel build //lte/gateway/release:python_executables_tar //lte/gateway/release:dhcp_helper_cli_tar

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM $OS_DIST:$OS_RELEASE AS gateway_python
ARG OS_DIST
ARG CPU_ARCH
ARG OS_RELEASE
ARG EXTRA_REPO

ENV VIRTUAL_ENV=/build
ENV TZ=Europe/Paris

ARG JSONPOINTER_VERSION=1.13
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
  apt-transport-https \
  ca-certificates \
  docker.io \
  ethtool \
  inetutils-ping \
  iproute2 \
  iptables \
  isc-dhcp-client \
  linux-headers-generic \
  net-tools \
  netcat \
  nghttp2-proxy \
  python3-pip \
  python3-venv \
  redis-server \
  sudo \
  wget \
  && rm -rf /var/lib/apt/lists/*

ENV PATH="/magma/orc8r/gateway/python/scripts/:/magma/lte/gateway/python/scripts/:$PATH"

# Add the magma apt repo
COPY keys/linux_foundation_registry_key.asc /etc/apt/trusted.gpg.d/magma.asc
RUN echo "deb https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main" > /etc/apt/sources.list.d/magma.list

RUN echo "deb https://packages.fluentbit.io/ubuntu/focal focal main" > /etc/apt/sources.list.d/tda.list
RUN wget -qO - https://packages.fluentbit.io/fluentbit.key | apt-key add -

RUN apt-get update && apt-get install -y \
  bcc-tools \
  libopenvswitch \
  openvswitch-datapath-dkms \
  openvswitch-common \
  openvswitch-switch \
  td-agent-bit=1.7.8 \
  wireguard \
  && rm -rf /var/lib/apt/lists/*

COPY ./lte/gateway/python/dhcp_helper_cli/dhcp_helper_cli.py /usr/local/bin/
RUN mkdir -p /var/opt/magma/

# Copy the build artifacts.
COPY --from=builder /magma/bazel-bin/lte/gateway/release/magma_python_executables.tar.gz \
                    /tmp/magma_python_executables.tar.gz
RUN tar -xf /tmp/magma_python_executables.tar.gz --directory / && \
    rm /tmp/magma_python_executables.tar.gz

COPY --from=builder /magma/bazel-bin/lte/gateway/release/magma-dhcp-cli_1.9.0-VERSION-SUFFIX_amd64.tar.gz \
                    /tmp/magma-dhcp-cli_1.9.0-VERSION-SUFFIX_amd64.tar.gz
RUN tar -xf /tmp/magma-dhcp-cli_1.9.0-VERSION-SUFFIX_amd64.tar.gz --directory / && \
    rm /tmp/magma-dhcp-cli_1.9.0-VERSION-SUFFIX_amd64.tar.gz

WORKDIR /usr/local/bin

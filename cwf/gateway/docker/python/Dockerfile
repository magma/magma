# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# -----------------------------------------------------------------------------
# Builder image
# -----------------------------------------------------------------------------
FROM ubuntu:focal AS builder

ARG DEB_PORT=amd64
# workaround to avoid interactive tzdata configurtaion
ARG DEBIAN_FRONTEND=noninteractive

# Install the build deps with apt.
RUN apt-get -y update && apt-get -y install \
    apt-transport-https \
    apt-utils \
    git \
    libsystemd-dev \
    python3.8 \
    python3-distutils \
    python3-pip \
    pkg-config \
    systemd \
    wget

ENV MAGMA_ROOT /magma
ENV PIP_CACHE_HOME ~/.pipcache

# Download Bazel
RUN  wget -P /usr/sbin --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-"${DEB_PORT}" && \
  chmod +x /usr/sbin/bazelisk-linux-"${DEB_PORT}" && \
  ln -s /usr/sbin/bazelisk-linux-"${DEB_PORT}" /usr/sbin/bazel

# Generate python proto bindings.
COPY cwf/protos $MAGMA_ROOT/cwf/protos
COPY cwf/swagger $MAGMA_ROOT/cwf/swagger
COPY feg/protos $MAGMA_ROOT/feg/protos
COPY feg/swagger $MAGMA_ROOT/feg/swagger
COPY lte/gateway/configs $MAGMA_ROOT/lte/gateway/configs
COPY lte/gateway/python $MAGMA_ROOT/lte/gateway/python
COPY lte/gateway/deploy/roles/magma/files/patches $MAGMA_ROOT/lte/gateway/deploy/roles/magma/files/patches
COPY lte/protos $MAGMA_ROOT/lte/protos
COPY lte/swagger $MAGMA_ROOT/lte/swagger
COPY orc8r/gateway/configs $MAGMA_ROOT/orc8r/gateway/configs
COPY orc8r/gateway/python $MAGMA_ROOT/orc8r/gateway/python
COPY orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY orc8r/swagger $MAGMA_ROOT/orc8r/swagger
COPY protos $MAGMA_ROOT/protos

# Copy Bazel files
COPY WORKSPACE.bazel BUILD.bazel .bazelignore .bazelrc .bazelversion ${MAGMA_ROOT}/
COPY bazel/ ${MAGMA_ROOT}/bazel
COPY ./lte/gateway/release/BUILD.bazel ${MAGMA_ROOT}/lte/gateway/release/BUILD.bazel
COPY ./lte/gateway/release/deb_dependencies.bzl ${MAGMA_ROOT}/lte/gateway/release/deb_dependencies.bzl

WORKDIR /magma
RUN bazel build //lte/gateway/release:cwf_python_executables_tar

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM ubuntu:focal AS lte_gateway_python

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https
COPY cwf/gateway/deploy/roles/ovs/files/magma-preferences /etc/apt/preferences.d/

# Add the magma apt repo
COPY keys/linux_foundation_registry_key.asc /etc/apt/trusted.gpg.d/magma.asc

#[TODO_GPG_KEY_WORKAROUND]Temporary steps to be removed once GPG public key issue has been fixed.
RUN echo "deb [trusted=yes] https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main" \
    > /etc/apt/sources.list.d/magma.list
RUN add-apt-repository "deb http://archive.ubuntu.com/ubuntu/ focal-proposed restricted main multiverse universe"

#[TODO_GPG_KEY_WORKAROUND]Temporary steps to be removed once GPG public key issue has been fixed.
RUN echo "Acquire::AllowInsecureRepositories true;" > /etc/apt/apt.conf.d/99AllowInsecureRepositories && \
    echo "APT::Get::AllowUnauthenticated true;" >> /etc/apt/apt.conf.d/99AllowInsecureRepositories

RUN apt-get -y update && apt-get -y install \
    curl \
    libc-ares2 \
    libev4 \
    libffi-dev \
    libjansson4 \
    libjemalloc2 \
    libssl-dev \
    libsystemd-dev \
    magma-nghttpx=1.31.1-1 \
    net-tools \
    openssl \
    pkg-config \
    python-cffi \
    python3-pip \
    python3.8 \
    python3.8-dev \
    redis-server \
    iptables \
    git \
    automake \
    gcc \
    libtool \
    libcap-ng-dev \
    linux-headers-generic \
    netcat \
    iputils-ping \
    bcc-tools

# Copy the build artifacts.
COPY --from=builder /magma/bazel-bin/lte/gateway/release/cwf_python_executables.tar.gz \
                    /tmp/cwf_python_executables.tar.gz
RUN tar -xf /tmp/cwf_python_executables.tar.gz --directory / && \
    rm /tmp/cwf_python_executables.tar.gz

# Temporary workaround to restore uplink bridge flows
RUN mkdir -p /var/opt/magma/scripts
COPY cwf/gateway/deploy/roles/cwag/files/add_uplink_bridge_flows.sh /var/opt/magma/scripts

# Install OVS via Magma bionic pkg repo
RUN git clone --depth 1 --single-branch --branch v2.12.0 https://github.com/openvswitch/ovs.git

COPY cwf/gateway/deploy/roles/ovs/files/0001-Add-custom-IPDR-fields-for-IPFIX-export.patch /tmp
COPY cwf/gateway/deploy/roles/ovs/files/0002-ovs-Handle-spaces-in-ovs-arguments.patch /tmp
COPY cwf/gateway/deploy/roles/ovs/files/0003-Add-pdp_start_epoch-custom-field-to-IPFIX-export.patch /tmp
COPY cwf/gateway/deploy/roles/ovs/files//0004-ovsdb-idlc.in-dict-changes.patch /tmp
WORKDIR ovs
RUN git apply /tmp/0001-Add-custom-IPDR-fields-for-IPFIX-export.patch
RUN git apply /tmp/0002-ovs-Handle-spaces-in-ovs-arguments.patch
RUN git apply /tmp/0003-Add-pdp_start_epoch-custom-field-to-IPFIX-export.patch
RUN git apply /tmp/0004-ovsdb-idlc.in-dict-changes.patch
RUN ./boot.sh
RUN ./configure --prefix=/usr --localstatedir=/var --sysconfdir=/etc
RUN make
RUN make install

# Copy the configs.
COPY lte/gateway/configs /etc/magma
COPY orc8r/gateway/configs/templates /etc/magma/templates
RUN mkdir -p /var/opt/magma/configs

WORKDIR /

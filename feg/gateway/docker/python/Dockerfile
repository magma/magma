################################################################################
# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

FROM ubuntu:focal AS builder

ARG DEB_PORT=amd64
# workaround to avoid interactive tzdata configuration
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
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64/

# Download Bazel
RUN  wget -P /usr/sbin --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-"${DEB_PORT}" && \
  chmod +x /usr/sbin/bazelisk-linux-"${DEB_PORT}" && \
  ln -s /usr/sbin/bazelisk-linux-"${DEB_PORT}" /usr/sbin/bazel

COPY cwf/protos $MAGMA_ROOT/cwf/protos
COPY cwf/swagger $MAGMA_ROOT/cwf/swagger
COPY feg/protos $MAGMA_ROOT/feg/protos
COPY feg/swagger $MAGMA_ROOT/feg/swagger
COPY lte/gateway/configs $MAGMA_ROOT/lte/gateway/configs
COPY lte/protos $MAGMA_ROOT/lte/protos
COPY lte/swagger $MAGMA_ROOT/lte/swagger
COPY orc8r/gateway/python $MAGMA_ROOT/orc8r/gateway/python
COPY orc8r/gateway/configs $MAGMA_ROOT/orc8r/gateway/configs
COPY orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY orc8r/swagger $MAGMA_ROOT/orc8r/swagger
COPY orc8r/tools/ansible/roles/fluent_bit/files $MAGMA_ROOT/orc8r/tools/ansible/roles/fluent_bit/files
COPY protos $MAGMA_ROOT/protos

# Copy Bazel files
COPY WORKSPACE.bazel BUILD.bazel .bazelignore .bazelrc .bazelversion ${MAGMA_ROOT}/
COPY bazel/ ${MAGMA_ROOT}/bazel
COPY ./orc8r/gateway/python/BUILD.bazel ${MAGMA_ROOT}/orc8r/gateway/python/BUILD.bazel
COPY ./lte/gateway/release/BUILD.bazel ${MAGMA_ROOT}/lte/gateway/release/BUILD.bazel
COPY ./lte/gateway/release/deb_dependencies.bzl ${MAGMA_ROOT}/lte/gateway/release/deb_dependencies.bzl
COPY ./lte/gateway/configs/templates/BUILD.bazel ${MAGMA_ROOT}/lte/gateway/configs/templates/BUILD.bazel
COPY ./orc8r/gateway/configs/templates/BUILD.bazel ${MAGMA_ROOT}/orc8r/gateway/configs/templates/BUILD.bazel

WORKDIR /magma
RUN bazel build //lte/gateway/release:feg_python_executables_tar

# -----------------------------------------------------------------------------
# Production image
# -----------------------------------------------------------------------------
FROM ubuntu:focal AS gateway_python
ARG MAGMA_BUILD_BRANCH=unknown
ARG MAGMA_BUILD_TAG=unknown
ARG MAGMA_BUILD_COMMIT_HASH=unknonw
ARG MAGMA_BUILD_COMMIT_DATE=unknown

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https curl

COPY keys/linux_foundation_registry_key.asc /etc/apt/trusted.gpg.d/magma.asc

#[TODO_GPG_KEY_WORKAROUND]Temporary steps to be removed once GPG public key issue has been fixed.
#RUN add-apt-repository "deb https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main"
RUN echo "deb [trusted=yes] https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main" \
    > /etc/apt/sources.list.d/magma.list

#[TODO_GPG_KEY_WORKAROUND]Temporary steps to be removed once GPG public key issue has been fixed.
RUN echo "Acquire::AllowInsecureRepositories true;" > /etc/apt/apt.conf.d/99AllowInsecureRepositories && \
    echo "APT::Get::AllowUnauthenticated true;" >> /etc/apt/apt.conf.d/99AllowInsecureRepositories

RUN curl -L http://packages.fluentbit.io/fluentbit.key > /tmp/fluentbit.key
RUN apt-key add /tmp/fluentbit.key 
#[TODO_GPG_KEY_WORKAROUND]Temporary steps to be removed once GPG public key issue has been fixed.
#&& add-apt-repository "deb https://packages.fluentbit.io/ubuntu/focal focal main"
RUN echo "deb [trusted=yes] https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main" \
    > /etc/apt/sources.list.d/magma.list

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install \
    iproute2 \
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
    iputils-ping \
    pkg-config \
    python-cffi \
    python3-pip \
    python3.8 \
    python3.8-dev \
    redis-server \
    git \
    netcat \
    td-agent-bit

# Install docker.
RUN curl -sSL https://get.docker.com/ > /tmp/get_docker.sh && \
    sh /tmp/get_docker.sh && \
    rm /tmp/get_docker.sh

# Copy the build artifacts.
COPY --from=builder /magma/bazel-bin/lte/gateway/release/feg_python_executables.tar.gz \
                    /tmp/feg_python_executables.tar.gz
RUN tar -xf /tmp/feg_python_executables.tar.gz --directory / && \
    rm /tmp/feg_python_executables.tar.gz

WORKDIR /usr/local/bin

# Copy the configs to overwrite the LTE configs.
# WARNING: With the current setup this needs to happen
# after the installation of the Bazel build artifacts!
COPY feg/gateway/configs /etc/magma

COPY orc8r/gateway/configs/templates /etc/magma/templates

RUN mkdir -p /var/opt/magma/configs
RUN mkdir -p /var/opt/magma/fluent-bit

# Add commit information
ENV MAGMA_BUILD_BRANCH $MAGMA_BUILD_BRANCH
ENV MAGMA_BUILD_TAG $MAGMA_BUILD_TAG
ENV MAGMA_BUILD_COMMIT_HASH $MAGMA_BUILD_COMMIT_HASH
ENV MAGMA_BUILD_COMMIT_DATE $MAGMA_BUILD_COMMIT_DATE

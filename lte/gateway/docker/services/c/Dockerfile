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
# Builder image for C binaries and Magma proto files
# -----------------------------------------------------------------------------
ARG CPU_ARCH=x86_64
ARG DEB_PORT=amd64
ARG OS_DIST=ubuntu
ARG OS_RELEASE=focal
ARG EXTRA_REPO=https://linuxfoundation.jfrog.io/artifactory/magma-packages-test
ARG CLANG_VERSION=3.8
ARG FEATURES=mme_oai

FROM $OS_DIST:$OS_RELEASE AS builder
ARG CPU_ARCH
ARG DEB_PORT
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO
ARG CLANG_VERSION

ENV MAGMA_ROOT /magma
ENV C_BUILD /build/c
ENV OAI_BUILD $C_BUILD/oai
ENV TZ=Europe/Paris

ENV CCACHE_DIR ${MAGMA_ROOT}/.cache/gateway/ccache
ENV MAGMA_DEV_MODE 0
ENV XDG_CACHE_HOME ${MAGMA_ROOT}/.cache

# [TODO_GPG_KEY_WORKAROUND] Temporary steps to be removed once GPG public key issue has been fixed.
RUN echo "Acquire::AllowInsecureRepositories true;" > /etc/apt/apt.conf.d/99AllowInsecureRepositories \
    && echo "APT::Get::AllowUnauthenticated true;" >> /etc/apt/apt.conf.d/99AllowInsecureRepositories

RUN apt-get update && \
  # Setup necessary tools for adding the Magma repository
  apt-get install -y apt-utils software-properties-common apt-transport-https gnupg wget && \
  # Download Bazel
  wget -P /usr/sbin --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-"${DEB_PORT}" && \
  chmod +x /usr/sbin/bazelisk-linux-"${DEB_PORT}" && \
  ln -s /usr/sbin/bazelisk-linux-"${DEB_PORT}" /usr/sbin/bazel

# Install dependencies required for building
RUN apt-get update && apt-get install -y \
  autoconf \
  autogen \
  build-essential \
  ccache \
  check \
  cmake \
  curl \
  git \
  libboost-chrono-dev \
  libboost-context-dev \
  libboost-program-options-dev \
  libboost-filesystem-dev \
  libboost-regex-dev \
  libc++-dev \
  libconfig-dev \
  libcurl4-openssl-dev \
  libczmq-dev \
  libdouble-conversion-dev \
  libgflags-dev \
  libgmp3-dev \
  libgoogle-glog-dev \
  libmnl-dev \
  libpcap-dev \
  libprotoc-dev \
  libsctp-dev \
  libsqlite3-dev \
  libssl-dev \
  libtspi-dev \
  libtool \
  libxml2-dev \
  libxslt-dev \
  libyaml-cpp-dev \
  ninja-build \
  nlohmann-json3-dev \
  pkg-config \
  protobuf-compiler \
  python3-pip \
  sudo \
  unzip \
  uuid-dev \
  && rm -rf /var/lib/apt/lists/*

# Add the magma apt repo
COPY keys/linux_foundation_registry_key.asc /etc/apt/trusted.gpg.d/magma.asc
#[TODO_GPG_KEY_WORKAROUND]Temporary steps to be removed once GPG public key issue has been fixed.
#RUN echo "deb https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main" > /etc/apt/sources.list.d/magma.list
RUN echo "deb [trusted=yes] https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main" \
    > /etc/apt/sources.list.d/magma.list
RUN apt-get update && apt-get install -y \
  grpc-dev \
  libfolly-dev \
  liblfds710 \
  magma-cpp-redis \
  magma-libfluid \
  oai-asn1c \
  oai-freediameter \
  oai-gnutls \
  oai-nettle \
  prometheus-cpp-dev \
  && rm -rf /var/lib/apt/lists/* \
  && rm /etc/apt/sources.list.d/magma.list

ENV MAGMA_ROOT /magma
WORKDIR /magma

# Copy Bazel files at root and third_party
COPY WORKSPACE.bazel BUILD.bazel .bazelignore .bazelrc .bazelversion ${MAGMA_ROOT}/
COPY bazel/ ${MAGMA_ROOT}/bazel
COPY third_party/build/patches/libfluid/ ${MAGMA_ROOT}/third_party/build/patches/libfluid/

# Build external dependencies first. This will help not rebuilt all dependencies triggered by Magma changes.
RUN bazel build \
  @com_github_grpc_grpc//:grpc++ \
  @com_google_protobuf//:protobuf \
  @prometheus_cpp//:prometheus-cpp \
  @yaml-cpp//:yaml-cpp \
  @github_nlohmann_json//:json \
  @sentry_native//:sentry

# Copy proto files
COPY feg/protos ${MAGMA_ROOT}/feg/protos
COPY feg/gateway/services/aaa/protos ${MAGMA_ROOT}/feg/gateway/services/aaa/protos
COPY lte/protos ${MAGMA_ROOT}/lte/protos
COPY orc8r/protos ${MAGMA_ROOT}/orc8r/protos
COPY protos ${MAGMA_ROOT}/protos

# Build session_manager c code
COPY orc8r/gateway/c/common ${MAGMA_ROOT}/orc8r/gateway/c/common
COPY lte/gateway/c ${MAGMA_ROOT}/lte/gateway/c

COPY lte/gateway/python/scripts ${MAGMA_ROOT}/lte/gateway/python/scripts
COPY lte/gateway/docker ${MAGMA_ROOT}/lte/gateway/docker
COPY lte/gateway/docker/mme/configs/ ${MAGMA_ROOT}/lte/gateway/docker/configs/

RUN bazel build  \
  --config=production \
  //lte/gateway/c/sctpd/src:sctpd \
  //lte/gateway/c/connection_tracker/src:connectiond \
  //lte/gateway/c/li_agent/src:liagentd \
  //lte/gateway/c/session_manager:sessiond \
  //lte/gateway/c/core:agw_of

# Prepare config file
COPY lte/gateway/configs ${MAGMA_ROOT}/lte/gateway/configs

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM $OS_DIST:$OS_RELEASE AS gateway_c
ARG CPU_ARCH
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO
ARG MAGMA_VERSION=master

ENV MAGMA_ROOT /magma
ENV C_BUILD /build/c
ENV TZ=Europe/Paris
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# [TODO_GPG_KEY_WORKAROUND] Temporary steps to be removed once GPG public key issue has been fixed.
RUN echo "Acquire::AllowInsecureRepositories true;" > /etc/apt/apt.conf.d/99AllowInsecureRepositories \
    && echo "APT::Get::AllowUnauthenticated true;" >> /etc/apt/apt.conf.d/99AllowInsecureRepositories

# Install runtime dependencies
RUN apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  apt-transport-https \
  apt-utils \
  ca-certificates \
  gnupg \
  iproute2 \
  iptables \
  libgoogle-glog-dev \
  libidn11-dev \
  libmnl-dev \
  libprotoc-dev \
  libsctp-dev \
  libtspi1 \
  libyaml-cpp-dev \
  net-tools \
  netcat \
  openssl \
  psmisc \
  sudo \
  tshark \
  tzdata \
  wget \
  && rm -rf /var/lib/apt/lists/*

COPY keys/linux_foundation_registry_key.asc /etc/apt/trusted.gpg.d/magma.asc
#[TODO_GPG_KEY_WORKAROUND]Temporary steps to be removed once GPG public key issue has been fixed.
#RUN echo "deb https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main" > /etc/apt/sources.list.d/magma.list
RUN echo "deb [trusted=yes] https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main" \
    > /etc/apt/sources.list.d/magma.list

RUN apt-get update && apt-get install -y \
  libopenvswitch \
  openvswitch-common \
  openvswitch-datapath-dkms \
  openvswitch-switch \
  && rm -rf /var/lib/apt/lists/*

# Copy pre-built shared object files
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/liblsan.so.0 /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libprotobuf.so /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libboost* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libevent-* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libsnappy.* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libdouble-conversion.* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libicui18n.* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libicuuc.* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libicudata.* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libconfig.* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libczmq.* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libevent* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libzmq* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libsodium* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libpgm* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libnorm* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libgflags* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libgssapi_krb5* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libkrb5* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libk5crypto* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libkeyutils* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libcurl.so.* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/librtmp.so.* /usr/lib/"${CPU_ARCH}"-linux-gnu/
COPY --from=builder /usr/lib/"${CPU_ARCH}"-linux-gnu/libssh.so.* /usr/lib/"${CPU_ARCH}"-linux-gnu/


COPY --from=builder /usr/local/lib/liblfds710.so /usr/local/lib/
COPY --from=builder /usr/local/lib/libgrpc++.so /usr/local/lib/
COPY --from=builder /usr/local/lib/libfolly.so /usr/local/lib/
COPY --from=builder /usr/local/lib/libgrpc.so /usr/local/lib/
COPY --from=builder /usr/local/lib/libgpr.so /usr/local/lib/
COPY --from=builder /usr/lib/libnettle.so  /usr/local/lib/
COPY --from=builder /usr/lib/libfluid* /usr/local/lib/
COPY --from=builder /usr/lib/libgnutls.so.* /usr/local/lib/
COPY --from=builder /usr/lib/libhogweed.so.2 /usr/local/lib/

COPY --from=builder /usr/local/lib/libfdproto.so.6 /usr/local/lib/
COPY --from=builder /usr/local/lib/libfdcore.so.6 /usr/local/lib/
COPY --from=builder /usr/local/lib/libaddress_sorting.so /usr/local/lib/

# Copy the build artifacts.
COPY --from=builder ${MAGMA_ROOT}/bazel-bin/lte/gateway/c/session_manager/sessiond /usr/local/bin/sessiond
COPY --from=builder ${MAGMA_ROOT}/bazel-bin/lte/gateway/c/sctpd/src/sctpd /usr/local/bin/sctpd
COPY --from=builder ${MAGMA_ROOT}/bazel-bin/lte/gateway/c/connection_tracker/src/connectiond /usr/local/bin/connectiond
COPY --from=builder ${MAGMA_ROOT}/bazel-bin/lte/gateway/c/li_agent/src/liagentd /usr/local/bin/liagentd
COPY --from=builder ${MAGMA_ROOT}/bazel-bin/lte/gateway/c/core/agw_of /usr/local/bin/oai_mme

RUN ldconfig 2> /dev/null

# Copy the configs.
COPY lte/gateway/configs /etc/magma
COPY orc8r/gateway/configs/templates /etc/magma/templates
COPY lte/gateway/deploy/roles/magma/files/magma-create-gtp-port.sh /usr/local/bin/

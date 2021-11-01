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
ARG OS_DIST=ubuntu
ARG OS_RELEASE=focal
ARG EXTRA_REPO=https://facebookconnectivity.jfrog.io/artifactory/list/dev-focal/
ARG CLANG_VERSION=3.8
ARG FEATURES=mme_oai

FROM $OS_DIST:$OS_RELEASE AS builder
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO
ARG CLANG_VERSION

ENV MAGMA_ROOT /magma
ENV C_BUILD /build/c
ENV OAI_BUILD $C_BUILD/oai
ENV TZ=Europe/Paris

ENV CCACHE_DIR $MAGMA_ROOT/.cache/gateway/ccache
ENV MAGMA_DEV_MODE 0
ENV XDG_CACHE_HOME $MAGMA_ROOT/.cache

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https gnupg

COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb ${EXTRA_REPO} ${OS_RELEASE} main"

# Install dependencies required for building
RUN apt-get update && apt-get install -y \
  sudo \
  curl \
  wget \
  unzip \
  cmake \
  git \
  build-essential \
  autoconf \
  libtool \
  pkg-config \
  libgflags-dev \
  libgtest-dev \
  clang-11 \
  #clang-format-11 \
  libc++-dev \
  protobuf-compiler \
  grpc-dev \
  ninja-build \
  autogen \
  ccache \
  libprotoc-dev \
  libxml2-dev \
  libxslt-dev \
  libyaml-cpp-dev \
  nlohmann-json3-dev \
  magma-cpp-redis \
  libgoogle-glog-dev \
  prometheus-cpp-dev \
  libfolly-dev \
  libsctp-dev \
  libpcap-dev \
  libtins-dev \
  libmnl-dev \
  uuid-dev \
  libcurl4-openssl-dev \
  magma-libfluid \
  libdouble-conversion-dev \
  libboost-chrono-dev \
  libboost-context-dev \
  libboost-program-options-dev \
  libboost-filesystem-dev \
  libboost-regex-dev \
  check \
  libgtest-dev \
  liblfds710 \
  google-mock \
  libssl-dev \
  libsctp-dev \
  libtspi-dev \
  oai-asn1c \
  oai-nettle \
  oai-gnutls \
  oai-freediameter \
  prometheus-cpp-dev \
  libconfig-dev \
  libgmp3-dev \
  libczmq-dev

# Copy proto files
COPY feg/protos $MAGMA_ROOT/feg/protos
COPY feg/gateway/services/aaa/protos $MAGMA_ROOT/feg/gateway/services/aaa/protos
COPY lte/protos $MAGMA_ROOT/lte/protos
COPY orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY protos $MAGMA_ROOT/protos

# Build session_manager c code
COPY lte/gateway/Makefile $MAGMA_ROOT/lte/gateway/Makefile
COPY orc8r/gateway/c/common $MAGMA_ROOT/orc8r/gateway/c/common
COPY lte/gateway/c $MAGMA_ROOT/lte/gateway/c

COPY lte/gateway/python/scripts $MAGMA_ROOT/lte/gateway/python/scripts
COPY lte/gateway/docker $MAGMA_ROOT/lte/gateway/docker
COPY lte/gateway/docker/mme/configs/ $MAGMA_ROOT/lte/gateway/docker/configs/

ARG BUILD_TYPE=RelWithDebInfo
ENV BUILD_TYPE=$BUILD_TYPE
RUN make -C $MAGMA_ROOT/lte/gateway/ build_session_manager BUILD_TYPE="${BUILD_TYPE}"
RUN make -C $MAGMA_ROOT/lte/gateway/ build_sctpd BUILD_TYPE="${BUILD_TYPE}"
RUN make -C $MAGMA_ROOT/lte/gateway/ build_connection_tracker BUILD_TYPE="${BUILD_TYPE}"
RUN make -C $MAGMA_ROOT/lte/gateway/ build_li_agent BUILD_TYPE="${BUILD_TYPE}"
RUN make -C $MAGMA_ROOT/lte/gateway/ build_oai BUILD_TYPE="${BUILD_TYPE}"

# Prepare config file
COPY lte/gateway/configs $MAGMA_ROOT/lte/gateway/configs

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM $OS_DIST:$OS_RELEASE AS gateway_c
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO
ARG MAGMA_VERSION=master

ENV MAGMA_ROOT /magma
ENV C_BUILD /build/c
ENV TZ=Europe/Paris

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
  apt-utils \
  apt-transport-https \
  ca-certificates \
  wget \
  gnupg \
  sudo \
  netcat \
  libyaml-cpp-dev \
  libgoogle-glog-dev \
  libprotoc-dev \
  libtins-dev \
  libmnl-dev \
  libsctp-dev \
  psmisc \
  openssl \
  net-tools \
  tshark \
  tzdata \
  iproute2 \
  iptables \
  && rm -rf /var/lib/apt/lists/*

RUN echo "deb https://facebookconnectivity.jfrog.io/artifactory/list/dev-focal/ focal main" > /etc/apt/sources.list.d/magma.list
RUN wget -qO - https://facebookconnectivity.jfrog.io/artifactory/api/security/keypair/default-gpg-key/public | apt-key add -

RUN apt-get update && apt-get install -y libopenvswitch openvswitch-datapath-dkms openvswitch-common openvswitch-switch

# Copy pre-built shared object files
COPY --from=builder /usr/lib/x86_64-linux-gnu/liblsan.so.0 /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libprotobuf.so /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libboost* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libevent-* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libsnappy.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libdouble-conversion.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libicui18n.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libicuuc.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libicudata.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libconfig.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libczmq.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libevent* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libzmq* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libsodium* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libpgm* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libnorm* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libgflags* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libgssapi_krb5* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libkrb5* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libk5crypto* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libkeyutils* /usr/lib/x86_64-linux-gnu/


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
COPY --from=builder $C_BUILD/session_manager/sessiond /usr/local/bin/sessiond
COPY --from=builder $C_BUILD/li_agent/src/liagentd /usr/local/bin/liagentd
COPY --from=builder $C_BUILD/sctpd/src/sctpd /usr/local/bin/sctpd
COPY --from=builder $C_BUILD/connection_tracker/src/connectiond /usr/local/bin/connectiond
COPY --from=builder $C_BUILD/core/oai/oai_mme/mme /usr/local/bin/oai_mme

RUN ldconfig 2> /dev/null

# Copy the configs.
COPY lte/gateway/configs /etc/magma
COPY orc8r/gateway/configs/templates /etc/magma/templates
COPY lte/gateway/deploy/roles/magma/files/magma-create-gtp-port.sh /usr/local/bin/

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

# debian stretch image
# ARG OS_DIST=debian
# ARG OS_RELEASE=stretch
# ARG EXTRA_REPO=https://facebookconnectivity.jfrog.io/artifactory/list/dev/
# ARG CLANG_VERSION=3.8

# ubuntu bionic image
ARG OS_DIST=ubuntu
ARG OS_RELEASE=bionic
ARG EXTRA_REPO=https://facebookconnectivity.jfrog.io/artifactory/cwf-prod
ARG CLANG_VERSION=3.9

# Stretch is required for c build
FROM $OS_DIST:$OS_RELEASE AS builder
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO
ARG CLANG_VERSION

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https gnupg
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub

RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb ${EXTRA_REPO} ${OS_RELEASE} main"


# Install dependencies required for building
RUN apt-get -y update && apt-get -y install \
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
  clang-${CLANG_VERSION} \
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
  nlohmann-json-dev \
  magma-cpp-redis \
  libgoogle-glog-dev \
  prometheus-cpp-dev \
  libfolly-dev \
  magma-libfluid \
  libdouble-conversion-dev \
  libboost-chrono-dev \
  libboost-context-dev \
  libboost-program-options-dev \
  libboost-filesystem-dev \
  libboost-regex-dev

ENV MAGMA_ROOT /magma
ENV C_BUILD /build/c
ENV OAI_BUILD $C_BUILD/oai

ENV CCACHE_DIR $MAGMA_ROOT/.cache/gateway/ccache
ENV MAGMA_DEV_MODE 1
ENV XDG_CACHE_HOME $MAGMA_ROOT/.cache

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
ARG BUILD_TYPE=RelWithDebInfo
ENV BUILD_TYPE=$BUILD_TYPE
RUN make -C $MAGMA_ROOT/lte/gateway/ build_session_manager BUILD_TYPE="${BUILD_TYPE}"

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM $OS_DIST:$OS_RELEASE AS gateway_c
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https gnupg
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub

RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb $EXTRA_REPO $OS_RELEASE main"


# Install runtime dependencies
RUN apt-get -y update && apt-get -y install \
  curl \
  sudo \
  # install prometheus
  prometheus-cpp-dev \
  # install openvswitch
  magma-libfluid \
  # install lxml
  python3-lxml \
  bridge-utils \
  # install yaml parser
  libyaml-cpp-dev \
  libgoogle-glog-dev \
  # folly deps
  libfolly-dev \
  libdouble-conversion-dev \
  libboost-chrono-dev \
  libboost-context-dev \
  libboost-program-options-dev \
  libboost-filesystem-dev \
  libboost-regex-dev \
  nlohmann-json-dev \
  redis-server \
  python-redis \
  magma-cpp-redis \
  grpc-dev \
  protobuf-compiler \
  libprotoc-dev \
  netcat

# Copy the build artifacts.
COPY --from=builder /build/c/session_manager/sessiond /usr/local/bin/sessiond

# Copy the configs.
COPY lte/gateway/configs /etc/magma
COPY orc8r/gateway/configs/templates /etc/magma/templates
RUN mkdir -p /var/opt/magma/configs

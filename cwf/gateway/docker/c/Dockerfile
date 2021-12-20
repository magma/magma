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

# ubuntu bionic image
ARG OS_DIST=ubuntu OS_RELEASE=bionic EXTRA_REPO=https://facebookconnectivity.jfrog.io/artifactory/cwf-prod

# Stretch is required for c build
FROM $OS_DIST:$OS_RELEASE AS builder
ARG OS_DIST OS_RELEASE EXTRA_REPO

RUN apt-get update && \
    # Setup necessary tools for adding the Magma repository
    apt-get install -y apt-utils software-properties-common apt-transport-https gnupg wget && \
    # Download Bazel
    wget -P /usr/sbin --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-amd64 && \
    chmod +x /usr/sbin/bazelisk-linux-amd64 && \
    ln -s /usr/sbin/bazelisk-linux-amd64 /usr/sbin/bazel

# Add the magma apt repo
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb ${EXTRA_REPO} ${OS_RELEASE} main"

# Install dependencies required for building
RUN apt-get -y update && apt-get -y install \
  sudo \
  curl \
  unzip \
  cmake \
  git \
  gcc \
  g++ \
  build-essential \
  # folly deps
  libfolly-dev \
  libgoogle-glog-dev \
  libgflags-dev \
  libevent-dev \
  libdouble-conversion-dev \
  libiberty-dev \
  libdouble-conversion-dev \
  libboost-chrono-dev \
  libboost-context-dev \
  libboost-program-options-dev \
  libboost-filesystem-dev \
  libboost-regex-dev

ENV MAGMA_ROOT /magma
WORKDIR /magma

COPY . ${MAGMA_ROOT}

# Build external dependencies first. This will help not rebuilt all dependencies triggered by Magma changes.
RUN bazel build \
  @com_github_grpc_grpc//:grpc++ \
  @com_github_google_glog//:glog \
  @protobuf//:protobuf \
  @prometheus_cpp//:prometheus-cpp \
  @yaml-cpp//:yaml-cpp \
  @github_nlohmann_json//:json

RUN bazel --bazelrc=${MAGMA_ROOT}/bazel/bazelrcs/cwag.bazelrc build //lte/gateway/c/session_manager:sessiond

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM $OS_DIST:$OS_RELEASE AS gateway_c

# Copy runtime dependencies
COPY --from=builder /usr/lib/x86_64-linux-gnu/libboost* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libevent-* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libgflags* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libglog* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libsnappy.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libdouble-conversion.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libicui18n.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libicuuc.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libicudata.* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libssl* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/openssl* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libcrypto* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/lib/x86_64-linux-gnu/libunwind* /usr/lib/x86_64-linux-gnu/
COPY --from=builder /usr/local/lib/libfolly.so /usr/local/lib/

RUN ldconfig 2> /dev/null

# Copy the build artifacts.
COPY --from=builder /magma/bazel-bin/lte/gateway/c/session_manager/sessiond /usr/local/bin/sessiond

# Copy the configs.
COPY lte/gateway/configs /etc/magma
COPY orc8r/gateway/configs/templates /etc/magma/templates
RUN mkdir -p /var/opt/magma/configs

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
ARG OS_DIST=ubuntu
ARG OS_RELEASE=bionic

# Stretch is required for c build
FROM $OS_DIST:$OS_RELEASE AS builder
ARG OS_DIST
ARG OS_RELEASE

ENV TZ=America/New_York \
    DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get upgrade -y 

RUN apt-get install -y --no-install-recommends \
        apt-transport-https \
        apt-utils \
        build-essential \
        ca-certificates \
        curl \
        gcc \
        git \
        gnupg2 \
        g++ \
        wget

# Install bazel
WORKDIR /usr/sbin
RUN wget --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-amd64 && \
    chmod +x bazelisk-linux-amd64 && \
    ln -s /usr/sbin/bazelisk-linux-amd64 /usr/sbin/bazel

# Dependencies for rules_folly @ https://github.com/storypku/rules_folly
RUN apt-get -y install --no-install-recommends \
    autoconf \
    automake \
    libtool \
    libssl-dev

# Update shared library configuration
RUN ldconfig -v

ENV MAGMA_ROOT /magma

# Copy proto files
COPY feg/protos ${MAGMA_ROOT}/feg/protos
COPY feg/gateway/services/aaa/protos ${MAGMA_ROOT}/feg/gateway/services/aaa/protos
COPY lte/protos ${MAGMA_ROOT}/lte/protos
COPY orc8r/protos ${MAGMA_ROOT}/orc8r/protos
COPY protos ${MAGMA_ROOT}/protos

# Copy bazel related files
COPY WORKSPACE.bazel ${MAGMA_ROOT}/WORKSPACE.bazel
COPY .bazelignore ${MAGMA_ROOT}/
COPY .bazelrc ${MAGMA_ROOT}/
COPY .bazelversion ${MAGMA_ROOT}/
COPY *.bzl ${MAGMA_ROOT}/
COPY *.bazel ${MAGMA_ROOT}/
COPY third_party ${MAGMA_ROOT}/third_party

# Mount bazel cache
COPY .bazel-cache ${MAGMA_ROOT}/
COPY .bazel-cache-repo ${MAGMA_ROOT}/

# Copy session_manager source code
COPY orc8r/gateway/c/common ${MAGMA_ROOT}/orc8r/gateway/c/common
COPY lte/gateway/c ${MAGMA_ROOT}/lte/gateway/c

# Build session_manager
WORKDIR ${MAGMA_ROOT}
# Sentry isn't supported on this base image (breakpad will not compile)
RUN bazel build //lte/gateway/c/session_manager:sessiond --define=disable_sentry_native=1

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM $OS_DIST:$OS_RELEASE AS gateway_c
ARG OS_DIST
ARG OS_RELEASE

# Copy the build artifacts.
COPY --from=builder /magma/bazel-bin/lte/gateway/c/session_manager/sessiond /usr/local/bin/sessiond

# Copy the configs.
COPY lte/gateway/configs /etc/magma
COPY orc8r/gateway/configs/templates /etc/magma/templates
RUN mkdir -p /var/opt/magma/configs

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

ARG OS_DIST=ubuntu OS_RELEASE=focal EXTRA_REPO=https://linuxfoundation.jfrog.io/artifactory/magma-packages-test

FROM ghcr.io/magma/magma/bazel-cache-plain:latest AS builder

ENV MAGMA_ROOT /magma
WORKDIR /magma

# Copy Bazel files at root and third_party
COPY WORKSPACE.bazel BUILD.bazel .bazelignore .bazelrc .bazelversion ${MAGMA_ROOT}/
COPY bazel/ ${MAGMA_ROOT}/bazel

# Copy proto files
COPY feg/protos ${MAGMA_ROOT}/feg/protos
COPY feg/gateway/services/aaa/protos ${MAGMA_ROOT}/feg/gateway/services/aaa/protos
COPY lte/protos ${MAGMA_ROOT}/lte/protos
COPY orc8r/protos ${MAGMA_ROOT}/orc8r/protos
COPY protos ${MAGMA_ROOT}/protos

# Build session_manager c code
COPY orc8r/gateway/c/common ${MAGMA_ROOT}/orc8r/gateway/c/common
COPY lte/gateway/c/session_manager ${MAGMA_ROOT}/lte/gateway/c/session_manager

RUN bazel build //lte/gateway/c/session_manager:sessiond --define=disable_sentry_native=1

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

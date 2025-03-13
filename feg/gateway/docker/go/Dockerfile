# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# -----------------------------------------------------------------------------
# Development image for test, precommit, etc.
# -----------------------------------------------------------------------------
ARG baseImage="ubuntu:focal"
FROM ${baseImage} AS base

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https
# Add the magma apt repo
COPY keys/linux_foundation_registry_key.asc /etc/apt/trusted.gpg.d/magma.asc

#[TODO_GPG_KEY_WORKAROUND]Temporary steps to be removed once GPG public key issue has been fixed.
#RUN add-apt-repository "deb https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main"
RUN echo "deb [trusted=yes] https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main" \
    > /etc/apt/sources.list.d/magma.list

#[TODO_GPG_KEY_WORKAROUND]Temporary steps to be removed once GPG public key issue has been fixed.
RUN echo "Acquire::AllowInsecureRepositories true;" > /etc/apt/apt.conf.d/99AllowInsecureRepositories && \
    echo "APT::Get::AllowUnauthenticated true;" >> /etc/apt/apt.conf.d/99AllowInsecureRepositories

# Install the runtime deps.
RUN apt-get update && apt-get install -y \
    bzr \
    curl \
    daemontools \
    gcc \
    git \
    libc-ares-dev \
    libev-dev \
    libevent-dev \
    libffi-dev \
    libjansson-dev \
    libjemalloc-dev \
    libssl-dev \
    libsystemd-dev \
    magma-nghttpx=1.31.1-1 \
    make \
    net-tools \
    pkg-config \
    python-cffi \
    python3-pip \
    redis-server \
    rsyslog \
    sudo \
    unzip \
    vim \
    virtualenv

# Golang
WORKDIR /usr/local
ARG GOLANG_VERSION="1.21.0"
RUN GO_TARBALL="go${GOLANG_VERSION}.linux-amd64.tar.gz" \
 && curl https://golang.org/dl/${GO_TARBALL} --remote-name --location \
 && tar -xzf ${GO_TARBALL} \
 && ln -s /usr/local/go/bin/go /usr/local/bin/go \
 && ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt \
 && rm ${GO_TARBALL}

# Install protobuf compiler.y
RUN curl -Lfs https://github.com/protocolbuffers/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip && \
    unzip protoc3.zip -d protoc3 && \
    mv protoc3/bin/protoc /bin/protoc && \
    chmod a+rx /bin/protoc && \
    mv protoc3/include/google /usr/include/ && \
    chmod -R a+Xr /usr/include/google && \
    rm -rf protoc3.zip protoc3

ENV GOBIN /var/opt/magma/bin
ENV MAGMA_ROOT /magma
ENV PIP_CACHE_HOME ~/.pipcache
ENV PYTHON_BUILD /build/python
ENV PATH ${PYTHON_BUILD}/bin:${PATH}:${GOBIN}
ENV GO111MODULE on
# Use public go modules proxy
ENV GOPROXY https://proxy.golang.org

RUN printenv > /etc/environment

COPY feg/radius/lib/go/ $MAGMA_ROOT/feg/radius/lib/go
COPY feg/radius/src/go.* $MAGMA_ROOT/feg/radius/src/
COPY orc8r/lib/go/go.* $MAGMA_ROOT/orc8r/lib/go/
WORKDIR $MAGMA_ROOT/feg/radius/src
RUN go mod download

# Copy just the go.mod and go.sum files to download the golang deps.
# This step allows us to cache the downloads, and prevents reaching out to
# the internet unless any of the go.mod or go.sum files are changed.
COPY lte/cloud/go/go.* $MAGMA_ROOT/lte/cloud/go/
COPY feg/cloud/go/go.* $MAGMA_ROOT/feg/cloud/go/
COPY feg/cloud/go/protos/go.* $MAGMA_ROOT/feg/cloud/go/protos/
COPY feg/gateway/go.* $MAGMA_ROOT/feg/gateway/
COPY orc8r/lib/go/protos/go.* $MAGMA_ROOT/orc8r/lib/go/protos/
COPY orc8r/cloud/go/go.* $MAGMA_ROOT/orc8r/cloud/go/
COPY orc8r/gateway/go/go.* $MAGMA_ROOT/orc8r/gateway/go/
WORKDIR $MAGMA_ROOT/feg/gateway
RUN go mod download

# Install protoc-gen-go
RUN go install github.com/golang/protobuf/protoc-gen-go \
 && go get -u github.com/vektra/mockery/v2/.../ \
 && go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen

# Symlink python scripts.
RUN ln -s /build/python/bin/generate_service_config.py /usr/local/bin/generate_service_config.py \
 && ln -s /build/python/bin/generate_nghttpx_config.py /usr/local/bin/generate_nghttpx_config.py

# Build the code.
COPY feg $MAGMA_ROOT/feg
COPY lte/cloud $MAGMA_ROOT/lte/cloud
COPY orc8r/lib/go $MAGMA_ROOT/orc8r/lib/go
COPY orc8r/cloud $MAGMA_ROOT/orc8r/cloud
COPY orc8r/gateway/go $MAGMA_ROOT/orc8r/gateway/go

# -----------------------------------------------------------------------------
# Builder image with binaries
# -----------------------------------------------------------------------------
FROM base AS builder
# Enable make gen if proto gen is required
# RUN make -C $MAGMA_ROOT/feg/gateway gen
RUN make -C $MAGMA_ROOT/feg/gateway build

WORKDIR $MAGMA_ROOT/feg/radius/src
RUN ./run.sh build

# -----------------------------------------------------------------------------
# Go-cache base image
# -----------------------------------------------------------------------------
FROM ${baseImage} AS gocache
COPY --from=builder /root/.cache /root/.cache

# -----------------------------------------------------------------------------
# Production image
# -----------------------------------------------------------------------------
FROM ${baseImage} AS gateway_go
ARG MAGMA_BUILD_BRANCH=unknown
ARG MAGMA_BUILD_TAG=unknown
ARG MAGMA_BUILD_COMMIT_HASH=unknonw
ARG MAGMA_BUILD_COMMIT_DATE=unknown

# Install envdir.
RUN apt-get -y update && apt-get -y install daemontools netcat gettext musl

ENV MAGMA_ROOT /magma
# Copy the build artifacts.
COPY --from=builder /var/opt/magma/bin /var/opt/magma/bin
COPY --from=builder $MAGMA_ROOT/feg/radius/src/radius /var/opt/magma/bin/radius
COPY --from=builder $MAGMA_ROOT/feg/radius/src/config/samples/radius.cwf.config.json.template /etc/magma/templates/radius.conf.template
COPY --from=builder $MAGMA_ROOT/feg/radius/src/config/samples/radius.cwf_analytics_disabled.config.json.template /etc/magma/templates/radius.analytics_disabled.conf.template

# Copy the configs.
COPY feg/gateway/configs /etc/magma

# Create empty envdir directory
RUN mkdir -p /var/opt/magma/envdir
RUN mkdir -p /var/opt/magma/configs
RUN mkdir -p /var/opt/magma/tmp

# Add commit information
ENV MAGMA_BUILD_BRANCH $MAGMA_BUILD_BRANCH
ENV MAGMA_BUILD_TAG $MAGMA_BUILD_TAG
ENV MAGMA_BUILD_COMMIT_HASH $MAGMA_BUILD_COMMIT_HASH
ENV MAGMA_BUILD_COMMIT_DATE $MAGMA_BUILD_COMMIT_DATE

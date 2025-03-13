# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# ------------------------------------------------------------------------------
# Base: for tests, precommit, codegen, etc.
# ------------------------------------------------------------------------------
ARG PLATFORM=linux/amd64

FROM --platform=$PLATFORM ubuntu:focal AS base

ENV GO111MODULE on
ENV GOPATH ${USER}/go
ENV GOBIN /build/bin
ENV PATH ${PATH}:${GOBIN}:${GOPATH}/bin

# Used in Makefiles
ENV MAGMA_ROOT /src/magma

# Apt runtime deps
ARG DEBIAN_FRONTEND=noninteractive
ARG TZ=UTC
RUN apt-get update && apt-get install -y \
  bzr \
  curl \
  gcc \
  git \
  make \
  openjdk-8-jre-headless \
  unzip \
  vim

# Golang
WORKDIR /usr/local
ARG GOLANG_VERSION="1.20.1"
RUN GO_TARBALL="go${GOLANG_VERSION}.linux-amd64.tar.gz" \
 && curl https://linuxfoundation.jfrog.io/artifactory/magma-blob/${GO_TARBALL} --remote-name --location \
 && tar -xzf ${GO_TARBALL} \
 && ln -s /usr/local/go/bin/go /usr/local/bin/go \
 && ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt \
 && rm ${GO_TARBALL}

# Install goimports
# RUN go get golang.org/x/tools/cmd/goimports

# Protobuf compiler
# Apt has 2.x but we need 3.x
# See: https://grpc.io/docs/protoc-installation/
RUN curl -Lfs https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip -o protoc3.zip && \
    unzip protoc3.zip -d protoc3 && \
    mv protoc3/bin/protoc /bin/protoc && \
    chmod a+rx /bin/protoc && \
    mv protoc3/include/google /usr/include/ && \
    chmod -R a+Xr /usr/include/google && \
    rm -rf protoc3.zip protoc3

# ------------------------------------------------------------------------------
# Gocache: cache Go modules
# ------------------------------------------------------------------------------
FROM base AS gocache

ARG MAGMA_MODULES="orc8r lte feg cwf dp"
RUN echo "export GOCACHE_MODULES=\"$(for m in $MAGMA_MODULES ; do echo -n /gomod/src/magma/$m ; echo -n ' ' ; done)\"" >> /etc/profile.d/env.sh

COPY gomod /gomod
RUN cat /etc/profile.d/env.sh
RUN . /etc/profile.d/env.sh && for m in $GOCACHE_MODULES ; do cd ${m}/cloud/go && echo ${m}/cloud/go && go mod download ; done

# ------------------------------------------------------------------------------
# Src: different src depending on MAGMA_MODULES
# ------------------------------------------------------------------------------
FROM gocache AS src

ARG MAGMA_MODULES="orc8r lte feg cwf dp"
RUN echo "export MAGMA_MODULES=\"$(for m in $MAGMA_MODULES ; do echo -n /src/magma/$m ; echo -n ' ' ; done)\"" >> /etc/profile.d/env.sh

# Source code
COPY src /src
WORKDIR /src/magma/orc8r/cloud
RUN . /etc/profile.d/env.sh && make tools

# Configs
COPY configs /etc/magma/configs

# ------------------------------------------------------------------------------
# Builder: compile src
# ------------------------------------------------------------------------------
FROM src AS builder

RUN . /etc/profile.d/env.sh && make build

# ------------------------------------------------------------------------------
# Production
# ------------------------------------------------------------------------------
FROM ubuntu:focal

# Apt runtime deps
RUN apt-get update && apt-get install -y \
  daemontools \
  netcat \
  openssl \
  supervisor \
  unzip \
  wget \
&& rm -rf /var/lib/apt/lists/*

# Swagger UI
# See: https://github.com/swagger-api/swagger-ui
ARG SWAGGER_UI_VERSION="3.52.2"
RUN cd /tmp && \
  wget "https://github.com/swagger-api/swagger-ui/archive/v$SWAGGER_UI_VERSION.zip" && \
  unzip "v$SWAGGER_UI_VERSION.zip" -d swagger-ui && \
  mkdir -p /var/opt/magma/static/swagger-ui && \
  cp -r "swagger-ui/swagger-ui-$SWAGGER_UI_VERSION/dist" /var/opt/magma/static/swagger-ui

# Script to wait for DB to be reachable
COPY src/magma/orc8r/cloud/docker/wait-for-it.sh /usr/local/bin

# Empty envdir for overriding in production
RUN mkdir -p /var/opt/magma/envdir

# Build artifacts
ARG SWAGGER_FILES=src/magma/orc8r/cloud/go/services/obsidian/swagger
COPY --from=builder /${SWAGGER_FILES}/v1/index.html /var/opt/magma/static/swagger/v1/ui/index.html
COPY --from=builder /${SWAGGER_FILES}/v1/css/sidebar.css /var/opt/magma/static/swagger/v1/static/sidebar.css
COPY --from=builder /${SWAGGER_FILES}/v1/swagger.yml /var/opt/magma/static/swagger/v1/spec/swagger.yml
COPY --from=builder src/magma/orc8r/cloud/swagger /etc/magma/swagger
COPY --from=builder /build/bin /var/opt/magma/bin

# Supervisor configs
ARG CNTLR_FILES=src/magma/orc8r/cloud/docker/controller
COPY ${CNTLR_FILES}/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ${CNTLR_FILES}/supervisor_logger.py /usr/local/lib/python3.8/dist-packages/supervisor_logger.py

# Scripts for dev mode
COPY ${CNTLR_FILES}/create_test_controller_certs /usr/local/bin/create_test_controller_certs

# Configs
COPY configs /etc/magma/configs

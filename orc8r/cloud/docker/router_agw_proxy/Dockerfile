# ------------------------------------------------------------------------------
# Base: for tests, precommit, codegen, etc.
# ------------------------------------------------------------------------------
ARG PLATFORM=linux/amd64

FROM --platform=$PLATFORM ubuntu:xenial as base

ENV GO111MODULE on
ENV GOPATH ${USER}/go
ENV GOBIN /build/bin
ENV PATH ${PATH}:${GOBIN}:${GOPATH}/bin

# Used in Makefiles
ENV MAGMA_ROOT /src/magma

# Apt runtime deps
RUN apt-get update && apt-get install -y \
  bzr \
  curl \
  gcc \
  git \
  make \
  openjdk-8-jre-headless \
  unzip \
  vim \
  ca-certificates && update-ca-certificates \
  && rm -rf /var/lib/apt/lists/*



# Golang 1.18
WORKDIR /usr/local
ARG GOLANG_VERSION="1.18.3"
RUN GO_TARBALL="go${GOLANG_VERSION}.linux-amd64.tar.gz" \
 && curl https://linuxfoundation.jfrog.io/artifactory/magma-blob/${GO_TARBALL} --remote-name --location \
 && tar -xzf ${GO_TARBALL} \
 && ln -s /usr/local/go/bin/go /usr/local/bin/go \
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
FROM base as gocache

ARG MAGMA_MODULES="orc8r lte feg cwf dp"
RUN echo "export GOCACHE_MODULES=\"$(for m in $MAGMA_MODULES ; do echo -n /gomod/src/magma/$m ; echo -n ' ' ; done)\"" >> /etc/profile.d/env.sh

COPY gomod /gomod
RUN cat /etc/profile.d/env.sh
RUN . /etc/profile.d/env.sh && for m in $GOCACHE_MODULES ; do cd ${m}/cloud/go && echo ${m}/cloud/go && go mod download ; done

# ------------------------------------------------------------------------------
# Src: different src depending on MAGMA_MODULES
# ------------------------------------------------------------------------------
FROM gocache as src

ARG MAGMA_MODULES="orc8r lte feg cwf dp"
RUN echo "export MAGMA_MODULES=\"$(for m in $MAGMA_MODULES ; do echo -n /src/magma/$m ; echo -n ' ' ; done)\"" >> /etc/profile.d/env.sh

# Source code
COPY src /src
WORKDIR /src/magma/orc8r/cloud
RUN . /etc/profile.d/env.sh && make tools

# Configs
COPY configs /etc/magma/configs
# ----------------------------------------------------------------------
# Builder: build only router_agw_proxy
# ----------------------------------------------------------------------
FROM src as builder

RUN . /etc/profile.d/env.sh && make build
WORKDIR /src/magma/orc8r/cloud/go/services/router_agw_proxy
RUN go build -o /router_agw_proxy main.go

# ----------------------------------------------------------------------
# Production: minimal runtime image
# ----------------------------------------------------------------------
FROM debian:bullseye-slim

# Install CA certificates for HTTPS/TLS trust
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /router_agw_proxy .

# Scripts for dev mode
ARG CNTLR_FILES=src/magma/orc8r/cloud/docker/controller
COPY ${CNTLR_FILES}/create_test_controller_certs /usr/local/bin/create_test_controller_certs

# Configs
COPY configs /etc/magma/configs

EXPOSE 9097
CMD ["./router_agw_proxy"]

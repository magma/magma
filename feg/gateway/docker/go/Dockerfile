# -----------------------------------------------------------------------------
# Development image for test, precommit, etc.
# -----------------------------------------------------------------------------
FROM ubuntu:xenial as base

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb https://magma.jfrog.io/magma/list/dev/ xenial main"

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

# Install Golang 1.11.
WORKDIR /usr/local
RUN curl https://dl.google.com/go/go1.11.linux-amd64.tar.gz -O && \
    tar xvf go1.11.linux-amd64.tar.gz && \
    cp -r go/bin/* /usr/local/bin/

# Install protobuf compiler.
RUN curl -Lfs https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip && \
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
RUN printenv > /etc/environment

# Copy just the go.mod and go.sum files to download the golang deps.
# This step allows us to cache the downloads, and prevents reaching out to
# the internet unless any of the go.mod or go.sum files are changed.
COPY lte/cloud/go/go.* $MAGMA_ROOT/lte/cloud/go/
COPY feg/cloud/go/go.* $MAGMA_ROOT/feg/cloud/go/
COPY feg/cloud/go/protos/go.* $MAGMA_ROOT/feg/cloud/go/protos/
COPY feg/gateway/go.* $MAGMA_ROOT/feg/gateway/
COPY feg/gateway/third-party/go/src/github.com/fiorix/go-diameter/go.* $MAGMA_ROOT/feg/gateway/third-party/go/src/github.com/fiorix/go-diameter/
COPY feg/gateway/services/eap/client/go.* $MAGMA_ROOT/feg/gateway/services/eap/client/
COPY orc8r/cloud/go/go.* $MAGMA_ROOT/orc8r/cloud/go/
WORKDIR $MAGMA_ROOT/feg/gateway
RUN go mod download

# Copy the configs.
COPY feg/gateway/deploy/roles/feg_dev/files/magma_*.service  /etc/systemd/system/
COPY feg/gateway/deploy/roles/feg_services/files/magma_*.service /etc/systemd/system/
COPY lte/gateway/deploy/roles/magma/files/systemd/magma_magmad.service /etc/systemd/system/
COPY orc8r/tools/ansible/roles/gateway_services/files/magma.service /etc/systemd/system/magma@.service
COPY orc8r/tools/ansible/roles/gateway_services/files/magma_control_proxy.service /etc/systemd/system/
COPY orc8r/cloud/deploy/roles/proxy/files/magma_headers.rb /etc/nghttpx/magma_headers.rb

# Rename configs to start with 'magma@' instead of 'magma_'.
RUN for f in /etc/systemd/system/magma_*.service; do mv $f $(echo $f | sed s/magma_/magma@/); done

# Symlink python scripts.
RUN ln -s /build/python/bin/generate_service_config.py /usr/local/bin/generate_service_config.py
RUN ln -s /build/python/bin/generate_nghttpx_config.py /usr/local/bin/generate_nghttpx_config.py

# -----------------------------------------------------------------------------
# Builder image with binaries and plugins
# -----------------------------------------------------------------------------
FROM base as builder

# Build the code.
COPY feg $MAGMA_ROOT/feg
COPY lte/cloud $MAGMA_ROOT/lte/cloud
COPY orc8r/cloud $MAGMA_ROOT/orc8r/cloud
RUN make -C $MAGMA_ROOT/feg/gateway build

# -----------------------------------------------------------------------------
# Production image
# -----------------------------------------------------------------------------
FROM ubuntu:xenial AS gateway_go

# Install envdir.
RUN apt-get -y update && apt-get -y install daemontools

# Copy the build artifacts.
COPY --from=builder /var/opt/magma/bin /var/opt/magma/bin

# Copy the configs.
COPY feg/gateway/configs /etc/magma
# Add docker config overrides
COPY feg/gateway/docker/configs /etc/magma

# Create empty envdir directory
RUN mkdir -p /var/opt/magma/envdir

RUN mkdir -p /var/opt/magma/configs

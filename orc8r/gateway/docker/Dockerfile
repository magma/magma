# Builder image to generate proto files
# -----------------------------------------------------------------------------
FROM ubuntu:focal AS builder

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install curl make virtualenv zip \
  apt-utils software-properties-common apt-transport-https

# Install protobuf compiler.
RUN curl -Lfs https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip && \
  unzip protoc3.zip -d protoc3 && \
  mv protoc3/bin/protoc /usr/bin/protoc && \
  chmod a+rx /usr/bin/protoc && \
  cp -r protoc3/include/google /usr/include/ && \
  chmod -R a+Xr /usr/include/google && \
  rm -rf protoc3.zip protoc3

# Install python3.8 since it's not native on focal
RUN apt-add-repository "ppa:deadsnakes/ppa"
RUN apt-get -y update && apt-get -y install python3.8

ENV MAGMA_ROOT /magma
ENV PYTHON_BUILD /build/python
ENV PIP_CACHE_HOME ~/.pipcache

# Generate python proto bindings.
COPY lte/gateway/python/defs.mk $MAGMA_ROOT/lte/gateway/python/defs.mk
COPY lte/gateway/python/Makefile $MAGMA_ROOT/lte/gateway/python/Makefile
COPY lte/protos $MAGMA_ROOT/lte/protos
COPY orc8r/gateway/python/python.mk $MAGMA_ROOT/orc8r/gateway/python/python.mk
COPY orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY protos $MAGMA_ROOT/protos
RUN make -C $MAGMA_ROOT/lte/gateway/python protos

# -----------------------------------------------------------------------------
# Production image
# -----------------------------------------------------------------------------
FROM ubuntu:focal AS gateway_python

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
COPY cwf/gateway/deploy/roles/ovs/files/magma-preferences /etc/apt/preferences.d/
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb https://facebookconnectivity.jfrog.io/artifactory/list/dev-focal/ focal main"

# Install python3.8 since it's not native on focal
RUN apt-add-repository "ppa:deadsnakes/ppa"
RUN apt-get -y update && apt-get -y install python3.8

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install \
  curl \
  fabric \
  libc-ares2 \
  libev4 \
  libffi-dev \
  libjansson4 \
  libjemalloc2 \
  libssl-dev \
  libsystemd-dev \
  magma-nghttpx=1.31.1-1 \
  net-tools \
  openssh-server \
  openssl \
  pkg-config \
  python-cffi \
  python3-pip \
  python3.8-dev \
  redis-server \
  network-manager

RUN curl -sSL https://get.docker.com/ > /tmp/get_docker.sh && \
    sh /tmp/get_docker.sh && \
    rm /tmp/get_docker.sh

# Install python code.
COPY orc8r/gateway/python /tmp/orc8r
RUN python3.8 -m pip install --no-cache-dir /tmp/orc8r

# Copy the build artifacts.
COPY --from=builder /build/python/gen /usr/local/lib/python3.8/dist-packages/

# Copy the configs.
COPY orc8r/gateway/configs /etc/magma

COPY orc8r/gateway/configs/templates /etc/magma/templates

RUN mkdir -p /var/opt/magma/configs

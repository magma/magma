# -----------------------------------------------------------------------------
# Builder image to generate OVS debian packages and Magma proto files
# -----------------------------------------------------------------------------
FROM ubuntu:xenial AS builder

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb https://magma.jfrog.io/magma/list/dev/ xenial main"

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install curl make virtualenv zip tar

# TODO: Fetch necessary OVS debian packages from pkgrepo
# Fetch OVS v2.8.1 tarball
RUN curl -Lfs https://www.openvswitch.org/releases/openvswitch-2.8.1.tar.gz | tar xvz

# Install OVS dependencies for building
RUN apt-get -y update && apt-get -y install \
  graphviz \
  autoconf \
  automake \
  bzip2 \
  debhelper \
  dh-autoreconf \
  libssl-dev \
  libtool \
  openssl \
  procps \
  python-all \
  python-twisted-conch \
  python-zopeinterface \
  python-six \
  build-essential \
  fakeroot

# Build OVS v2.8.1 debian packages
RUN cd openvswitch-2.8.1 && DEB_BUILD_OPTIONS='parallel=8 nocheck' fakeroot debian/rules binary

# Install protobuf compiler.
RUN curl -Lfs https://github.com/google/protobuf/releases/download/v3.1.0/protoc-3.1.0-linux-x86_64.zip -o protoc3.zip && \
  unzip protoc3.zip -d protoc3 && \
  mv protoc3/bin/protoc /usr/bin/protoc && \
  chmod a+rx /usr/bin/protoc && \
  cp -r protoc3/include/google /usr/include/ && \
  chmod -R a+Xr /usr/include/google && \
  rm -rf protoc3.zip protoc3

ENV MAGMA_ROOT /magma
ENV PYTHON_BUILD /build/python
ENV PIP_CACHE_HOME ~/.pipcache

# Generate python proto bindings.
COPY feg/protos $MAGMA_ROOT/feg/protos
COPY lte/gateway/python/defs.mk $MAGMA_ROOT/lte/gateway/python/defs.mk
COPY lte/gateway/python/Makefile $MAGMA_ROOT/lte/gateway/python/Makefile
COPY lte/protos $MAGMA_ROOT/lte/protos
COPY orc8r/gateway/python/python.mk $MAGMA_ROOT/orc8r/gateway/python/python.mk
COPY orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY protos $MAGMA_ROOT/protos
RUN make -C $MAGMA_ROOT/lte/gateway/python protos

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM ubuntu:xenial AS lte_gateway_python

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb https://magma.jfrog.io/magma/list/dev/ xenial main"

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install \
  curl \
  libc-ares2 \
  libev4 \
  libevent-openssl-2.0-5 \
  libffi-dev \
  libjansson4 \
  libjemalloc1 \
  libssl-dev \
  libsystemd-dev \
  magma-nghttpx=1.31.1-1 \
  net-tools \
  openssl \
  pkg-config \
  python-cffi \
  python3-aioeventlet \
  python3-pip \
  redis-server

# Copy OVS debian packages from builder
COPY --from=builder libopenvswitch_2.8.1-1_amd64.deb /tmp/ovs/
COPY --from=builder openvswitch-common_2.8.1-1_amd64.deb /tmp/ovs/
COPY --from=builder openvswitch-switch_2.8.1-1_amd64.deb /tmp/ovs/
COPY --from=builder python-openvswitch_2.8.1-1_all.deb /tmp/ovs/

# Install OVS debian packages
RUN apt-get -y install \
  /tmp/ovs/libopenvswitch_2.8.1-1_amd64.deb \
  /tmp/ovs/openvswitch-common_2.8.1-1_amd64.deb \
  /tmp/ovs/openvswitch-switch_2.8.1-1_amd64.deb \
  /tmp/ovs/python-openvswitch_2.8.1-1_all.deb

# Install orc8r python (magma.common required for lte python)
COPY orc8r/gateway/python /tmp/orc8r
RUN pip3 install /tmp/orc8r

# Install lte python
COPY lte/gateway/python /tmp/lte
RUN pip3 install /tmp/lte

# Copy the build artifacts.
COPY --from=builder /build/python/gen /usr/local/lib/python3.5/dist-packages/

# Copy the configs.
COPY lte/gateway/configs /etc/magma
COPY orc8r/gateway/configs/templates /etc/magma/templates
COPY orc8r/cloud/docker/proxy/magma_headers.rb /etc/nghttpx/magma_headers.rb
RUN mkdir -p /var/opt/magma/configs

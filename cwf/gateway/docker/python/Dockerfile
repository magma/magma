# -----------------------------------------------------------------------------
# Builder image for Magma proto files
# -----------------------------------------------------------------------------
FROM ubuntu:xenial AS builder

# Install the runtime deps from apt.
RUN apt-get -y update && apt-get -y install curl make virtualenv zip

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
COPY cwf/gateway/deploy/roles/ovs/files/magma-preferences /etc/apt/preferences.d/
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
    redis-server \
    iptables \
    git \
    automake \
    gcc \
    libtool \
    libcap-ng-dev \
    linux-headers-generic \
    netcat
    
RUN pip3 install \
    Cython \
    fire \
    envoy \
    glob2 \
    ryu \
    flask \
    aiodns \
    pymemoize \
    wsgiserver \
    pycrypto \
    six \
    eventlet \
    h2 \
    hpack \
    docker \
    redis \
    redis-collections \
    aiohttp \
    Jinja2 \
    netifaces \
    pylint \
    PyYAML \
    pytz \
    snowflake \
    systemd-python \
    itsdangerous \
    click \
    pycares \
    python-dateutil \
    jsonpickle 

# Install OVS via Magma bionic pkg repo
RUN git clone --single-branch --branch v2.12.0 https://github.com/openvswitch/ovs.git

COPY cwf/gateway/deploy/roles/ovs/files/0001-Add-custom-IPDR-fields-for-IPFIX-export.patch /tmp
COPY cwf/gateway/deploy/roles/ovs/files/0002-ovs-Handle-spaces-in-ovs-arguments.patch /tmp
COPY cwf/gateway/deploy/roles/ovs/files/0003-Add-pdp_start_epoch-custom-field-to-IPFIX-export.patch /tmp

WORKDIR ovs
RUN git apply /tmp/0001-Add-custom-IPDR-fields-for-IPFIX-export.patch
RUN git apply /tmp/0002-ovs-Handle-spaces-in-ovs-arguments.patch
RUN git apply /tmp/0003-Add-pdp_start_epoch-custom-field-to-IPFIX-export.patch
RUN ./boot.sh
RUN ./configure --prefix=/usr --localstatedir=/var --sysconfdir=/etc
RUN make
RUN make install

# Install orc8r python (magma.common required for lte python)
COPY orc8r/gateway/python /tmp/orc8r
RUN pip3 install /tmp/orc8r

# Install lte python
COPY lte/gateway/python /tmp/lte
RUN pip3 install /tmp/lte

# Copy the configs.
COPY lte/gateway/configs /etc/magma
COPY orc8r/gateway/configs/templates /etc/magma/templates
COPY orc8r/cloud/docker/proxy/magma_headers.rb /etc/nghttpx/magma_headers.rb
RUN mkdir -p /var/opt/magma/configs

WORKDIR /

# Copy the build artifacts.
COPY --from=builder /build/python/gen /usr/local/lib/python3.5/dist-packages/

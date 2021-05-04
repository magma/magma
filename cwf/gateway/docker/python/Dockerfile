# -----------------------------------------------------------------------------
# Builder image for Magma proto files
# -----------------------------------------------------------------------------
FROM ubuntu:focal AS builder

# workaround to avoid interactive tzdata configurtaion
ARG DEBIAN_FRONTEND=noninteractive

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
COPY cwf/protos $MAGMA_ROOT/cwf/protos
COPY feg/protos $MAGMA_ROOT/feg/protos
COPY lte/gateway/python/defs.mk $MAGMA_ROOT/lte/gateway/python/defs.mk
COPY lte/gateway/python/Makefile $MAGMA_ROOT/lte/gateway/python/Makefile
COPY lte/protos $MAGMA_ROOT/lte/protos
COPY orc8r/gateway/python $MAGMA_ROOT/orc8r/gateway/python
COPY orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY protos $MAGMA_ROOT/protos
ENV PROTO_LIST orc8r_protos lte_protos feg_protos cwf_protos
RUN make -C $MAGMA_ROOT/orc8r/gateway/python protos

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM ubuntu:focal AS lte_gateway_python

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
COPY cwf/gateway/deploy/roles/ovs/files/magma-preferences /etc/apt/preferences.d/
# TODO: change repo to focal once there will be such
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb https://facebookconnectivity.jfrog.io/artifactory/list/dev-focal/ focal main" && \
    apt-add-repository "deb http://archive.ubuntu.com/ubuntu/ focal-proposed restricted main multiverse universe"

 # Install python3.8 since it's not native on focal
RUN apt-add-repository "ppa:deadsnakes/ppa"
RUN apt-get -y update && apt-get -y install python3.8
RUN apt-get -y update && apt-get -y install \
    curl \
    libc-ares2 \
    libev4 \
    libffi-dev \
    libjansson4 \
    libjemalloc2 \
    libssl-dev \
    libsystemd-dev \
    magma-nghttpx=1.31.1-1 \
    net-tools \
    openssl \
    pkg-config \
    python-cffi \
    python3-pip \
    python3.8-dev \
    redis-server \
    iptables \
    git \
    automake \
    gcc \
    libtool \
    libcap-ng-dev \
    linux-headers-generic \
    netcat \
    iputils-ping

RUN python3.8 -m pip install --no-cache-dir \
    Cython \
    fire \
    envoy \
    glob2 \
    ryu \
    flask \
    aiodns \
    pymemoize \
    wsgiserver \
    pycryptodome \
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
    aioeventlet>=0.4 \ 
    jsonpickle

# TODO: aioeventlet>=0.4 was manually added to fix regression tracked with GH issue 6152
# Fix to pull the patched aioeventlet build deployed by the aioeventlet_build.sh script

COPY cwf/gateway/deploy/roles/ovs/files/nx_actions.py /usr/local/lib/python3.8/dist-packages/ryu/ofproto/

# Temporary workaround to restore uplink bridge flows
RUN mkdir -p /var/opt/magma/scripts
COPY cwf/gateway/deploy/roles/cwag/files/add_uplink_bridge_flows.sh /var/opt/magma/scripts

# Install OVS via Magma bionic pkg repo
RUN git clone --depth 1 --single-branch --branch v2.12.0 https://github.com/openvswitch/ovs.git

COPY cwf/gateway/deploy/roles/ovs/files/0001-Add-custom-IPDR-fields-for-IPFIX-export.patch /tmp
COPY cwf/gateway/deploy/roles/ovs/files/0002-ovs-Handle-spaces-in-ovs-arguments.patch /tmp
COPY cwf/gateway/deploy/roles/ovs/files/0003-Add-pdp_start_epoch-custom-field-to-IPFIX-export.patch /tmp
COPY cwf/gateway/deploy/roles/ovs/files//0004-ovsdb-idlc.in-dict-changes.patch /tmp
WORKDIR ovs
RUN git apply /tmp/0001-Add-custom-IPDR-fields-for-IPFIX-export.patch
RUN git apply /tmp/0002-ovs-Handle-spaces-in-ovs-arguments.patch
RUN git apply /tmp/0003-Add-pdp_start_epoch-custom-field-to-IPFIX-export.patch
RUN git apply /tmp/0004-ovsdb-idlc.in-dict-changes.patch
RUN ./boot.sh
RUN ./configure --prefix=/usr --localstatedir=/var --sysconfdir=/etc
RUN make
RUN make install

# Install orc8r python (magma.common required for lte python)
COPY orc8r/gateway/python /tmp/orc8r
RUN python3.8 -m pip install --no-cache-dir /tmp/orc8r

# Install lte python
COPY lte/gateway/python /tmp/lte
RUN python3.8 -m pip install --no-cache-dir /tmp/lte

# Copy the configs.
COPY lte/gateway/configs /etc/magma
COPY orc8r/gateway/configs/templates /etc/magma/templates
RUN mkdir -p /var/opt/magma/configs

WORKDIR /

# Copy the build artifacts.
COPY --from=builder /build/python/gen /usr/local/lib/python3.8/dist-packages/

# update aioh2 since there is no pushed package, but master is fixed
RUN python3.8 -m pip install --force-reinstall git+https://github.com/URenko/aioh2.git

# patching aioeventlet
RUN mkdir /patches
COPY lte/gateway/deploy/roles/magma/files/patches/aioeventlet.py38.patch /patches
RUN patch -N -s -f /usr/local/lib/python3.8/dist-packages/aioeventlet.py < patches/aioeventlet.py38.patch

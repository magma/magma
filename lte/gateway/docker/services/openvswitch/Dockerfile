ARG CPU_ARCH=x86_64
ARG DEB_PORT=amd64
ARG OS_DIST=ubuntu
ARG OS_RELEASE=focal
ARG EXTRA_REPO=https://linuxfoundation.jfrog.io/artifactory/magma-packages-test

FROM $OS_DIST:$OS_RELEASE AS gateway_ovs
ARG CPU_ARCH
ARG DEB_PORT
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO

ENV LINUX_HEADERS_VER 5.4.0-186-generic

## Install packages
RUN apt-get -q update && \
    apt-get -y -q --no-install-recommends install \
    # Setup necessary tools for adding the Magma repository
    apt-utils \
    ca-certificates \
    apt-transport-https \
    # networking utils
    iptables \
    iproute2 \
    iputils-arping \
    iputils-clockdiff \
    iputils-ping \
    iputils-tracepath \
    bridge-utils \
    ifupdown \
    vim \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

## Add the magma apt repo and install the custom openvswitch modules
COPY keys/linux_foundation_registry_key.asc /etc/apt/trusted.gpg.d/magma.asc
RUN echo "deb https://linuxfoundation.jfrog.io/artifactory/magma-packages focal-1.8.0 main" > /etc/apt/sources.list.d/magma.list
RUN apt-get update && apt-get install -y --no-install-recommends \
    libopenvswitch \
    openvswitch-common \
    openvswitch-switch \
    linux-headers-${LINUX_HEADERS_VER} \
    openvswitch-datapath-dkms \
    && rm -rf /var/lib/apt/lists/*

# Copy necessary files for ovs
COPY --chmod=755 lte/gateway/docker/services/openvswitch/healthcheck.sh /usr/local/bin/
COPY --chmod=755 lte/gateway/docker/services/openvswitch/entrypoint.sh /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ] 

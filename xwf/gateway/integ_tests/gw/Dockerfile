FROM ubuntu:bionic

# need to run update twice, this is not bug
RUN apt-get update
RUN apt-get update && apt-get install -y \
    git \
    sudo \
    iptables \
    ansible \
    dnsmasq \
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
    python-zope.interface \
    python-six \
    build-essential \
    fakeroot \
    jq \
    isc-dhcp-server \
    net-tools \
    iproute2 \
    wget \
    curl

WORKDIR /code
COPY xwf/gateway/deploy ./xwf/gateway/deploy

RUN ANSANSIBLE_CONFIG=xwf/gateway/ansible.cfg ansible-playbook xwf/gateway/deploy/xwf.yml -i "localhost," -c local \
  -e repo=cwf-prod -e distribution=bionic --tags "install" -v

COPY xwf ./xwf
COPY orc8r ./orc8r
COPY cwf ./cwf

# Create snowflake to be mounted into containers
RUN touch /etc/snowflake

# Placing configs in the appropriate place...
RUN mkdir -p /var/opt/magma
RUN mkdir -p /var/opt/magma/configs
RUN mkdir -p /var/opt/magma/certs
RUN mkdir -p /etc/magma
RUN mkdir -p /var/opt/magma/docker

COPY xwf/gateway/integ_tests/gw/entrypoint.sh .
CMD ./entrypoint.sh

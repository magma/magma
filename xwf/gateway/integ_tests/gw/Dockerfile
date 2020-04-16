FROM ubuntu:bionic

RUN apt-get update
RUN apt-get update \
  &&apt-get install -y \
  git \
  software-properties-common \
  sudo netplan \
  apt-utils \
  iptables \
  apt-transport-https \
  ca-certificates \
  curl \
  gnupg-agent \
  ansible \
  net-tools \
  docker.io \
  docker-compose \
  isc-dhcp-server \
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
  python-zope.interface \
  python-six \
  build-essential \
  fakeroot \
  dnsmasq \
  dialog

WORKDIR /code
COPY cwf ./cwf
COPY lte ./lte
COPY orc8r ./orc8r
COPY xwf ./xwf
RUN ANSIBLE_CONFIG=cwf/gateway/ansible.cfg  ansible-playbook cwf/gateway/deploy/cwag.yml -i "localhost," -c local --tags "install" -v

# Create snowflake to be mounted into containers
RUN touch /etc/snowflake

# Placing configs in the appropriate place...
RUN mkdir -p /var/opt/magma
RUN mkdir -p /var/opt/magma/configs
RUN mkdir -p /var/opt/magma/certs
RUN mkdir -p /etc/magma
RUN mkdir -p /var/opt/magma/docker

COPY orc8r/tools/docker/recreate_services.sh /var/opt/magma/docker/
COPY orc8r/tools/docker/recreate_services_cron /etc/cron.d/

COPY xwf/gateway/integ_tests/gw/entrypoint.sh .
CMD ./entrypoint.sh

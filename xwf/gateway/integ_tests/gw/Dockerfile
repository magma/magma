FROM centos:7

RUN yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
RUN yum update -y
RUN yum install -y \
  git \
  sudo \
  iptables \
  ansible \
  yum-plugin-versionlock \
  net-tools \
  initscripts \
  graphviz \
  bzip2 \
  openssl \
  procps \
  python-six \
  dnsmasq \
  dialog \
  wget \
  jq \
  dhcp

WORKDIR /code
COPY xwf/gateway/deploy ./xwf/gateway/deploy
RUN ANSANSIBLE_CONFIG=xwf/gateway/ansible.cfg ansible-playbook xwf/gateway/deploy/xwf.yml -i "localhost," -c local --tags "install" -v

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

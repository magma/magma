FROM ubuntu:xenial
ENV MAGMA_ROOT /root
ENV DOCKER_VERSION "17.03.0-ce"
RUN     apt-get update \
             && apt-get install -y \
             curl \
             docker \
             git \
             net-tools \
             python3-pip \
             python3-apt \
             unzip \
             vim \
             wget \
             ruby-dev &&\
             curl -L -o /tmp/docker-$DOCKER_VERSION.tgz https://get.docker.com/builds/Linux/x86_64/docker-$DOCKER_VERSION.tgz &&\
             tar -xz -C /tmp -f /tmp/docker-$DOCKER_VERSION.tgz &&\
             mv /tmp/docker/* /usr/bin &&\
             curl -L https://github.com/docker/compose/releases/download/1.25.4/docker-compose-`uname -s`-`uname -m` > ~/docker-compose &&\
             chmod +x ~/docker-compose &&\
             mv ~/docker-compose /usr/local/bin/docker-compose

VOLUME /var/run/docker.sock /var/run/docker.sock
WORKDIR /root

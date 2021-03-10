FROM ubuntu:latest
ENV MAGMA_ROOT /root/magma
RUN apt-get update && \
    apt-get -y install curl make virtualenv zip rsync git software-properties-common python3-pip python-dev && \
    apt-get -y upgrade openssh-server && \
    alias python=python3 && \
    pip3 install --no-cache-dir ansible

RUN mkdir -p /etc/ansible && \
    mkdir -p $MAGMA_ROOT

# Generating an ssh key
RUN mkdir -p ~/.ssh &&  \
    cd ~/.ssh && \
    ssh-keygen -N "" -f ~/.ssh/id_rsa && \
    chmod 600 id_rsa

ARG CACHE_DATE=2019-09-05
RUN git clone --depth 1 https://github.com/facebookincubator/magma.git $MAGMA_ROOT

RUN cp $MAGMA_ROOT/lte/gateway/docker/deploy/ssh_config /etc/ssh/ssh_config

WORKDIR $MAGMA_ROOT

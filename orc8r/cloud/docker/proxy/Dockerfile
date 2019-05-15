FROM ubuntu:xenial

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https
COPY src/magma/orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb https://magma.jfrog.io/magma/list/dev/ xenial main"

# Install the deps from apt
RUN apt-get update && \
    apt-get install -y \
        libssl-dev libev-dev libevent-dev libjansson-dev libjemalloc-dev libc-ares-dev magma-nghttpx=1.31.1-1 \
        daemontools \
        supervisor \
        python3-pip

# Install python3 deps from pip
RUN pip3 install PyYAML jinja2

# Create an empty envdir for overriding in production
RUN mkdir -p /var/opt/magma/envdir

ARG PROXY_FILES=src/magma/orc8r/cloud/docker/proxy

# Copy the scripts and configs from the context
COPY configs /etc/magma/configs
COPY ${PROXY_FILES}/templates /etc/magma/templates
COPY ${PROXY_FILES}/magma_headers.rb /etc/nghttpx/magma_headers.rb
COPY ${PROXY_FILES}/run_nghttpx.py /usr/local/bin/run_nghttpx.py
COPY ${PROXY_FILES}/create_test_proxy_certs /usr/local/bin/create_test_proxy_certs

# Copy the supervisor configs
COPY ${PROXY_FILES}/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
CMD ["/usr/bin/supervisord"]
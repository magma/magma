FROM ubuntu:xenial

# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y apt-utils software-properties-common apt-transport-https
COPY magma/orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
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

# Copy the scripts and configs from the context
COPY configs /etc/magma/configs
COPY magma/orc8r/cloud/deploy/roles/proxy/templates /etc/magma/templates
COPY magma/orc8r/cloud/deploy/files/envdir /var/opt/magma/envdir
COPY magma/orc8r/cloud/deploy/roles/proxy/files/magma_headers.rb /etc/nghttpx/magma_headers.rb
COPY magma/orc8r/cloud/deploy/roles/proxy/files/run_nghttpx.py /usr/local/bin/run_nghttpx.py
COPY magma/orc8r/cloud/scripts/create_test_proxy_certs /usr/local/bin/create_test_proxy_certs

# Override the dev environments to suit Docker
# TODO: Remove this step after removing vagrant
RUN echo "controller" > /var/opt/magma/envdir/PROXY_BACKENDS

# Copy the supervisor configs
COPY magma/orc8r/cloud/deploy/roles/proxy/files/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
CMD ["/usr/bin/supervisord"]
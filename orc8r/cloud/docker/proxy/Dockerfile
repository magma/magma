FROM ubuntu:xenial

ARG CNTLR_FILES=src/magma/orc8r/cloud/docker/controller
ARG PROXY_FILES=src/magma/orc8r/cloud/docker/proxy

# Install the runtime deps from apt
# Add the magma apt repo
RUN apt-get update && \
    apt-get install -y \
        daemontools netcat openssl supervisor wget unzip \
        apt-utils software-properties-common apt-transport-https
COPY src/magma/orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/jfrog.pub
RUN apt-key add /tmp/jfrog.pub && \
    apt-add-repository "deb https://magma.jfrog.io/magma/list/dev/ xenial main"

# Install the proxy deps from apt
RUN apt-get update && \
    apt-get install -y \
        libssl-dev libev-dev libevent-dev libjansson-dev libjemalloc-dev libc-ares-dev magma-nghttpx=1.31.1-1 \
        daemontools \
        supervisor \
        python3-pip \
        squid

# Install python3 deps from pip
RUN pip3 install PyYAML jinja2

# Create an empty envdir for overriding in production
RUN mkdir -p /var/opt/magma/envdir
# We only care to collect and report metrics on Proxy, and not much else
ENV CONTROLLER_SERVICES METRICSD
# TODO: fix dispatcher to work on containers
ENV SERVICE_HOST_NAME localhost
# TODO: fix metricsd for fetching host level metrics
ENV HOST_NAME dev

# Copy the configs
COPY configs /etc/magma/configs

# Copy proxy scripts and configs from the context
COPY ${PROXY_FILES}/templates /etc/magma/templates
COPY ${PROXY_FILES}/magma_headers.rb /etc/nghttpx/magma_headers.rb
COPY ${PROXY_FILES}/run_nghttpx.py /usr/local/bin/run_nghttpx.py
COPY ${PROXY_FILES}/squid.conf /etc/squid/squid.conf
COPY ${PROXY_FILES}/create_test_proxy_certs /usr/local/bin/create_test_proxy_certs

# Copy the build artifacts
COPY --from=orc8r_controller /var/opt/magma/bin /var/opt/magma/bin
COPY --from=orc8r_controller /var/opt/magma/plugins /var/opt/magma/plugins

# Copy the supervisor configs
COPY ${PROXY_FILES}/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ${CNTLR_FILES}/supervisor_logger.py /usr/local/lib/python2.7/dist-packages/supervisor_logger.py
CMD ["/usr/bin/supervisord"]

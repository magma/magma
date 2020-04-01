ARG NGHTTP_VERSION=1.40.0

FROM ubuntu:bionic as builder

# From nghttp2 readme, minus libsystemd-dev
RUN apt-get update && \
    apt-get install -y git curl g++ make binutils autoconf automake autotools-dev libtool pkg-config \
      zlib1g-dev libcunit1-dev openssl libssl-dev libxml2-dev libev-dev libevent-dev libjansson-dev \
      libc-ares-dev libjemalloc-dev ruby bison \
      cython python3-dev python-setuptools

# Compile nghttp2 with serial number patch
ARG NGHTTP_VERSION
WORKDIR /tmp
RUN git clone https://github.com/nghttp2/nghttp2.git
WORKDIR /tmp/nghttp2
RUN git checkout v${NGHTTP_VERSION}

RUN git submodule update --init
COPY /src/magma/orc8r/cloud/docker/proxy/0001-Fix-get_x509_serial-for-long-serial-numbers.patch /tmp
RUN git apply /tmp/0001-Fix-get_x509_serial-for-long-serial-numbers.patch
RUN autoreconf -i
RUN automake
RUN autoconf
RUN ./configure --with-mruby --enable-app --disable-examples --disable-python-bindings
RUN make
RUN make install DESTDIR=/tmp/nghttp2/install

FROM ubuntu:bionic

ARG CNTLR_FILES=src/magma/orc8r/cloud/docker/controller
ARG PROXY_FILES=src/magma/orc8r/cloud/docker/proxy
ARG NGHTTP_VERSION

# Runtime dependencies for nghttpx, supervisord, etc
RUN apt-get update && \
    apt-get install -y \
        daemontools netcat openssl supervisor wget unzip \
        libssl-dev libev-dev libevent-dev libjansson-dev libjemalloc-dev libc-ares-dev \
        python3-pip \
        squid

# Install nghttpx
COPY --from=builder /tmp/nghttp2/install/usr /usr
RUN ldconfig

# Install python3 deps from pip
RUN pip3 install PyYAML jinja2

# Create an empty envdir for overriding in production
RUN mkdir -p /var/opt/magma/envdir

# Copy the configs
COPY configs /etc/magma/configs

# Copy proxy scripts and configs from the context
COPY ${PROXY_FILES}/templates /etc/magma/templates
COPY ${PROXY_FILES}/magma_headers.rb /etc/nghttpx/magma_headers.rb
COPY ${PROXY_FILES}/run_nghttpx.py /usr/local/bin/run_nghttpx.py
COPY ${PROXY_FILES}/squid.conf /etc/squid/squid.conf
COPY ${PROXY_FILES}/create_test_proxy_certs /usr/local/bin/create_test_proxy_certs

# Copy the supervisor configs
COPY ${PROXY_FILES}/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY ${CNTLR_FILES}/supervisor_logger.py /usr/local/lib/python2.7/dist-packages/supervisor_logger.py
CMD ["/usr/bin/supervisord"]

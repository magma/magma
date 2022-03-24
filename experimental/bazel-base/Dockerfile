################################################################
# Builder Image (can also be used as developer's image)
################################################################
FROM ubuntu:focal as bazel_builder

ENV TZ=America/New_York \
    DEBIAN_FRONTEND=noninteractive

#  MAGMA_ROOT is needed by python tests (e.g. freedomfi_one_tests in enodebd)
ENV MAGMA_ROOT=/magma

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
        apt-transport-https \
        apt-utils \
        build-essential \
        ca-certificates \
        curl \
        gcc \
        git \
        gnupg2 \
        g++ \
        # dependency of mobilityd (tests)
        iproute2 \
        # dependency of python services (e.g. magmad)
        iputils-ping \
        # dependencies of FreeDiameter
        bison \
        flex \
        libgcrypt-dev \
        libgmp3-dev \
        libidn11-dev \
        libsctp1 \
        libsctp-dev \
        # dependencies of oai/mme
        libczmq-dev \
        libconfig-dev \
        # dependency of @sentry_native//:sentry
        libcurl4-openssl-dev \
        # dependency of sctpd
        libsctp-dev \
        libssl-dev \
        # dependency of pip systemd
        libsystemd-dev \
        lld \
        # dependency of python services (e.g. magmad)
        net-tools \
        # dependency of python services (e.g. pipelined)
        netbase \
        python3 \
        python-is-python3 \
        # dependency of python services (e.g. magmad)
        systemd \
        unzip \
        # dependency of liagent
        uuid-dev \
        vim \
        wget \
        zip

# Install depencies from facebook repository
COPY orc8r/tools/ansible/roles/pkgrepo/files/jfrog.pub /tmp/
RUN apt-key add /tmp/jfrog.pub && \
    echo "deb https://facebookconnectivity.jfrog.io/artifactory/list/dev-focal/ focal main" > /etc/apt/sources.list.d/fbc.list && \
    apt-get update -y && \
    apt-get install -y bcc-tools

# Install bazel
WORKDIR /usr/sbin
RUN wget --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-amd64 && \
    chmod +x bazelisk-linux-amd64 && \
    ln -s /usr/sbin/bazelisk-linux-amd64 /usr/sbin/bazel

# Install Folly as a static library in the container
RUN apt-get install -y --no-install-recommends cmake

## Install Fmt (Folly Dep)
RUN git clone https://github.com/fmtlib/fmt.git && cd fmt && \
    mkdir _build && cd _build && \
    cmake .. -DFMT_TEST=0 && \
    make -j"$(nproc)" && \
    make install && \
    cd / && \
    rm -rf fmt

RUN apt-get install -y --no-install-recommends \
    libgoogle-glog-dev \
    libgflags-dev \
    libboost-all-dev \
    libevent-dev \
    libdouble-conversion-dev \
    libiberty-dev

# Facebook Folly C++ lib
# Note: "Because folly does not provide any ABI compatibility guarantees from
#        commit to commit, we generally recommend building folly as a static library."
# Here we checkout the hash for v2021.02.22.00 (arbitrary recent version)
RUN git clone https://github.com/facebook/folly && cd folly && \
    git checkout tags/v2021.02.15.00 && \
    mkdir _build && cd _build && \
    cmake .. && \
    make -j"$(nproc)" && \
    make install && \
    cd / && \
    rm -rf folly

RUN apt-get install -y --no-install-recommends \
    # used by libtins and connection tracker
    libpcap-dev=1.9.1-3 \
    # used by connection tracker
    libmnl-dev=1.0.4-2

RUN apt-get -y install --no-install-recommends \
    libtool=2.4.6-14

##### liblfds
# https://www.liblfds.org/mediawiki/index.php?title=r7.1.0:Building_Guide_(liblfds)
RUN git clone https://github.com/liblfds/liblfds.git && \
    cd liblfds/liblfds/liblfds7.1.0/liblfds710/build/gcc_gnumake/ && \
    make -j"$(nproc)" && \
    make ar_install && \
    cd / && \
    rm -rf liblfds

# TODO(GH9710): Generate asn1c with Bazel - also this repo is really old :o
RUN git clone https://gitlab.eurecom.fr/oai/asn1c.git && \
    cd asn1c && \
    git checkout f12568d617dbf48497588f8e227d70388fa217c9 && \
    autoreconf -iv && \
    ./configure && \
    make -j"$(nproc)" && \
    make install

##### NETTLE and GNUTLS
# TODO Upgrade these - requires us to update our use of libnettle due to API migration.
#  see https://gist.github.com/electronjoe/a899e4bfbc2904cb353444386296c38e
# Note the CFLAGS define below due to glibc deprecation of critical flag,
#  see https://github.com/rdslw/openwrt/blob/e5d47f32131849a69a9267de51a30d6be1f0d0ac/tools/bison/patches/110-glibc-change-work-around.patch
RUN wget --progress=dot:giga http://ftp.ntua.gr/mirror/gnu/nettle/nettle-2.5.tar.gz && \
    tar -xf nettle-2.5.tar.gz && \
    cd nettle-2.5 && \
    mkdir build && \
    cd build/ && \
    ../configure --disable-openssl --enable-shared --libdir=/usr/lib && \
    make -j"$(nproc)" && \
    make install && \
    ldconfig -v && \
    cd / && \
    wget --progress=dot:giga https://www.gnupg.org/ftp/gcrypt/gnutls/v3.1/gnutls-3.1.23.tar.xz && \
    tar xf gnutls-3.1.23.tar.xz && \
    cd gnutls-3.1.23 && \
    mkdir build && cd build && \
    CFLAGS=-D_IO_ftrylockfile ../configure --libdir=/usr/lib --with-libnettle-prefix=/usr && \
    make -j"$(nproc)" && \
    make install && \
    ldconfig -v && \
    cd / && \
    rm -rf nettle* && \
    rm -rf gnutls*

##### FreeDiameter
COPY lte/gateway/c/core/oai/patches/ /tmp/
RUN git clone https://github.com/OPENAIRINTERFACE/opencord.org.freeDiameter.git freediameter && \
    cd freediameter && \
    patch -p1 < /tmp/0001-opencoord.org.freeDiameter.patch && \
    mkdir build && \
    cd build && \
    cmake -DDISABLE_SCTP:BOOL=ON ../ && \
    make -j"$(nproc)" && \
    make install && \
    rm -rf /freediameter /tmp/*

# Update shared library configuration
RUN ldconfig -v

RUN ln -s /magma/bazel/bazelrcs/docker.bazelrc /etc/bazelrc

WORKDIR /magma

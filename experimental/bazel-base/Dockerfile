################################################################
# Builder Image (can also be used as developer's image)
################################################################
FROM ubuntu:focal as bazel_builder

ENV TZ=America/New_York \
    DEBIAN_FRONTEND=noninteractive

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

# TODO(GH9710): Generate asn1c with Bazel - also this repo is really old :o
RUN git clone https://gitlab.eurecom.fr/oai/asn1c.git && \
    cd asn1c && \
    git checkout f12568d617dbf48497588f8e227d70388fa217c9 && \
    autoreconf -iv && \
    ./configure && \
    make -j"$(nproc)" && \
    make install

# Update shared library configuration
RUN ldconfig -v

RUN ln -s /magma/bazel/bazelrcs/docker.bazelrc /etc/bazelrc

WORKDIR /magma

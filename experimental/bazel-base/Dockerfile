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
        # dependency of @sentry_native//:sentry
        libcurl4-openssl-dev \
        gcc \
        git \
        gnupg2 \
        g++ \
        python-dev \
        zip \
        unzip \
        vim \
        wget \
        libssl-dev

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

# libtins is required to build the connection_tracker
RUN apt-get install -y --no-install-recommends \
    cmake=3.16.3-1ubuntu1 \
    libpcap-dev=1.9.1-3 \
    libmnl-dev=1.0.4-2

# dependency of @sentry_native//:sentry
RUN apt-get -y install --no-install-recommends \
    libcurl4-openssl-dev=7.68.0-1ubuntu2.7

# TODO(@themarwhal): this might not be hard to bazelify if we can generate the config
RUN git clone --branch v4.3 https://github.com/mfontanini/libtins.git && \
    cd libtins && \
    mkdir build && \
    cd build && \
    cmake ../ -DLIBTINS_ENABLE_CXX11=1 && \
    make -j"$(nproc)" && \
    make install && \
    cd / && \
    rm -rf libtins && \
    # symlink to /usr/lib to match Magma VM
    ln -s /usr/local/lib/libtins.so /usr/lib/libtins.so

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

RUN ln -s /magma/experimental/bazel-base/bazelrcs/docker.bazelrc /etc/bazelrc

WORKDIR /magma

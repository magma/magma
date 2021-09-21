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
    cmake .. && \
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

##### Facebook Folly C++ lib
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


WORKDIR /magma

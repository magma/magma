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
        # TODO(@themarwhal) Remove these Folly dependencies below - using these for now to get SessionD to build
        libboost-all-dev \
        libevent-dev \
        libdouble-conversion-dev \
        # libgoogle-glog-dev see folly BUILD.bazel file for how to specifiy custom paths
        # libgflags-dev
        libiberty-dev \
        liblz4-dev \
        liblzma-dev \
        libsnappy-dev \
        zlib1g-dev \
        binutils-dev \
        libjemalloc-dev \
        libssl-dev \
        pkg-config \
        libunwind-dev


WORKDIR /usr/sbin
RUN wget --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-amd64 && \
    chmod +x bazelisk-linux-amd64 && \
    ln -s /usr/sbin/bazelisk-linux-amd64 /usr/sbin/bazel

WORKDIR /magma

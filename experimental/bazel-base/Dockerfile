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
        # TODO(@themarwhal) Remove these Folly dependencies below - GH8577
        # The full list of recommended dependencies are listed at https://github.com/facebook/folly/blob/master/README.md#ubuntu-1604-lts
        # Omitting the following as they are not needed
        # - Already built by bazel: libgoogle-glog-dev, libgflags-dev
        # - Not used in Magma: ibjemalloc-dev, pkg-config, libunwind-dev, liblz4-dev
        #                      liblzma-dev, libsnappy-dev, zlib1g-dev, binutils-dev
        libboost-all-dev \
        libevent-dev \
        libdouble-conversion-dev \
        libiberty-dev \
        libssl-dev 


WORKDIR /usr/sbin
RUN wget --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-amd64 && \
    chmod +x bazelisk-linux-amd64 && \
    ln -s /usr/sbin/bazelisk-linux-amd64 /usr/sbin/bazel

WORKDIR /magma

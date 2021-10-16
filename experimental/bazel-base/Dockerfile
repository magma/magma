################################################################
# Builder Image (can also be used as developer's image)
################################################################
FROM ubuntu:focal as bazel_builder

ENV TZ=America/New_York \
    DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get upgrade -y 

RUN apt-get install -y --no-install-recommends \
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
        vim \
        wget

# Install bazel
WORKDIR /usr/sbin
RUN wget --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-amd64 && \
    chmod +x bazelisk-linux-amd64 && \
    ln -s /usr/sbin/bazelisk-linux-amd64 /usr/sbin/bazel

# Dependencies for rules_folly @ https://github.com/storypku/rules_folly
RUN apt-get -y install --no-install-recommends \
    autoconf=2.69-11.1 \
    automake=1:1.16.1-4ubuntu6 \
    libtool=2.4.6-14 \
    libssl-dev=1.1.1f-1ubuntu2.8

# dependency of @sentry_native//:sentry
RUN apt-get -y install --no-install-recommends libcurl4-openssl-dev

# libtins is required to build the connection_tracker
RUN apt-get install -y --no-install-recommends \
    cmake=3.16.3-1ubuntu1 \
    libpcap-dev=1.9.1-3 \
    libmnl-dev=1.0.4-2

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

# Update shared library configuration
RUN ldconfig -v

RUN ln -s /magma/experimental/bazel-base/bazelrcs/docker.bazelrc /etc/bazelrc

WORKDIR /magma

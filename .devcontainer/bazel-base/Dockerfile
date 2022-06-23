################################################################
# Builder Image (can also be used as developer's image)
################################################################
FROM ubuntu:focal as bazel_builder

ARG DEB_PORT=amd64

ARG PYTHON_VERSION=3.8

ENV TZ=UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

#  MAGMA_ROOT is needed by python tests (e.g. freedomfi_one_tests in enodebd)
ENV MAGMA_ROOT=/workspaces/magma

RUN echo "Install general purpose packages" && \
    apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
        apt-transport-https \
        apt-utils \
        # dependencies of FreeDiameter
        bison \
        build-essential \
        ca-certificates \
        cmake \
        curl \
        gcc \
        git \
        gnupg2 \
        g++ \
        # dependency of mobilityd (tests)
        iproute2 \
        # dependency of python services (e.g. magmad)
        iputils-ping \
        flex \
        libconfig-dev \
        # dependency of @sentry_native//:sentry
        libcurl4-openssl-dev \
        # dependencies of oai/mme
        libczmq-dev \
        libgcrypt-dev \
        libgmp3-dev \
        libidn11-dev \
        libsctp1 \
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
        python${PYTHON_VERSION} \
        python-is-python3 \
        software-properties-common \
        # dependency of python services (e.g. magmad)
        systemd \
        unzip \
        # dependency of liagent
        uuid-dev \
        vim \
        wget \
        zip

RUN apt-get install -y --no-install-recommends \
        libtool=2.4.6-14 && \
    echo "Install Folly dependencies" && \
    apt-get install -y --no-install-recommends \
        libgoogle-glog-dev \
        libgflags-dev \
        libboost-all-dev \
        libevent-dev \
        libdouble-conversion-dev \
        libboost-chrono-dev \
        libiberty-dev && \
    echo "Install libtins and connection tracker dependencies" && \
    apt-get install -y --no-install-recommends \
        libpcap-dev=1.9.1-3 \
        libmnl-dev=1.0.4-2

## Install Fmt (Folly Dep)
RUN git clone https://github.com/fmtlib/fmt.git && \
    cd fmt && \
    mkdir _build && \
    cd _build && \
    cmake -DBUILD_SHARED_LIBS=ON -DFMT_TEST=0 .. && \
    make -j"$(nproc)" && \
    make install && \
    cd / && \
    rm -rf fmt

# Facebook Folly C++ lib
# Note: "Because folly does not provide any ABI compatibility guarantees from
#        commit to commit, we generally recommend building folly as a static library."
# Here we checkout the hash for v2021.02.22.00 (arbitrary recent version)
RUN git clone --depth 1 --branch v2021.02.15.00 https://github.com/facebook/folly && \
    cd /folly && \
    mkdir _build && \
    cd _build && \
    cmake -DBUILD_SHARED_LIBS=ON .. && \
    make -j"$(nproc)" && \
    make install && \
    cd / && \
    rm -rf folly

# setup magma artifactories and install magma dependencies
RUN wget -qO - https://artifactory.magmacore.org:443/artifactory/api/gpg/key/public | apt-key add - && \
    add-apt-repository 'deb https://artifactory.magmacore.org/artifactory/debian-test focal-ci main' && \
    add-apt-repository 'deb https://artifactory.magmacore.org/artifactory/debian-test focal-1.7.0 main' && \
    apt-get update -y && \
    apt-get install -y --no-install-recommends \
        bcc-tools \
        liblfds710 \
        oai-asn1c \
        oai-gnutls \
        oai-nettle \
        oai-freediameter

# Install bazel
WORKDIR /usr/sbin
RUN wget --progress=dot:giga https://github.com/bazelbuild/bazelisk/releases/download/v1.10.0/bazelisk-linux-"${DEB_PORT}" && \
    chmod +x bazelisk-linux-"${DEB_PORT}" && \
    ln -s /usr/sbin/bazelisk-linux-"${DEB_PORT}" /usr/sbin/bazel

# Update shared library configuration
RUN ldconfig -v

RUN ln -s /workspaces/magma/.bazel-cache /var/cache/bazel-cache
RUN ln -s /workspaces/magma/.bazel-cache-repo /var/cache/bazel-cache-repo

WORKDIR /workspaces/magma

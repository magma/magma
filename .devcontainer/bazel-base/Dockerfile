################################################################
# Builder Image (can also be used as developer's image)
################################################################
FROM ubuntu:focal as bazel_builder

ARG DEB_PORT=amd64

ARG PYTHON_VERSION=3.8

ARG USER_UID=1000
ARG USER_GID=$USER_UID

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
        bison \
        build-essential \
        ca-certificates \
        cmake \
        curl \
        gcc \
        git \
        gnupg2 \
        g++ \
        iproute2 `# dependency of mobilityd (tests)` \
        iputils-ping `# dependency of python services (e.g. magmad)` \
        flex \
        libconfig-dev \
        libcurl4-openssl-dev `# dependency of @sentry_native//:sentry` \
        libczmq-dev \
        libgcrypt-dev \
        libgmp3-dev \
        libidn11-dev \
        libsctp1 \
        libsqlite3-dev \
        libsctp-dev `# dependency of sctpd` \
        libssl-dev \
        libsystemd-dev `# dependency of pip systemd` \
        lld-12 \
        net-tools `# dependency of python services (e.g. magmad)` \
        netbase `# dependency of python services (e.g. pipelined)` \
        python${PYTHON_VERSION} \
        python-is-python3 \
        python3-distutils `# dependency of bazel pip_parse rule` \
        software-properties-common \
        systemd `# dependency of python services (e.g. magmad)` \
        unzip \
        uuid-dev `# dependency of liagent` \
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

# setup magma artifactories and install magma dependencies
RUN wget -qO - https://artifactory.magmacore.org:443/artifactory/api/gpg/key/public | apt-key add - && \
    add-apt-repository 'deb https://artifactory.magmacore.org/artifactory/debian-test focal-ci main' && \
    apt-get update -y && \
    apt-get install -y --no-install-recommends \
        bcc-tools \
        libfolly-dev \
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

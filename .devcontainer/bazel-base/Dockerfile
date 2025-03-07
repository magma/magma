# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

################################################################
# Builder Image (can also be used as developer's image)
################################################################
FROM ubuntu:focal AS bazel_builder

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
COPY keys/linux_foundation_registry_key.asc /etc/apt/trusted.gpg.d/magma.asc
RUN add-apt-repository 'deb https://linuxfoundation.jfrog.io/artifactory/magma-packages-test focal-ci main' && \
    apt-get update -y && \
    apt-get install -y --no-install-recommends \
        bcc-tools \
        libfolly-dev \
        liblfds710 \
        magma-libfluid \
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

# Convenience alias to build AGW services
RUN echo 'alias magma-build-agw="bazel build //... --build_tag_filters=service,util_script,-mme_oai"' >> ~/.bashrc

WORKDIR /workspaces/magma

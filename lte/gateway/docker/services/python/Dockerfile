################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

# -----------------------------------------------------------------------------
# Builder image for Python binaries and Magma proto files
# -----------------------------------------------------------------------------
ARG CPU_ARCH=x86_64
ARG DEB_PORT=amd64
ARG OS_DIST=ubuntu
ARG OS_RELEASE=focal
ARG EXTRA_REPO=https://artifactory.magmacore.org/artifactory/debian-test

FROM $OS_DIST:$OS_RELEASE AS builder
ARG CPU_ARCH
ARG DEB_PORT
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO

ENV MAGMA_DEV_MODE 0
ENV TZ=Europe/Paris
ENV MAGMA_ROOT=/magma
ENV PYTHON_BUILD=/build
ENV PIP_CACHE_HOME="~/.pipcache"
ENV SWAGGER_CODEGEN_DIR /var/tmp/codegen/modules/swagger-codegen-cli/target
ENV SWAGGER_CODEGEN_JAR ${SWAGGER_CODEGEN_DIR}/swagger-codegen-cli.jar
ARG CODEGEN_VERSION=2.2.3
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
  docker.io \
  git \
  lsb-release \
  libsystemd-dev \
  libprotobuf-dev \
  pkg-config \
  python3-dev \
  python3-eventlet \
  python3-pip \
  python3-protobuf \
  python3-pystemd \
  ruby \
  ruby-dev \
  sudo \
  openjdk-8-jdk \
  openjdk-8-jre-headless \
  virtualenv \
  wget \
  && rm -rf /var/lib/apt/lists/*

RUN gem install fpm

#COPY . $MAGMA_ROOT/
COPY ./lte/gateway/python/magma $MAGMA_ROOT/lte/gateway/python/magma
COPY ./third_party/build/bin/aioeventlet_build.sh $MAGMA_ROOT/third_party/build/bin/aioeventlet_build.sh
COPY ./third_party/build/lib/util.sh $MAGMA_ROOT/third_party/build/lib/util.sh
WORKDIR /var/tmp/
RUN /magma/third_party/build/bin/aioeventlet_build.sh && \
  dpkg -i python3-aioeventlet*

RUN mkdir -p ${SWAGGER_CODEGEN_DIR}; \
  wget --no-verbose https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${CODEGEN_VERSION}/swagger-codegen-cli-${CODEGEN_VERSION}.jar -O ${SWAGGER_CODEGEN_JAR}

WORKDIR /magma/lte/gateway/python
COPY \
    ./lte/gateway/python/Makefile \
    ./lte/gateway/python/defs.mk \
    ./lte/gateway/python/setup.py \
    $MAGMA_ROOT/lte/gateway/python/
COPY ./lte/gateway/python/load_tests $MAGMA_ROOT/lte/gateway/python/load_tests
COPY ./lte/gateway/python/scripts $MAGMA_ROOT/lte/gateway/python/scripts
COPY ./orc8r/gateway/python/magma $MAGMA_ROOT/orc8r/gateway/python/magma
COPY \
    ./orc8r/gateway/python/python.mk \
    ./orc8r/gateway/python/Makefile \
    ./orc8r/gateway/python/defs.mk \
    ./orc8r/gateway/python/setup.py \
     $MAGMA_ROOT/orc8r/gateway/python/
COPY ./orc8r/gateway/python/scripts $MAGMA_ROOT/orc8r/gateway/python/scripts
COPY ./protos $MAGMA_ROOT/protos
COPY ./lte/protos $MAGMA_ROOT/lte/protos
COPY ./orc8r/protos $MAGMA_ROOT/orc8r/protos
COPY ./feg/protos $MAGMA_ROOT/feg/protos
COPY ./dp/protos $MAGMA_ROOT/dp/protos
COPY ./lte/swagger $MAGMA_ROOT/lte/swagger
COPY ./orc8r/swagger $MAGMA_ROOT/orc8r/swagger
RUN make buildenv

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM $OS_DIST:$OS_RELEASE AS gateway_python
ARG OS_DIST
ARG CPU_ARCH
ARG OS_RELEASE
ARG EXTRA_REPO

ENV VIRTUAL_ENV=/build
ENV TZ=Europe/Paris

ARG JSONPOINTER_VERSION=1.13
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
  apt-transport-https \
  ca-certificates \
  docker.io \
  ethtool \
  inetutils-ping \
  iproute2 \
  iptables \
  isc-dhcp-client \
  linux-headers-generic \
  net-tools \
  netcat \
  nghttp2-proxy \
  python3-pip \
  python3-venv \
  python3-eventlet \
  python3-pystemd \
  python3-jinja2 \
  redis-server \
  sudo \
  virtualenv \
  wget \
  && rm -rf /var/lib/apt/lists/*

RUN python3 -m venv $VIRTUAL_ENV

ENV PATH="/magma/orc8r/gateway/python/scripts/:/magma/lte/gateway/python/scripts/:$VIRTUAL_ENV/bin:$PATH"

RUN echo "deb https://artifactory.magmacore.org/artifactory/debian-test focal-ci main" > /etc/apt/sources.list.d/magma.list
RUN wget -qO - https://artifactory.magmacore.org:443/artifactory/api/gpg/key/public | apt-key add -

RUN echo "deb https://packages.fluentbit.io/ubuntu/focal focal main" > /etc/apt/sources.list.d/tda.list
RUN wget -qO - https://packages.fluentbit.io/fluentbit.key | apt-key add -

RUN apt-get update && apt-get install -y \
  bcc-tools \
  libopenvswitch \
  openvswitch-datapath-dkms \
  openvswitch-common \
  openvswitch-switch \
  td-agent-bit \
  wireguard \
  && rm -rf /var/lib/apt/lists/*

COPY ./orc8r/gateway/python/scripts/ /usr/local/bin
COPY ./lte/gateway/python/scripts/ /usr/local/bin
COPY ./lte/gateway/configs/templates /etc/magma/templates/
COPY ./orc8r/gateway/configs/templates/nghttpx.conf.template /etc/magma/templates/nghttpx.conf.template
COPY ./lte/gateway/deploy/roles/magma/files/set_irq_affinity /usr/local/bin/set_irq_affinity

COPY --from=builder /build /build
COPY --from=builder /magma/lte/gateway/python/magma /magma/lte/gateway/python/magma
COPY --from=builder /magma/orc8r/gateway/python/magma /magma/orc8r/gateway/python/magma
COPY --from=builder /var/tmp/python3-aioeventlet* /var/tmp/

RUN chmod -R +x /usr/local/bin/generate* /usr/local/bin/set_irq_affinity /usr/local/bin/checkin_cli.py && \
  dpkg -i /var/tmp/python3-aioeventlet* && \
  pip install jsonpointer>$JSONPOINTER_VERSION && \
  mkdir -p /var/opt/magma/




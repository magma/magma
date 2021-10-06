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
# Builder image for C binaries and Magma proto files
# -----------------------------------------------------------------------------
ARG OS_DIST=ubuntu
ARG OS_RELEASE=focal
ARG EXTRA_REPO=https://facebookconnectivity.jfrog.io/artifactory/list/dev-focal/

FROM $OS_DIST:$OS_RELEASE AS builder
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO

ENV MAGMA_DEV_MODE 0
ENV TZ=Europe/Paris
ENV MAGMA_ROOT=/magma
ENV PYTHON_BUILD=/build
ENV PIP_CACHE_HOME="~/.pipcache"
ENV SWAGGER_CODEGEN_JAR=/var/tmp/codegen/modules/swagger-codegen-cli/target/swagger-codegen-cli.jar

ARG MVN_VERSION=3.5.4
ARG CODEGEN_VERSION=v2.2.3
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
  wget \
  ruby \
  sudo \
  ruby-dev \
  python3-pip \
  python3-dev \
  python3-eventlet \
  python3-pystemd \
  python3-protobuf \
  git \
  virtualenv \
  lsb-release \
  openjdk-8-jre-headless \
  openjdk-8-jdk \
  pkg-config \
  libsystemd-dev \
  libprotobuf-dev

RUN cd /usr/local/bin && ln -s /usr/bin/python3 python
RUN gem install fpm

COPY . $MAGMA_ROOT/

WORKDIR /var/tmp/
RUN /magma/third_party/build/bin/aioeventlet_build.sh && \
    dpkg -i python3-aioeventlet*

RUN wget http://mirrors.ibiblio.org/apache/maven/maven-3/$MVN_VERSION/binaries/apache-maven-$MVN_VERSION-bin.tar.gz && \
    tar -xvf apache-maven-$MVN_VERSION-bin.tar.gz

RUN git clone https://github.com/swagger-api/swagger-codegen /var/tmp/codegen && \
    cd /var/tmp/codegen && \
    git checkout $CODEGEN_VERSION

WORKDIR /var/tmp/codegen
RUN ../apache-maven-$MVN_VERSION/bin/mvn clean package -DskipTests

WORKDIR /magma/lte/gateway/python
RUN make buildenv

# -----------------------------------------------------------------------------
# Dev/Production image
# -----------------------------------------------------------------------------
FROM $OS_DIST:$OS_RELEASE AS gateway_python
ARG OS_DIST
ARG OS_RELEASE
ARG EXTRA_REPO

ENV VIRTUAL_ENV=/build
ENV TZ=Europe/Paris

ARG JSONPOINTER_VERSION=1.13
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
  apt-transport-https \
  ca-certificates \
  netcat \
  sudo \
  python3-pip \
  python3-venv \
  virtualenv \
  python3-eventlet \
  python3-pystemd \
  python3-jinja2 \
  nghttp2-proxy \
  net-tools \
  inetutils-ping \
  redis-server \
  wget \
  ethtool \
  linux-headers-`uname -r` \
  iptables \
  iproute2

RUN cd /usr/local/bin && ln -s /usr/bin/python3 pytho && \
  python3 -m venv $VIRTUAL_ENV

ENV PATH="/magma/orc8r/gateway/python/scripts/:/magma/lte/gateway/python/scripts/:$VIRTUAL_ENV/bin:$PATH"

RUN echo "deb https://facebookconnectivity.jfrog.io/artifactory/list/dev-focal/ focal main" > /etc/apt/sources.list.d/magma.list
RUN wget -qO - https://facebookconnectivity.jfrog.io/artifactory/api/security/keypair/default-gpg-key/public | apt-key add -

RUN echo "deb https://packages.fluentbit.io/ubuntu/focal focal main" > /etc/apt/sources.list.d/tda.list
RUN wget -qO - https://packages.fluentbit.io/fluentbit.key | apt-key add -

RUN apt-get update && apt-get install -y \
  td-agent-bit \
  libopenvswitch \
  openvswitch-datapath-dkms \
  openvswitch-common \
  openvswitch-switch

COPY --from=builder /build /build
COPY --from=builder /magma /magma
COPY --from=builder /magma/orc8r/gateway/python/scripts/ /usr/local/bin
COPY --from=builder /magma/lte/gateway/python/scripts/ /usr/local/bin
COPY --from=builder /var/tmp/python3-aioeventlet* /var/tmp/
COPY --from=builder /magma/lte/gateway/configs/templates /etc/magma/templates/
COPY --from=builder /magma/orc8r/gateway/configs/templates/nghttpx.conf.template /etc/magma/templates/nghttpx.conf.template
COPY --from=builder /magma/orc8r/gateway/python/scripts/generate_nghttpx_config.py /usr/local/bin/generate_nghttpx_config.py
COPY --from=builder /magma/orc8r/gateway/python/scripts/generate_service_config.py /usr/local/bin/generate_service_config.py
COPY --from=builder /magma/orc8r/gateway/python/scripts/generate_fluent_bit_config.py /usr/local/bin/generate_fluent_bit_config.py
COPY --from=builder /magma/lte/gateway/deploy/roles/magma/files/set_irq_affinity /usr/local/bin/set_irq_affinity

RUN chmod -R +x /usr/local/bin/generate* /usr/local/bin/set_irq_affinity /usr/local/bin/checkin_cli.py && \
  dpkg -i /var/tmp/python3-aioeventlet* && \
  pip install jsonpointer>$JSONPOINTER_VERSION && \
  mkdir -p /var/opt/magma/

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
ARG OS_DIST=ubuntu
ARG OS_RELEASE=focal


FROM $OS_DIST:$OS_RELEASE AS builder

ARG MAGMA_RELEASE=master
ARG DEBIAN_FRONTEND=noninteractive
ARG CODEGEN_VERSION=v2.2.3
ARG MVN_VERSION=3.5.4

ENV VIRTUAL_ENV=/build/venv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"
ENV MAGMA_ROOT=/magma
ENV S1AP_TESTER_ROOT=/s1ap-tester
ENV S1AP_TESTER_SRC=/S1APTester
ENV PYTHON_BUILD=/build
ENV PIP_CACHE_HOME="~/.pipcache"
ENV TZ=Europe/Paris
ENV SWAGGER_CODEGEN_JAR=/var/tmp/codegen/modules/swagger-codegen-cli/target/swagger-codegen-cli.jar
ENV S1AP_TESTER_CFG=$MAGMA_ROOT/lte/gateway/python/integ_tests/data/s1ap_tester_cfg

# ENV PYTHON_VERSION=3.8

RUN apt-get update && apt-get install -y \
  git \
  python3-pip \
  python3-venv \
  python3-protobuf \
  python3-eventlet \
  python3-pystemd \
  libsystemd-dev \
  virtualenv \
  lsb-release \
  apt-utils \
  sudo \
  libprotobuf-dev \
  openjdk-8-jre-headless \
  openjdk-8-jdk \
  ruby \
  ruby-dev \
  wget \
  pkg-config \
  libmysqlclient-dev \
  rsync

RUN gem install fpm

RUN git clone https://github.com/magma/magma $MAGMA_ROOT
WORKDIR /magma
RUN git checkout $MAGMA_RELEASE

WORKDIR /var/tmp/
RUN /magma/third_party/build/bin/aioeventlet_build.sh && \
    dpkg -i python3-aioeventlet*

RUN wget http://mirrors.ibiblio.org/apache/maven/maven-3/$MVN_VERSION/binaries/apache-maven-$MVN_VERSION-bin.tar.gz && \
    tar -xvf apache-maven-$MVN_VERSION-bin.tar.gz

RUN git clone https://github.com/swagger-api/swagger-codegen /var/tmp/codegen
WORKDIR /var/tmp/codegen
RUN git checkout $CODEGEN_VERSION

RUN ../apache-maven-$MVN_VERSION/bin/mvn clean package -DskipTests

RUN python3 -m venv $VIRTUAL_ENV
RUN cd /usr/local/bin && ln -s /usr/bin/python3 python

WORKDIR $MAGMA_ROOT/lte/gateway/python
RUN make

RUN mkdir -p $S1AP_TESTER_ROOT/bin /var/tmp/test_results/

WORKDIR $MAGMA_ROOT/lte/gateway/python/integ_tests

RUN make

RUN $MAGMA_ROOT/lte/gateway/deploy/roles/magma_test/files/clone_s1_tester.sh

RUN touch /usr/include/stropts.h

RUN $MAGMA_ROOT/lte/gateway/deploy/roles/magma_test/files/build_s1_tester.sh

RUN pip3 install --cache-dir $PIP_CACHE_HOME pyparsing

RUN python3 $MAGMA_ROOT/lte/gateway/deploy/roles/magma_test/files/c_parser.py

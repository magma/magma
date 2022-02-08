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
ARG GIT_VERSION=1:2.25.1-1ubuntu3.2
ARG BUILD_ESSENTIAL_VERSION=12.8ubuntu1.1
ARG PIP_VERSION=20.0.2-5ubuntu1.6
ARG VENV_VERSION=3.8.2-0ubuntu2
ARG PROTOBUF_VERSION=3.6.1.3-2ubuntu5
ARG EVENTLET_VERSION=0.25.1-2ubuntu1.1
ARG PYSTEMD_VERSION=0.7.0-2build1
ARG LIBSYSTEMD_DEV_VERSION=245.4-4ubuntu3.15
ARG VIRTUALENV_VERSION=20.0.17-1ubuntu0.4
ARG LSB_RELEASE_VERSION=11.1.0ubuntu2
ARG APT_UTILS_VERSION=2.0.6
ARG SUDO_VERSION=1.8.31-1ubuntu1.2
ARG LIBPROTOBUF_VERSION=3.6.1.3-2ubuntu5
ARG JRE_VERSION=8u312-b07-0ubuntu1~20.04
ARG JDK_VERSION=8u312-b07-0ubuntu1~20.04
ARG RUBY_VERSION=1:2.7+1
ARG RUBY_DEV_VERSION=1:2.7+1
ARG WGET_VERSION=1.20.3-1ubuntu2
ARG PKG_CONFIG_VERSION=0.29.1-0ubuntu4
ARG LIBMYSQL_VERSION=8.0.28-0ubuntu0.20.04.3
ARG RSYNC_VERSION=3.1.3-8ubuntu0.1
ARG BINUTILS_VERSION=2.34-6ubuntu1.3

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

RUN apt-get update && apt-get install -y --no-install-recommends \
  git=$GIT_VERSION \
  build-essential=$BUILD_ESSENTIAL_VERSION \
  python3-pip=$PIP_VERSION \
  python3-venv=$VENV_VERSION \
  python3-protobuf=$PROTOBUF_VERSION \
  python3-eventlet=$EVENTLET_VERSION \
  python3-pystemd=$PYSTEMD_VERSION \
  libsystemd-dev=$LIBSYSTEMD_DEV_VERSION \
  virtualenv=$VIRTUALENV_VERSION \
  lsb-release=$LSB_RELEASE_VERSION \
  apt-utils=$APT_UTILS_VERSION \
  sudo=$SUDO_VERSION \
  libprotobuf-dev=$LIBPROTOBUF_VERSION \
  openjdk-8-jre-headless=$JRE_VERSION \
  openjdk-8-jdk=$JDK_VERSION \
  ruby=$RUBY_VERSION \
  ruby-dev=$RUBY_DEV_VERSION \
  wget=$WGET_VERSION \
  pkg-config=$PKG_CONFIG_VERSION \
  libmysqlclient-dev=$LIBMYSQL_VERSION \
  rsync=$RSYNC_VERSION \
  binutils=$BINUTILS_VERSION \
  && apt-get clean \
  && gem install fpm

RUN git clone https://github.com/magma/magma $MAGMA_ROOT && cd $MAGMA_ROOT && git checkout $MAGMA_RELEASE

WORKDIR /var/tmp/
RUN /magma/third_party/build/bin/aioeventlet_build.sh && dpkg -i python3-aioeventlet*

RUN wget http://mirrors.ibiblio.org/apache/maven/maven-3/$MVN_VERSION/binaries/apache-maven-$MVN_VERSION-bin.tar.gz && \
  tar -xvf apache-maven-$MVN_VERSION-bin.tar.gz && git clone https://github.com/swagger-api/swagger-codegen /var/tmp/codegen && \
  cd /var/tmp/codegen && git checkout $CODEGEN_VERSION && ../apache-maven-$MVN_VERSION/bin/mvn clean package -DskipTests
RUN python3 -m venv $VIRTUAL_ENV && cd /usr/local/bin && ln -s /usr/bin/python3 python

WORKDIR $MAGMA_ROOT/lte/gateway/python
RUN make && mkdir -p $S1AP_TESTER_ROOT/bin /var/tmp/test_results/
WORKDIR $MAGMA_ROOT/lte/gateway/python/integ_tests
RUN make && rm -rf /var/lib/apt/lists/* && $MAGMA_ROOT/lte/gateway/deploy/roles/magma_test/files/clone_s1_tester.sh && \
  touch /usr/include/stropts.h && $MAGMA_ROOT/lte/gateway/deploy/roles/magma_test/files/build_s1_tester.sh
RUN pip3 install --cache-dir $PIP_CACHE_HOME pyparsing && python3 $MAGMA_ROOT/lte/gateway/deploy/roles/magma_test/files/c_parser.py

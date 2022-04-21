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
ARG S1_IP=192.168.60.144
ARG S1_PORT=62462

ARG GIT_VERSION=2.25.1-1ubuntu3.2
ARG PYTHON3_PIP_VERSION=20.0.2-5ubuntu1.6
ARG PYTHON3_VENV_VERSION=3.8.2-0ubuntu2
ARG VIRTUALENV_VERSION=20.0.17-1ubuntu0.4
ARG ETHTOOL_VERSION=5.4-1
ARG IPERF3_VERSION=3.7-3

ENV MAGMA_ROOT=/magma

RUN apt-get update && apt-get install -y --no-install-recommends \
  git=$GIT_VERSION \
  python3-pip=$PYTHON3_PIP_VERSION \
  python3-venv=$PYTHON3_VENV_VERSION \
  virtualenv=$VIRTUALENV_VERSION \
  ethtool=$ETHTOOL_VERSION \
  iperf3=$IPERF3_VERSION \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

RUN pip3 install git+https://github.com/illu89/iperf3-python#egg=iperf3-python \
  pyroute2 \
  scapy-python3

RUN git clone https://github.com/magma/magma $MAGMA_ROOT && cd $MAGMA_ROOT && git checkout $MAGMA_RELEASE

RUN mkdir /usr/lib/python3.8/util

RUN cp $MAGMA_ROOT/lte/gateway/deploy/roles/trfserver/files/traffic_server.py /usr/local/bin/traffic_server.py
RUN cp $MAGMA_ROOT/lte/gateway/deploy/roles/trfserver/files/util/traffic_messages.py /usr/lib/python3.8/util/

ENTRYPOINT /usr/local/bin/traffic_server.py $S1_IP $S1_PORT

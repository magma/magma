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

ENV MAGMA_ROOT=/magma

RUN apt-get update && apt-get install -y --no-install-recommends \
  git \
  python3-pip \
  python3-venv \
  virtualenv \
  ethtool \
  iperf3 \
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

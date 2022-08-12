################################################################################
# Copyright 2022 The Magma Authors.

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

FROM $OS_DIST:$OS_RELEASE

ARG MAGMA_RELEASE=master
ARG DEBIAN_FRONTEND=noninteractive

ARG ETHTOOL_VERSION=1:5.4-1
ARG GIT_VERSION=1:2.25.1-1ubuntu3.5
ARG IPERF3_VERSION=3.7-3
ARG PYTHON_VERSION=3.8
ARG PYTHON3_PIP_VERSION=20.0.2-5ubuntu1.6
ARG PYTHON3_VENV_VERSION=3.8.2-0ubuntu2
ARG VIRTUALENV_VERSION=20.0.17-1ubuntu0.4

ENV MAGMA_ROOT=/magma

COPY . $MAGMA_ROOT/

RUN apt-get update && apt-get install -y --no-install-recommends \
  ethtool=$ETHTOOL_VERSION \
  git=$GIT_VERSION \
  iperf3=$IPERF3_VERSION \
  python3-pip=$PYTHON3_PIP_VERSION \
  python3-venv=$PYTHON3_VENV_VERSION \
  virtualenv=$VIRTUALENV_VERSION \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* && \
  pip3 install --no-cache-dir git+https://github.com/illu89/iperf3-python#egg=iperf3-python \
  pyroute2 \
  scapy-python3 && \
  mkdir -p /usr/lib/python${PYTHON_VERSION}/util && \
  cp $MAGMA_ROOT/lte/gateway/deploy/roles/trfserver/files/traffic_server.py /usr/local/bin/traffic_server.py && \
  cp $MAGMA_ROOT/lte/gateway/deploy/roles/trfserver/files/util/traffic_messages.py /usr/lib/python${PYTHON_VERSION}/util/

ENTRYPOINT ["/usr/local/bin/traffic_server.py"]
CMD ["192.168.60.144", "62462"]

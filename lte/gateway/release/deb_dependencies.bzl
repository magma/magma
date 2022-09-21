# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
External dependencies of the magma debian build.
"""

SCTPD_MIN_VERSION = "1.8.0"  # earliest version of sctpd with which this magma version is compatible

# Magma system dependencies: anything that we depend on at the top level, add
# here.
MAGMA_DEPS = [
    "grpc-dev (>= 1.15.0)",
    "lighttpd (>= 1.4.45)",
    "libxslt1.1",
    "nghttp2-proxy (>= 1.18.1)",
    "redis-server (>= 3.2.0)",
    "sudo",
    "dnsmasq (>= 2.7)",
    "net-tools",  # for ifconfig
    "python3-pip",
    "python3-apt",  # The version in pypi is abandoned and broken on stretch
    "libsystemd-dev",
    "libyaml-cpp-dev",  # install yaml parser
    "libgoogle-glog-dev",
    "python-redis",
    "magma-cpp-redis",
    "libfolly-dev",  # required for C++ services
    "libdouble-conversion-dev",  # required for folly
    "libboost-chrono-dev",  # required for folly
    "ntpdate",  # required for eventd time synchronization
    "tshark",  # required for call tracing
    "libtins-dev",  # required for Connection tracker
    "libmnl-dev",  # required for Connection tracker
    "getenvoy-envoy",  # for envoy dep
    "uuid-dev",  # for liagentd
    "libprotobuf17 (>= 3.0.0)",
    "nlohmann-json3-dev",
    "sentry-native",  # sessiond
    "td-agent-bit (>= 1.7.8)",
    # eBPF compile and load tools for kernsnoopd and AGW datapath
    # Ubuntu bcc lib (bpfcc-tools) is pretty old, use magma repo package
    "bcc-tools",
    "wireguard",
]

# OAI runtime dependencies
OAI_DEPS = [
    "libconfig9",
    "oai-asn1c",
    "oai-gnutls (>= 3.1.23)",
    "oai-nettle (>= 1.0.1)",
    "prometheus-cpp-dev (>= 1.0.2)",
    "liblfds710",
    "libsctp-dev",
    "magma-sctpd (>= {min_version})".format(min_version = SCTPD_MIN_VERSION),
    "libczmq-dev (>= 4.0.2-7)",
    "libasan5",
    "oai-freediameter (>= 0.0.2)",
]

# OVS runtime dependencies
OVS_DEPS = [
    "magma-libfluid (>= 0.1.0.7)",
    "libopenvswitch (>= 2.15.4-8)",
    "openvswitch-switch (>= 2.15.4-8)",
    "openvswitch-common (>= 2.15.4-8)",
    "openvswitch-datapath-dkms (>= 2.15.4-8)",
]

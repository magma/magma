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

SCTPD_MIN_VERSION = "1.9.0"  # earliest version of sctpd with which this magma version is compatible
DHCP_HELPER_CLI_MIN_VERSION = "1.9.0"  # earliest version of dhcp_helper_cli with which this magma version is compatible

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
    "magma-dhcp-cli (>= {min_version})".format(min_version = DHCP_HELPER_CLI_MIN_VERSION),
    "sentry-native",  # sessiond
    "td-agent-bit (>= 1.7.8)",
    # eBPF compile and load tools for kernsnoopd and AGW datapath
    # Ubuntu bcc lib (bpfcc-tools) is pretty old, use magma repo package
    "bcc-tools",
    "wireguard",
    "systemd",  # for postinstall script
    "psmisc",  # for postinstall script
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
    "libsystemd-dev",
]

# OVS runtime dependencies
OVS_DEPS = [
    "magma-libfluid (>= 0.1.0.7)",
    "libopenvswitch (>= 2.15.4-8)",
    "openvswitch-switch (>= 2.15.4-8)",
    "openvswitch-common (>= 2.15.4-8)",
    "openvswitch-datapath-dkms (>= 2.15.4-8)",
]

# Conflicts of the Bazel-built Magma package with upstream Ubuntu packages
# Actually the Magma package does not conflict with the official Ubuntu packages,
# but with our custom python packages created with pydep that go by the same name
# as some upstream packages.
MAGMA_CONFLICTS_UPSTREAM = [
    "python3-aiodns",
    "python3-aioeventlet",
    "python3-aiohttp",
    "python3-aiohttp-dbg",
    "python3-async-timeout",
    "python3-attr",
    "python3-bpfcc",
    "python3-certifi",
    "python3-cffi",
    "python3-cffi-backend",
    "python3-click",
    "python3-crypto",
    "python3-cryptography",
    "python3-dateutil",
    "python3-debtcollector",
    "python3-deprecated",
    "python3-django-memoize",
    "python3-dnspython",
    "python3-dpkt",
    "python3-docker",
    "python3-eventlet",
    "python3-fire",
    "python3-flask",
    "python3-greenlet",
    "python3-grpcio",
    "python3-h2",
    "python3-hiredis",
    "python3-hpack",
    "python3-hyperframe",
    "python3-idna",
    "python3-itsdangerous",
    "python3-jinja2",
    "python3-json-pointer",
    "python3-jsonpickle",
    "python3-jsonschema",
    "python3-lxml",
    "python3-markupsafe",
    "python3-markupsafe-dbg",
    "python3-memoize",
    "python3-msgpack",
    "python3-multidict",
    "python3-netaddr",
    "python3-netifaces",
    "python3-openvswitch",
    "python3-oslo.config",
    "python3-oslo.i18n",
    "python3-packaging",
    "python3-pbr",
    "python3-pkg-resources",
    "python3-priority",
    "python3-prometheus-client",
    "python3-protobuf",
    "python3-psutil",
    "python3-pycares",
    "python3-pycparser",
    "python3-pyroute2",
    "python3-pyrsistent",
    "python3-pystemd",
    "python3-redis",
    "python3-repoze.lru",
    "python3-requests",
    "python3-rfc3986",
    "python3-routes",
    "python3-ryu",
    "python3-scapy",
    "python3-sdnotify",
    "python3-sentry-sdk",
    "python3-setuptools",
    "python3-simplejson",
    "python3-six",
    "python3-sortedcontainers",
    "python3-stevedore",
    "python3-systemd",
    "python3-tinyrpc",
    "python3-tz",
    "python3-urllib3",
    "python3-webcolors",
    "python3-webob",
    "python3-websocket",
    "python3-werkzeug",
    "python3-wrapt",
    "python3-yaml",
    "python3-yarl",
    "python3-dpkt",
    "python3-sdnotify",
]

# Conflicts of the Magma 1.9 package with third party Ubuntu packages
MAGMA_CONFLICTS_THIRDPARTY = [
    "python3-pycryptodome",
    "python3-aiosignal",
    "python3-attrs",
    "python3-bravado-core",
    "python3-charset-normalizer",
    # python3-python-dateutil was built by us, while python3-dateutil was built by Ubuntu
    "python3-python-dateutil",
    "python3-envoy",
    "python3-frozenlist",
    "python3-jsonref",
    "python3-pymemoize",
    "python3-ovs",
    "python3-pyroute2",
    "python3-pyroute2.core",
    "python3-pyroute2.ndb",
    "python3-pyroute2.ethtool",
    "python3-pyroute2.ipdb",
    "python3-pyroute2.ipset",
    "python3-pyroute2.nftables",
    "python3-pyroute2.nslink",
    "python3-pytz",
    "python3-spyne",
    "python3-redis-collections",
    "python3-python-redis-lock",
    "python3-rfc3987",
    "python3-snowflake",
    "python3-strict-rfc3339",
    "python3-swagger-spec-validator",
    "python3-systemd-python",
    "python3-wsgiserver",
]

"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import os

from setuptools import setup

# We can use an environment variable to pass in the package version during
# build. Since we don't distribute this on its own, we don't really care what
# version this represents. 'None' defaults to 0.0.0.
VERSION = os.environ.get('PKG_VERSION', '0.0.0')

setup(
    name='lte',
    version=VERSION,
    packages=[
        'magma.enodebd',
        'magma.enodebd.data_models',
        'magma.enodebd.device_config',
        'magma.enodebd.devices',
        'magma.enodebd.devices.baicells_qrtb',
        'magma.enodebd.devices.experimental',
        'magma.enodebd.devices.freedomfi_one',
        'magma.enodebd.state_machines',
        'magma.enodebd.tr069',
        'magma.health',
        'magma.mobilityd',
        'magma.monitord',
        'magma.pipelined',
        'magma.pipelined.app',
        'magma.pipelined.ng_manager',
        'magma.pipelined.openflow',
        'magma.pipelined.qos',
        'magma.pipelined.ebpf',
        'magma.policydb',
        'magma.policydb.servicers',
        'magma.redirectd',
        'magma.redirectd.templates',
        'magma.smsd',
        'magma.health',
        'magma.subscriberdb',
        'magma.subscriberdb.crypto',
        'magma.subscriberdb.protocols',
        'magma.subscriberdb.protocols.diameter',
        'magma.subscriberdb.protocols.diameter.application',
        'magma.subscriberdb.store',
        'magma.subscriberdb.subscription',
        'magma.kernsnoopd',
        'load_tests',
    ],
    scripts=[
        'scripts/agw_health_cli.py',
        'scripts/config_stateless_agw.py',
        'scripts/create_oai_certs.py',
        'scripts/enodebd_cli.py',
        'scripts/fake_user.py',
        'scripts/feg_hello_cli.py',
        'scripts/generate_dnsd_config.py',
        'scripts/generate_oai_config.py',
        'scripts/ha_cli.py',
        'scripts/hello_cli.py',
        'scripts/mobility_cli.py',
        'scripts/mobility_dhcp_cli.py',
        'scripts/ocs_cli.py',
        'scripts/packet_tracer_cli.py',
        'scripts/packet_ryu_cli.py',
        'scripts/pcrf_cli.py',
        'scripts/pipelined_cli.py',
        'scripts/policydb_cli.py',
        'scripts/s6a_proxy_cli.py',
        'scripts/s6a_service_cli.py',
        'scripts/session_manager_cli.py',
        'scripts/sgs_cli.py',
        'scripts/sms_cli.py',
        'scripts/subscriber_cli.py',
        'scripts/spgw_service_cli.py',
        'scripts/cpe_monitoring_cli.py',
        'scripts/state_cli.py',
        'scripts/dp_probe_cli.py',
        'scripts/user_trace_cli.py',
        'scripts/icmpv6.py',
        'load_tests/loadtest_sessiond.py',
        'load_tests/loadtest_pipelined.py',
        'load_tests/loadtest_mobilityd.py',
        'load_tests/loadtest_subscriberdb.py',
    ],
    package_data={'magma.redirectd.templates': ['*.html']},
    install_requires=[
        'Cython>=0.29.32',
        'pystemd>=0.10.0',
        'fire>=0.4.0',
        'envoy>=0.0.3',
        'glob2>=0.7',
        # lxml required by spyne.
        'lxml==4.9.1',
        'ryu>=4.34',
        'spyne>=2.13,<2.14',
        'dpkt==1.9.8',
        'flask==1.1.4',
        'sentry_sdk>=1.5.0,<1.9',
        'aiodns>=3.0.0',
        'pymemoize>=1.0.3',
        'wsgiserver>=1.3',
        # pin recursive dependencies of ryu and others
        'charset-normalizer==2.0.0',
        'docker==4.1.0',
        'urllib3==1.26.11',
        'websocket-client>=1.3.2',
        'requests==2.28.1',
        'certifi>=2022.6.15',
        'idna==3.3',
        'python-dateutil>=2.8.2',
        'six>=1.16.0',
        'eventlet==0.30.2',
        'h2==3.2.0',
        'hpack>=3.0.0,<4.0.0',
        'freezegun>=1.2.1',
        'pycryptodome>=3.15.0',
        'pyroute2==0.6.13',
        'aiohttp>=3.8.1',
        'jsonpointer>=2.3',
        # TODO: (GH #12601) make magma compatible with ovs>=2.17.0
        'ovs==2.16.0',
        'prometheus-client>=0.3.1',
        'aioeventlet @ git+https://github.com/magma/deb-python-aioeventlet@86130360db113430370ed6c64d42aee3b47cd619',
        'sdnotify>=0.3.2',
    ],
    extras_require={
        'dev': [
            # Should be kept in sync with the version in python.mk
            'grpcio-tools>=1.46.3,<1.49.0',
            'iperf3>=0.1.11',
            'parameterized==0.8.1',
            'pytest==7.1.2',
            'pytest-cov==3.0.0',
            'scapy==2.4.5',
        ],
    },
)

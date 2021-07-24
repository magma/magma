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
import sys

from setuptools import setup


def os_release():
    release_info = {}
    with open('/etc/os-release', 'r') as f:
        for line in f:
            try:
                k, v = line.rstrip().split('=')
                release_info[k] = v.strip('"')
            except Exception:
                pass
    return release_info


# We can use an environment variable to pass in the package version during
# build. Since we don't distribute this on its own, we don't really care what
# version this represents. 'None' defaults to 0.0.0.
VERSION = os.environ.get('PKG_VERSION', None)

release_info = os_release()
if release_info.get('VERSION_CODENAME', '') == 'focal':
    setup(
        name='lte',
        version=VERSION,
        packages=[
            'magma.enodebd',
            'magma.enodebd.data_models',
            'magma.enodebd.device_config',
            'magma.enodebd.devices',
            'magma.enodebd.devices.experimental',
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
            'magma.pkt_tester',
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
        ],
        package_data={'magma.redirectd.templates': ['*.html']},
        install_requires=[
            'Cython>=0.29.1',
            'pystemd>=0.5.0',
            'fire>=0.2.0',
            'envoy>=0.0.3',
            'glob2>=0.7',
            # lxml required by spyne.
            'lxml==4.6.3',
            'ryu>=4.30',
            'spyne>=2.13.15',
            'scapy==2.4.4',
            'flask>=1.0.2',
            'sentry_sdk>=1.0.0',
            'aiodns>=1.1.1',
            'pymemoize>=1.0.2',
            'wsgiserver>=1.3',
            # pin recursive dependencies of ryu and others
            'chardet==3.0.4',
            'docker==4.0.2',
            'urllib3>=1.25.3',
            'websocket-client==0.56.0',
            'requests>=2.22.0',
            'certifi>=2019.6.16',
            'idna==2.8',
            'python-dateutil==2.8.1',
            'six>=1.12.0',
            'eventlet==0.30.2',
            'h2>=3.2.0',
            'hpack>=3.0',
            'freezegun>=0.3.15',
            'pycryptodome>=3.9.9',
            'pyroute2==0.5.14',
            'aiohttp==3.6.2',
            'json-pointer>=0.1.2',
            'ovs>=2.13',
            'prometheus-client>=0.3.1',
            'aioeventlet==0.5.1',   # aioeventlet-build.sh
        ],
        extras_require={
            'dev': [
                # Keep grpcio and grpcio-tools on same version for now
                # If you update this version here, you probably also want to
                # update it in lte/gateway/python/Makefile
                'grpcio-tools>=1.16.1',
                'nose==1.3.7',
                'iperf3',
            ],
        },
    )
    sys.exit(0)

# debian stretch packages:-
setup(
    name='lte',
    version=VERSION,
    packages=[
        'magma.enodebd',
        'magma.enodebd.data_models',
        'magma.enodebd.device_config',
        'magma.enodebd.devices',
        'magma.enodebd.devices.experimental',
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
        'magma.pkt_tester',
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
    ],
    package_data={'magma.redirectd.templates': ['*.html']},
    install_requires=[
        'Cython>=0.29.1',
        'pystemd>=0.5.0',
        'fire>=0.2.0',
        'envoy>=0.0.3',
        'glob2>=0.7',
        # lxml required by spyne.
        'lxml==4.6.3',
        'ryu>=4.30',
        'spyne>=2.12.16',
        'scapy==2.4.4',
        'flask>=1.0.2',
        'sentry_sdk>=1.0.0',
        'aioeventlet>=0.4',
        'aiodns>=1.1.1',
        'pymemoize>=1.0.2',
        'wsgiserver>=1.3',
        # pin recursive dependencies of ryu and others
        'chardet==3.0.4',
        'docker==4.0.2',
        'urllib3>=1.25.3',
        'websocket-client==0.56.0',
        'requests>=2.22.0',
        'certifi>=2019.6.16',
        'idna==2.8',
        'python-dateutil==2.8.1',
        'six>=1.12.0',
        'eventlet==0.30.2',
        'h2>=3.2.0',
        'hpack>=3.0',
        'freezegun>=0.3.15',
        'pycryptodome>=3.9.9',
        'pyroute2==0.5.14',
    ],
    extras_require={
        'dev': [
            # Keep grpcio and grpcio-tools on same version for now
            # If you update this version here, you probably also want to
            # update it in lte/gateway/python/Makefile
            'grpcio-tools>=1.16.1',
            'nose==1.3.7',
            'iperf3',
        ],
    },
)

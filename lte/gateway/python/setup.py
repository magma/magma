"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import os

from setuptools import setup

# We can use an environment variable to pass in the package version during
# build. Since we don't distribute this on its own, we don't really care what
# version this represents. 'None' defaults to 0.0.0.
VERSION = os.environ.get('PKG_VERSION', None)

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
        'magma.pipelined',
        'magma.pipelined.app',
        'magma.pipelined.openflow',
        'magma.pipelined.qos',
        'magma.pkt_tester',
        'magma.policydb',
        'magma.policydb.servicers',
        'magma.redirectd',
        'magma.redirectd.templates',
        'magma.subscriberdb',
        'magma.subscriberdb.crypto',
        'magma.subscriberdb.protocols',
        'magma.subscriberdb.protocols.diameter',
        'magma.subscriberdb.protocols.diameter.application',
        'magma.subscriberdb.store',
    ],
    scripts=[
        'scripts/agw_health_cli.py',
        'scripts/create_oai_certs.py',
        'scripts/enodebd_cli.py',
        'scripts/fake_user.py',
        'scripts/feg_hello_cli.py',
        'scripts/generate_oai_config.py',
        'scripts/hello_cli.py',
        'scripts/mobility_cli.py',
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
        'scripts/subscriber_cli.py',
        'scripts/spgw_service_cli.py',
    ],
    package_data={'magma.redirectd.templates': ['*.html']},
    install_requires=[
        'Cython>=0.29.1',
        'pystemd==0.5.0',
        'fire>=0.2.0',
        'envoy>=0.0.3',
        'glob2>=0.7',
        # lxml required by spyne.
        'lxml==4.2.1',
        'ryu>=4.30',
        'spyne==2.12.16',
        # scapy version 2.4.2 has an issue of not having LICENSE file in pypi
        # version resulting in error (this is a temporary fix)
        'scapy==2.4.3rc3',
        'flask>=1.0.2',
        'aiodns>=1.1.1',
        'pymemoize>=1.0.2',
        'wsgiserver>=1.3',
        'pycrypto>=2.6.1',
        # pin recursive dependencies of ryu and others
        'chardet==3.0.4',
        'docker==4.0.2',
        'urllib3==1.25.3',
        'websocket-client==0.56.0',
        'requests==2.22.0',
        'certifi==2019.6.16',
        'idna==2.8',
        'python-dateutil==2.8.1',
        'six>=1.12.0',
        'eventlet>=0.24',
        'h2>=3.2.0',
        'hpack>=3.0'
    ],
    extras_require={
        'dev': [
            # Keep grpcio and grpcio-tools on same version for now
            # If you update this version here, you probably also want to
            # update it in lte/gateway/python/Makefile
            'grpcio-tools==1.16.1',
            'nose==1.3.7',
            'pyroute2',
            'iperf3',
        ]
    },
)

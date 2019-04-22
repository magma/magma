#!/usr/bin/env python3
"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
import datetime
import sys

# NOTE: Uncomment following lines to get verbose logs on GRPC error.
# import os
# os.environ['GRPC_TRACE'] = 'all'
# os.environ['GRPC_VERBOSITY'] = 'DEBUG'
import psutil
import snowflake
from magma.common.cert_utils import load_cert
from magma.common.cert_validity import create_ssl_connection, \
    create_tcp_connection
from magma.common.service_registry import ServiceRegistry
from magma.configuration.service_configs import load_service_config
from orc8r.protos.magmad_pb2 import CheckinRequest, SystemStatus
from orc8r.protos.magmad_pb2_grpc import CheckindStub


def main():
    control_proxy_config = load_service_config('control_proxy')
    host = control_proxy_config['cloud_address']
    port = control_proxy_config['cloud_port']
    certfile = control_proxy_config['gateway_cert']
    keyfile = control_proxy_config['gateway_key']
    loop = asyncio.get_event_loop()

    err_suggestions = {
        'tcp':
            """
            - Verify hostname in /etc/magma/control_proxy.yml &
              /var/opt/magma/configs/control_proxy.yml.
            - Verify correct port.
            - Check DNS (nslookup hostname).
            - Make sure you are disconnected from VPN.
            - Check nghttpx is running on cloud VM.
            """,
        'certs':
            """
            - Delete certs:
                1. Delete gateway.key and gateway.crt in /var/opt/magma/certs.
                2. Restart magmad (sudo service magma@magmad restart).
            - Double check that cloud is up and gateway has been added with
              correct hardware id and key to allow bootstrap.
                1. Run show_gateway_info.py.
                2. Go to cloud swagger (EG. https://127.0.0.1:9443/apidocs).
                3. POST to add a new gateway. Fill JSON with corresponding
                   values from step 1.
            """,
        'ssl':
            """
            - Certificate may be valid but invalid for this host.
            - Delete certs:
                1. Delete gateway.key and gateway.crt in /var/opt/magma/certs.
                2. Restart magmad (sudo service magma@magmad restart).
            """,
        'direct_checkin':
            """
            - Verify checkind service is running on cloud VM.
            - Check logs for more information (sudo tail -f /var/log/syslog).
            """,
        'proxy_checkin':
            """
            - Verify control_proxy service is running.
            - Check logs for more information (sudo tail -f /var/log/syslog).
            """,
    }

    try:
        print('1. -- Testing TCP connection to %s:%d -- ' % (host, port))
        stage = 'tcp'
        loop.run_until_complete(
            create_tcp_connection(host, port, loop)
        )

        print('2. -- Testing Certificate -- ')
        stage = 'certs'
        test_check_cert(certfile)

        print('3. -- Testing SSL -- ')
        stage = 'ssl'
        loop.run_until_complete(
            create_ssl_connection(host, port, certfile, keyfile, loop)
        )

        print('4. -- Creating direct cloud checkin -- ')
        stage = 'direct_checkin'
        test_send_checkin(proxy_cloud_connections=False)

        print('5. -- Creating proxy cloud checkin -- ')
        stage = 'proxy_checkin'
        test_send_checkin(proxy_cloud_connections=True)

        print('Success!')

    except Exception as e:
        print('> Error: %s' % e,)
        print("Suggestions:", err_suggestions[stage])
        sys.exit(1)

    finally:
        loop.close()


def test_check_cert(certfile):
    """Determine whether cert is expired, soon expiring, or not yet valid."""
    cert = load_cert(certfile)

    now = datetime.datetime.utcnow()
    if now > cert.not_valid_after:
        raise Exception("Certificate has expired!")

    elif now + datetime.timedelta(hours=20) > cert.not_valid_after:
        print('> Certificate expiring soon: %s' % cert.not_valid_after)

    elif now < cert.not_valid_before:
        raise Exception('Certificate is not yet valid!')


def create_checkin_request():
    """Create request object to send with Checkin"""
    cpu = psutil.cpu_times()
    mem = psutil.virtual_memory()
    request = CheckinRequest(
        gateway_id=snowflake.snowflake(),
        magma_pkg_version='fake_version',
        system_status=SystemStatus(
            cpu_user=int(cpu.user * 1000),  # convert second to millisecond
            cpu_system=int(cpu.system * 1000),
            cpu_idle=int(cpu.idle * 1000),
            mem_total=mem.total,
            mem_available=mem.available,
            mem_used=mem.used,
            mem_free=mem.free,
        ),
    )
    return request


def test_send_checkin(proxy_cloud_connections=True):
    """Send checkin using either proxy or direct to cloud connection"""
    chan = ServiceRegistry.get_rpc_channel(
            'checkind', ServiceRegistry.CLOUD,
            proxy_cloud_connections=proxy_cloud_connections)
    client = CheckindStub(chan)
    request = create_checkin_request()

    checkin_timeout = 1000
    client.Checkin(request, checkin_timeout)


if __name__ == '__main__':
    main()

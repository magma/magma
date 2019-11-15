#!/usr/bin/env python3
"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.

Pre-run script for services to generate a nghttpx config from a jinja template
and the config/mconfig for the service.
"""

import logging

from generate_service_config import generate_template_config
from magma.configuration import load_service_config
from magma.configuration.mconfig_managers import load_service_mconfig
from orc8r.protos.mconfig.mconfigs_pb2 import FluentBit

CONFIG_OVERRIDE_DIR = '/var/opt/magma/tmp'


def _get_extra_tags():
    """
    Get the extra_tags specified in the FluentBit mconfig.
    """
    return load_service_mconfig('td-agent-bit', FluentBit()).extra_tags


def main():
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s',
    )

    control_proxy_config = load_service_config('control_proxy')
    host = control_proxy_config['fluentd_address']
    port = control_proxy_config['fluentd_port']
    cacert = control_proxy_config['rootca_cert']
    certfile = control_proxy_config['gateway_cert']
    keyfile = control_proxy_config['gateway_key']

    context = {
        'extra_tags': _get_extra_tags().items(),
        'host': host,
        'port': port,
        'cacert': cacert,
        'certfile': certfile,
        'keyfile': keyfile,
    }
    generate_template_config(
        'td-agent-bit', 'td-agent-bit', CONFIG_OVERRIDE_DIR, context.copy()
    )


if __name__ == '__main__':
    main()

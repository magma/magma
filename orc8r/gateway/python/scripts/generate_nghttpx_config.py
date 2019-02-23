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

import os
from magma.common.service_registry import ServiceRegistry
from magma.configuration.environment import is_dev_mode
from magma.configuration.service_configs import get_service_config_value

from generate_service_config import generate_template_config

CONFIG_OVERRIDE_DIR = '/var/opt/magma/tmp'


def get_context():
    """
    Create the context to be used for nghttpx, other than the one provided
    by the configs.
    """
    context = {}
    context['backends'] = []
    for service in ServiceRegistry.list_services():
        (ip_address, port) = ServiceRegistry.get_service_address(service)
        backend = {'service': service, 'ip': ip_address, 'port': port}
        context['backends'].append(backend)

    # We get the gateway cert after bootstrapping, but we do want nghttpx
    # to run before that for communication locally. Update the flag for
    # jinja to act upon.
    gateway_cert = get_service_config_value('control_proxy',
                                            'gateway_cert', None)
    if gateway_cert and os.path.exists(gateway_cert):
        context['use_gateway_cert'] = True
    else:
        context['use_gateway_cert'] = False

    context['dev_mode'] = is_dev_mode()

    context['allow_http_proxy'] = get_service_config_value(
        'control_proxy', 'allow_http_proxy', False)
    context['http_proxy'] = os.getenv('http_proxy', '')
    return context


def main():
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s')
    generate_template_config('control_proxy', 'nghttpx',
                             CONFIG_OVERRIDE_DIR, get_context())


if __name__ == "__main__":
    main()

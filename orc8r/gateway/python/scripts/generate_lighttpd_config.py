#!/usr/bin/env python3

"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging

from magma.common.misc_utils import get_ip_from_if
from magma.configuration.exceptions import LoadConfigError
from magma.configuration.mconfig_managers import load_service_mconfig
from magma.configuration.service_configs import load_service_config

from generate_service_config import generate_template_config

CONFIG_OVERRIDE_DIR = '/var/opt/magma/tmp'


def get_context():
    """
    Provide context to pass to Jinja2 for templating.
    """
    context = {}
    cfg = load_service_config("lighttpd")
    ip = "127.0.0.1"
    enable_caching = False
    try:
        mconfig = load_service_mconfig('lighttpd')
        enable_caching = mconfig.enable_caching
    except LoadConfigError:
        logging.info("Using default values for service 'lighttpd'")

    if enable_caching:
        ip = get_ip_from_if(cfg['interface'])

    context['interface_ip'] = ip
    context['store_root'] = cfg['store_root']

    return context


def main():
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s')

    generate_template_config('lighttpd', 'lighttpd',
                             CONFIG_OVERRIDE_DIR, get_context())


if __name__ == '__main__':
    main()

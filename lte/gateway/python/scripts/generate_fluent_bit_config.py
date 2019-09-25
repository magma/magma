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

CONFIG_OVERRIDE_DIR = "/var/opt/magma/tmp"


def main():
    logging.basicConfig(
        level=logging.INFO, format="[%(asctime)s %(levelname)s %(name)s] %(message)s"
    )
    generate_template_config("td-agent-bit", "td-agent-bit", CONFIG_OVERRIDE_DIR, {})


if __name__ == "__main__":
    main()

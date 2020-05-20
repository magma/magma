#!/usr/bin/env python3

#  Copyright (c) Facebook, Inc. and its affiliates.
#  All rights reserved.
#
#  This source code is licensed under the BSD-style license found in the
#  LICENSE file in the root directory of this source tree.
import os
from typing import Any, Dict

import jinja2
import yaml

CONFIGS_DIR = '/etc/magma/configs'
TEMPLATES_DIR = '/etc/magma/templates'
OUTPUT_DIR = '/etc/nginx'


def _load_services() -> Dict[Any, Any]:
    services = {}  # type: Dict[Any, Any]
    modules = os.listdir(CONFIGS_DIR)
    for module in modules:
        print("Loading registry for module: %s..." % module)
        filename = os.path.join(CONFIGS_DIR, module, "service_registry.yml")
        with open(filename) as file:
            registry = yaml.safe_load(file)
            if registry and "services" in registry:
                services.update(registry["services"])
    return services


def _generate_config(context: Dict[str, Any]) -> str:
    loader = jinja2.FileSystemLoader(TEMPLATES_DIR)
    env = jinja2.Environment(loader=loader)
    template = env.get_template("nginx.conf.j2")
    output = template.render(context)
    outfile = os.path.join(OUTPUT_DIR, "nginx.conf")
    with open(outfile, "w") as file:
        file.write(output)
    return outfile


def main():
    context = {
        'service_registry': _load_services(),
        'controller_hostname': os.environ['CONTROLLER_HOSTNAME'],
        'backend': os.environ['PROXY_BACKENDS'],
    }
    # We need to set the resolver when running in docker-compose
    # In k8s the /etc/resolv.conf has search domains so this isn't needed
    if os.environ['TEST_MODE'] == '1':
        context['resolver'] = '127.0.0.11'
    _generate_config(context)


if __name__ == '__main__':
    main()

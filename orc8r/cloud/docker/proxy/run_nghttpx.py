#!/usr/bin/env python3
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

import argparse
import logging
import os
import subprocess
from typing import Any, Dict

import jinja2
import yaml

CONFIGS_DIR = "/etc/magma/configs"
TEMPLATES_DIR = "/etc/magma/templates"
OUTPUT_DIR = "/etc/nghttpx"
OBSIDIAN_PORT = 9081


def _load_services() -> Dict[Any, Any]:
    """ Return the services from the registry configs of all modules """
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


def _generate_config(proxy_type: str, context: Dict[str, Any]) -> str:
    """ Generate the nghttpx config from the template """
    loader = jinja2.FileSystemLoader(TEMPLATES_DIR)
    env = jinja2.Environment(loader=loader)
    template = env.get_template("nghttpx_%s.conf.j2" % proxy_type)
    output = template.render(context)
    outfile = os.path.join(OUTPUT_DIR, "nghttpx_%s.conf" % proxy_type)
    with open(outfile, "w") as file:
        file.write(output)
    return outfile


def _run_nghttpx(conf: str) -> None:
    """ Runs the nghttpx process given the config file """
    try:
        subprocess.run(
            [
                "/usr/local/bin/nghttpx",
                "--conf=%s" % conf,
                "/var/opt/magma/certs/controller.key",
                "/var/opt/magma/certs/controller.crt",
            ], check=True,
        )
    except subprocess.CalledProcessError as err:
        exit(err.returncode)


def main() -> None:
    parser = argparse.ArgumentParser(description="Nghttpx runner")
    parser.add_argument("proxy_type", choices=["open", "clientcert"])
    args = parser.parse_args()

    # Create the jinja context
    context = {}  # type: Dict[str, Any]
    context["service_registry"] = _load_services()
    context["controller_hostname"] = os.environ["CONTROLLER_HOSTNAME"]
    context["proxy_backends"] = os.environ["PROXY_BACKENDS"]
    context["obsidian_port"] = OBSIDIAN_PORT
    context["env"] = os.environ

    # Generate the nghttpx config
    conf = _generate_config(args.proxy_type, context)

    # Run the nghttpx process
    _run_nghttpx(conf)
    logging.error("nghttpx restarting")


if __name__ == '__main__':
    main()

#!/usr/bin/env python3
"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Pre-run script for services to generate a config from a jinja template
and the config/mconfig for the service.
The templates should be in /etc/magma/templates and the final config will be in
/var/opt/magma/tmp.
"""

import argparse
import logging
import os

from jinja2 import Template
from magma.common.serialization_utils import write_to_file_atomically
from magma.configuration.exceptions import LoadConfigError
from magma.configuration.mconfig_managers import load_service_mconfig_as_json
from magma.configuration.service_configs import load_service_config
from snowflake import make_snowflake


def _get_template_filename(template):
    """
    Returns the location of the template file

    Templates will be looked up in these directories:

        - /etc/magma/templates (default)
        - /var/opt/magma/templates (override location)

    Args:
        template (str): Name of the template. The template file would
                        be /etc/magma/templates/{template}.conf.template
    """
    dirname = '/etc/magma/templates'
    override_dirname = '/var/opt/magma/templates'
    default_filename = os.path.join(
        dirname,
        '%s.conf.template' % template,
    )
    override_filename = os.path.join(
        override_dirname,
        '%s.conf.template' % template,
    )
    if os.path.exists(override_filename) and os.path.isfile(override_filename):
        return override_filename
    return default_filename


def _get_template_out_filename(template, dirname):
    """
    Returns the location of the final config.

    Args:
        template (str): Name of the template used. This will return the
                        output filename as /var/opt/magma/tmp/{template}.conf
        dirname (str): Path of the output file
    """
    return os.path.join(dirname, '%s.conf' % template)


def generate_template_config(service, template, out_dirname, context):
    """
    Generate the config from the jinja template.

    Args:
        service (str): Name of the magma service. Used for looking up the
                        config and mconfig
        template (str): Name of the input template, which is also used for
                        choosing the output filename
        out_dirname (str): Path of the output file
        context (map): Context to use for Jinja (the .yml config and mconfig
                        will be added into this context)
    """
    # Get the template and the output filenames
    template_filename = _get_template_filename(template)
    out_filename = _get_template_out_filename(template, out_dirname)
    logging.info(
        "Generating config file: [%s] using template: [%s]" % (
        out_filename, template_filename,
        ),
    )
    template_context = {}
    # Generate the content to use from the service yml config and mconfig.
    try:
        template_context.update(load_service_config(service))
    except LoadConfigError as err:
        logging.warning(err)

    template_context.update(context)
    try:
        mconfig = load_service_mconfig_as_json(service)
        template_context.update(mconfig)
    except LoadConfigError as err:
        logging.warning(err)

    # Export snowflake to template.
    # TODO: export a hardware-derived ID that can be used by a field tech
    # to easily identify a specific device.
    template_context.setdefault("snowflake", make_snowflake())

    # Create the config file based on the template
    template_str = open(template_filename, 'r').read()
    output = Template(template_str).render(template_context)
    os.makedirs(out_dirname, exist_ok=True)
    write_to_file_atomically(out_filename, output)


def parse_args():
    parser = argparse.ArgumentParser("Templatized service config generator")
    parser.add_argument("--service", required=True)
    parser.add_argument("--template", required=True)
    parser.add_argument("--output-path", default="/var/opt/magma/tmp")

    return parser.parse_args()


def main():
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s',
    )
    opts = parse_args()
    generate_template_config(opts.service, opts.template, opts.output_path, {})


if __name__ == "__main__":
    main()

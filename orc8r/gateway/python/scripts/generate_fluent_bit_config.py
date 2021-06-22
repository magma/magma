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

Pre-run script for services to generate a nghttpx config from a jinja template
and the config/mconfig for the service.
"""

import logging
import os

from generate_service_config import generate_template_config
from magma.configuration import load_service_config
from magma.configuration.mconfig_managers import load_service_mconfig
from orc8r.protos.mconfig.mconfigs_pb2 import FluentBit

CONFIG_OVERRIDE_DIR = '/var/opt/magma/tmp'


def main():
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s',
    )

    mc = load_service_mconfig('td-agent-bit', FluentBit())

    control_proxy_config = load_service_config('control_proxy')
    host = control_proxy_config['fluentd_address']
    port = control_proxy_config['fluentd_port']
    cacert = control_proxy_config['rootca_cert']
    certfile = control_proxy_config['gateway_cert']
    keyfile = control_proxy_config['gateway_key']

    context = {
        'host': host,
        'port': port,
        'cacert': cacert,
        'certfile': certfile,
        'keyfile': keyfile,

        'extra_tags': mc.extra_tags.items(),
        'throttle_rate': mc.throttle_rate or 1000,
        'throttle_window': mc.throttle_window or 5,
        'throttle_interval': mc.throttle_interval or '1m',
        'files': mc.files_by_tag.items(),
    }
    if certfile and os.path.exists(certfile):
        context['is_tls_enabled'] = True
    else:
        context['is_tls_enabled'] = False

    generate_template_config(
        'td-agent-bit', 'td-agent-bit', CONFIG_OVERRIDE_DIR, context.copy(),
    )


if __name__ == '__main__':
    main()

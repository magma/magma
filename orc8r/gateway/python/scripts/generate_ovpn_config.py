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
"""

import logging

from generate_service_config import generate_template_config
from magma.configuration.exceptions import LoadConfigError
from magma.configuration.mconfig_managers import load_service_mconfig

CONFIG_OVERRIDE_DIR = '/etc/openvpn/client'


def get_context():
    """
    Provide context to pass to Jinja2 for templating.
    """
    context = {}
    return context


def main():
    logging.basicConfig(
        level=logging.INFO,
        format='[%(asctime)s %(levelname)s %(name)s] %(message)s')

    generate_template_config('ovpn', 'client', CONFIG_OVERRIDE_DIR, get_context())


if __name__ == '__main__':
    main()

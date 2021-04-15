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

import os
import sys

from setuptools import setup


def os_release():
    release_info = {}
    with open('/etc/os-release', 'r') as f:
        for line in f:
            try:
                k, v = line.rstrip().split('=')
                release_info[k] = v.strip('"')
            except Exception:
                pass
    return release_info

# We can use an environment variable to pass in the package version during
# build. Since we don't distribute this on its own, we don't really care what
# version this represents. 'None' defaults to 0.0.0.
VERSION = os.environ.get('PKG_VERSION', None)

release_info = os_release()
if release_info.get('VERSION_CODENAME', '') == 'focal':
    setup(
        name='orc8r',
        version=VERSION,
        packages=[
            'magma.common',
            'magma.common.health',
            'magma.common.redis',
            'magma.configuration',
            'magma.directoryd',
            'magma.magmad',
            'magma.magmad.generic_command',
            'magma.magmad.check',
            'magma.magmad.check.kernel_check',
            'magma.magmad.check.machine_check',
            'magma.magmad.check.network_check',
            'magma.magmad.upgrade',
            'magma.state',
            'magma.eventd',
            'magma.ctraced',
        ],
        scripts=[
            'scripts/checkin_cli.py',
            'scripts/ctraced_cli.py',
            'scripts/directoryd_cli.py',
            'scripts/generate_lighttpd_config.py',
            'scripts/generate_nghttpx_config.py',
            'scripts/generate_service_config.py',
            'scripts/generate_fluent_bit_config.py',
            'scripts/health_cli.py',
            'scripts/magma_conditional_service.py',
            'scripts/magma_get_config.py',
            'scripts/magmad_cli.py',
            'scripts/service_util.py',
            'scripts/service303_cli.py',
            'scripts/show_gateway_info.py',
            'scripts/traffic_cli.py',
        ],
        install_requires=[
            'setuptools==49.6.0',
            'Cython>=0.29.1',
            'pystemd>=0.5.0',
            'docker>=4.0.2',
            'fire>=0.2.0',
            'glob2>=0.7',
            'aioh2>=0.2.2',
            'redis>=2.10.5',  # redis-py (Python bindings to redis)
            'redis-collections>=0.4.2',
            'python-redis-lock>=3.7.0',
            'aiohttp>=0.17.2',
            'grpcio>=1.16.1',
            'protobuf>=3.14.0',
            'Jinja2>=2.8',
            'netifaces>=0.10.4',
            'pylint>=1.7.1',
            'PyYAML>=3.12',
            'pytz>=2014.4',
            'prometheus_client==0.3.1',
            'sentry_sdk>=1.0.0',
            'snowflake>=0.0.3',
            'psutil==5.6.6',
            'cryptography>=1.9',
            'itsdangerous>=0.24',
            'click>=5.1',
            'pycares>=2.3.0',
            'python-dateutil>=1.4',
            # force same requests version as lte/gateway/python/setup.py
            'requests==2.22.0',
            'jsonpickle',
            'bravado-core==5.16.1',
            'jsonschema==3.1.0',
            "strict-rfc3339>=0.7",
            "rfc3987>=1.3.0",
            "webcolors>=1.11.1",
        ],
        extras_require={
            'dev': [
                "fakeredis[lua]"
            ],
        },
    )
    sys.exit(0)

# debian stretch packages:-
setup(
    name='orc8r',
    version=VERSION,
    packages=[
        'magma.common',
        'magma.common.health',
        'magma.common.redis',
        'magma.configuration',
        'magma.directoryd',
        'magma.magmad',
        'magma.magmad.generic_command',
        'magma.magmad.check',
        'magma.magmad.check.kernel_check',
        'magma.magmad.check.machine_check',
        'magma.magmad.check.network_check',
        'magma.magmad.upgrade',
        'magma.state',
        'magma.eventd',
        'magma.ctraced',
    ],
    scripts=[
        'scripts/checkin_cli.py',
        'scripts/ctraced_cli.py',
        'scripts/directoryd_cli.py',
        'scripts/generate_lighttpd_config.py',
        'scripts/generate_nghttpx_config.py',
        'scripts/generate_service_config.py',
        'scripts/generate_fluent_bit_config.py',
        'scripts/health_cli.py',
        'scripts/magma_conditional_service.py',
        'scripts/magma_get_config.py',
        'scripts/magmad_cli.py',
        'scripts/service_util.py',
        'scripts/service303_cli.py',
        'scripts/show_gateway_info.py',
        'scripts/traffic_cli.py',
    ],
    install_requires=[
        'setuptools==49.6.0',
        'Cython>=0.29.1',
        'pystemd>=0.5.0',
        'docker>=4.0.2',
        'fire>=0.2.0',
        'glob2>=0.7',
        'aioh2>=0.2.2',
        'redis>=2.10.5',  # redis-py (Python bindings to redis)
        'redis-collections>=0.4.2',
        'python-redis-lock>=3.7.0',
        'aiohttp>=0.17.2',
        'grpcio>=1.16.1',
        'protobuf>=3.14.0',
        'Jinja2>=2.8',
        'netifaces>=0.10.4',
        'pylint>=1.7.1',
        'PyYAML>=3.12',
        'pytz>=2014.4',
        'prometheus_client==0.3.1',
        'sentry_sdk>=1.0.0',
        'snowflake>=0.0.3',
        'psutil==5.6.6',
        'cryptography>=1.9',
        'systemd-python>=234',
        'itsdangerous>=0.24',
        'click>=5.1',
        'pycares>=2.3.0',
        'python-dateutil>=1.4',
        # force same requests version as lte/gateway/python/setup.py
        'requests==2.22.0',
        'jsonpickle',
        'bravado-core==5.16.1',
        'jsonschema==3.1.0',
        "strict-rfc3339>=0.7",
        "rfc3987>=1.3.0",
        "jsonpointer>=1.13",
        "webcolors>=1.11.1",
    ],
    extras_require={
        'dev': [
            "fakeredis[lua]"
        ],
    },
)

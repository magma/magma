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

from setuptools import setup

# We can use an environment variable to pass in the package version during
# build. Since we don't distribute this on its own, we don't really care what
# version this represents. 'None' defaults to 0.0.0.
VERSION = os.environ.get('PKG_VERSION', None)

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
        'setuptools>=70.0.0',
        'Cython>=0.29.32',
        'pystemd>=0.10.0',
        'docker>=6.0.0',
        'fire>=0.4.0',
        'glob2>=0.7',
        'aioh2>=0.2.2',
        'redis>=5.0.0',  # redis-py (Python bindings to redis)
        'redis-collections==0.11.0',
        'python-redis-lock>=3.7.0',
        'aiohttp>=3.8.1',
        'grpcio>=1.62.0',
        'protobuf==3.20.3',
        'Jinja2>=3.1.3',
        'markupsafe>=2.1.3',
        'netifaces>=0.11.0',
        'pylint>=1.7.1,<=2.14.0',
        'PyYAML>=6.0.1',
        'pytz>=2022.1',
        'prometheus_client>=0.20.0',
        'sentry_sdk>=1.5.0,<2.0',
        'snowflake>=0.0.3',
        'psutil>=5.9.5',
        'cryptography>=44.0.0',
        'itsdangerous>=2.1.0',
        'click>=8.1.0',
        'pycares>=4.2.1',
        'python-dateutil>=2.8.2',
        # force same requests version as lte/gateway/python/setup.py
        'requests>=2.32.0',
        'jsonpickle>=2.2.0',
        'bravado-core==5.17.0',
        'jsonschema>=4.20.0',
        "strict-rfc3339>=0.7",
        "rfc3987>=1.3.8",
        "webcolors>=1.12",
        'systemd-python>=234',
        "jsonpointer>=2.3",
    ],
    extras_require={
        'dev': [
            "lupa==1.13",
            "fakeredis[lua]==1.8.1",
        ],
    },
)

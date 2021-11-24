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

from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.kernsnoopd.snooper import Snooper

# MIN_COLLECT_INTERVAL overrides counter collection interval received from
# mconfig if it is below this threshold
MIN_COLLECT_INTERVAL = 5


def main():
    """
    Read configuration and run the Snooper job
    """
    # There is no mconfig for kernsnoopd
    service = MagmaService('kernsnoopd', None)

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name, sentry_mconfig=service.shared_mconfig.sentry_config)

    # Initialize and run the Snooper job
    Snooper(
        service.config['ebpf_programs'],
        max(MIN_COLLECT_INTERVAL, service.config['collect_interval']),
        ServiceRegistry,
        service.loop,
    )

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()

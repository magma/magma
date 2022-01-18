"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import importlib
import logging

from magma.protocol_controller.config import Config
from magma.protocol_controller.plugin import ProtocolPlugin

logging.basicConfig(
    level=logging.DEBUG,
    datefmt='%Y-%m-%d %H:%M:%S',
    format='%(asctime)s %(levelname)-8s %(message)s',
)


def get_plugin(config: Config) -> ProtocolPlugin:
    """
    Get Protocol controller plugin class object from configuration

    Parameters:
        config: Protocol controller configuration object

    Returns:
        ProtocolPlugin: protocol controller plugin object
    """
    plugin_module = importlib.import_module(
        '.'.join(config.PROTOCOL_PLUGIN.split('.')[:-1]),
    )
    plugin_class = getattr(
        plugin_module, config.PROTOCOL_PLUGIN.split('.')[-1],
    )
    return plugin_class()


if __name__ == "__main__":
    pc_plugin = get_plugin(Config)
    pc_plugin.initialize()

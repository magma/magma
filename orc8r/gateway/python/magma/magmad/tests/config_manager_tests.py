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
import asyncio
from unittest import TestCase
from unittest.mock import MagicMock, Mock, patch

from google.protobuf.any_pb2 import Any
from google.protobuf.json_format import MessageToJson
from magma.configuration.mconfig_managers import MconfigManagerImpl
from magma.magmad.config_manager import CONFIG_STREAM_NAME, ConfigManager
from orc8r.protos.mconfig.mconfigs_pb2 import MagmaD, MetricsD
from orc8r.protos.mconfig_pb2 import GatewayConfigs
from orc8r.protos.streamer_pb2 import DataUpdate


class ConfigManagerTest(TestCase):
    """
    Tests for the config manager class
    """
    @patch('magma.configuration.service_configs.load_service_config')
    def test_update(self, config_mock):
        """
        Test that mconfig updates are handled correctly
        """
        # Set up fixture data
        # Update will simulate gateway moving from
        # test_mconfig -> updated_mconfig
        test_mconfig = GatewayConfigs()
        updated_mconfig = GatewayConfigs()

        some_any = Any()
        magmad = MagmaD(log_level=1)
        some_any.Pack(magmad)
        test_mconfig.configs_by_key['magmad'].CopyFrom(some_any)
        updated_mconfig.configs_by_key['magmad'].CopyFrom(some_any)

        metricsd = MetricsD(log_level=2)
        some_any.Pack(metricsd)
        test_mconfig.configs_by_key['metricsd'].CopyFrom(some_any)
        metricsd = MetricsD(log_level=3)
        some_any.Pack(metricsd)
        updated_mconfig.configs_by_key['metricsd'].CopyFrom(some_any)

        # Set up mock dependencies
        config_mock.return_value = {
            'magma_services': ['magmad', 'metricsd'],
        }

        @asyncio.coroutine
        def _mock_restart_services(): return "blah"

        service_manager_mock = MagicMock()
        magmad_service_mock = MagicMock()
        mconfig_manager_mock = MconfigManagerImpl()

        load_mock = patch.object(
            mconfig_manager_mock,
            'load_mconfig', MagicMock(return_value=test_mconfig),
        )
        update_mock = patch.object(
            mconfig_manager_mock,
            'update_stored_mconfig', Mock(),
        )
        restart_service_mock = patch.object(
            service_manager_mock,
            'restart_services', MagicMock(wraps=_mock_restart_services),
        )
        processed_updates_mock = patch('magma.magmad.events.processed_updates')

        with load_mock as loader,\
                update_mock as updater, \
                restart_service_mock as restarter,\
                processed_updates_mock as processed_updates:
            loop = asyncio.new_event_loop()
            config_manager = ConfigManager(
                ['magmad', 'metricsd'], service_manager_mock,
                magmad_service_mock, mconfig_manager_mock,
                allow_unknown_fields=False,
                loop=loop,
            )

            # Process an empty set of updates
            updates = []
            config_manager.process_update(CONFIG_STREAM_NAME, updates, False)

            # No services should be updated or restarted due to empty updates
            restarter.assert_not_called()
            updater.assert_not_called()

            # Verify that config update restarts all services
            update_str = MessageToJson(updated_mconfig)
            updates = [
                DataUpdate(value=''.encode('utf-8'), key='some key'),
                DataUpdate(
                    value=update_str.encode('utf-8'),
                    key='last key',
                ),
            ]
            config_manager.process_update(CONFIG_STREAM_NAME, updates, False)

            # Only metricsd config was updated, hence should be restarted
            loader.assert_called_once_with()
            restarter.assert_called_once_with(['metricsd'])
            updater.assert_called_once_with(update_str)

            configs_by_service = {
                'magmad': updated_mconfig.configs_by_key['magmad'],
                'metricsd': updated_mconfig.configs_by_key['metricsd'],
            }
            processed_updates.assert_called_once_with(configs_by_service)

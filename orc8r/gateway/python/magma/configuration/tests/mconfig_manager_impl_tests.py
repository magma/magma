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

import unittest
from unittest import mock

from google.protobuf import any_pb2
from google.protobuf.json_format import MessageToJson
from magma.configuration import mconfig_managers
from magma.configuration.exceptions import LoadConfigError
from orc8r.protos.mconfig import mconfigs_pb2


class MconfigManagerImplTest(unittest.TestCase):
    @mock.patch('magma.configuration.service_configs.load_service_config')
    def test_load_mconfig(self, get_service_config_value_mock):
        # Fixture mconfig has 1 unrecognized service, 1 unregistered type
        magmad_fixture = mconfigs_pb2.MagmaD(
            checkin_interval=10,
            checkin_timeout=5,
            autoupgrade_enabled=True,
            autoupgrade_poll_interval=300,
            package_version='1.0.0-0',
            images=[],
            tier_id='default',
            feature_flags={'flag1': False},
        )
        magmad_fixture_any = any_pb2.Any()
        magmad_fixture_any.Pack(magmad_fixture)
        magmad_fixture_serialized = MessageToJson(magmad_fixture_any)
        fixture = '''
        {
            "configs_by_key": {
                "magmad": %s,
                "foo": {
                    "@type": "type.googleapis.com/magma.mconfig.NotAType",
                    "value": "test1"
                },
                "not_a_service": {
                    "@type": "type.googleapis.com/magma.mconfig.MagmaD",
                    "value": "test2"
                }
            }
        }
        ''' % magmad_fixture_serialized
        get_service_config_value_mock.return_value = {
            'magma_services': ['foo'],
        }

        with mock.patch('builtins.open', mock.mock_open(read_data=fixture)):
            manager = mconfig_managers.MconfigManagerImpl()
            with self.assertRaises(LoadConfigError):
                manager.load_mconfig()

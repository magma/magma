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

from google.protobuf.any_pb2 import Any
from magma.configuration import mconfigs
from orc8r.protos.mconfig import mconfigs_pb2


class MconfigsTest(unittest.TestCase):

    @mock.patch('magma.configuration.service_configs.load_service_config')
    def test_filter_configs_by_key(self, load_service_config_mock):
        # All services present, but 1 type not
        configs_by_key = {
            'magmad': {
                '@type': 'type.googleapis.com/magma.mconfig.MagmaD',
                'value': 'world'.encode(),
            },
            'directoryd': {
                '@type': 'type.googleapis.com/magma.mconfig.DirectoryD',
                'value': 'hello'.encode(),
            },
            'foo': {
                '@type': 'type.googleapis.com/magma.mconfig.Foo',
                'value': 'test'.encode(),
            },
        }

        # Directoryd not present
        load_service_config_mock.return_value = {
            'magma_services': ['mme', 'foo'],
        }
        actual = mconfigs.filter_configs_by_key(configs_by_key)
        expected = {
            'magmad': configs_by_key['magmad'],
            'foo': configs_by_key['foo'],
        }
        self.assertEqual(expected, actual)

        # No services present
        load_service_config_mock.return_value = {
            'magma_services': [],
        }
        actual = mconfigs.filter_configs_by_key(configs_by_key)
        expected = {'magmad': configs_by_key['magmad']}
        self.assertEqual(expected, actual)

        # Directoryd service present as a dynamic service
        load_service_config_mock.return_value = {
            'magma_services': [],
            'registered_dynamic_services': ['directoryd'],
        }
        actual = mconfigs.filter_configs_by_key(configs_by_key)
        expected = {
            'magmad': configs_by_key['magmad'],
            'directoryd': configs_by_key['directoryd'],
        }
        self.assertEqual(expected, actual)

        # Including 'shared_mconfig'
        configs_by_key['shared_mconfig'] = {
            '@type': 'type.googleapis.com/magma.mconfig.SharedMconfig',
            'value': 'shared'.encode(),
        }
        actual = mconfigs.filter_configs_by_key(configs_by_key)
        expected = {
            'magmad': configs_by_key['magmad'],
            'directoryd': configs_by_key['directoryd'],
            'shared_mconfig': configs_by_key['shared_mconfig'],
        }
        self.assertEqual(expected, actual)

    def test_unpack_mconfig_any(self):
        magmad_mconfig = mconfigs_pb2.MagmaD(
            checkin_interval=10,
            checkin_timeout=5,
            autoupgrade_enabled=True,
            autoupgrade_poll_interval=300,
            package_version='1.0.0-0',
            images=[],
            tier_id='default',
            feature_flags={'flag1': False},
        )
        magmad_any = Any(
            type_url='type.googleapis.com/magma.mconfig.MagmaD',
            value=magmad_mconfig.SerializeToString(),
        )
        actual = mconfigs.unpack_mconfig_any(magmad_any, mconfigs_pb2.MagmaD())
        self.assertEqual(magmad_mconfig, actual)

    def test_unpack_mconfig_directoryd(self):
        directoryd_mconfig = mconfigs_pb2.DirectoryD(
            log_level=5,
        )
        magmad_any = Any(
            type_url='type.googleapis.com/magma.mconfig.DirectoryD',
            value=directoryd_mconfig.SerializeToString(),
        )

        actual = mconfigs.unpack_mconfig_any(
            magmad_any, mconfigs_pb2.DirectoryD(),
        )
        self.assertEqual(directoryd_mconfig, actual)

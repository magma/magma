"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
from unittest import mock

from google.protobuf.any_pb2 import Any
from magma.configuration import mconfigs
from orc8r.protos.mconfig import mconfigs_pb2


class MconfigsTest(unittest.TestCase):

    @mock.patch('magma.configuration.service_configs.get_service_config_value')
    def test_filter_configs_by_key(self, get_service_config_value_mock):
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

        get_service_config_value_mock.return_value = ['mme', 'foo']
        actual = mconfigs.filter_configs_by_key(configs_by_key)
        expected = {
            'magmad': configs_by_key['magmad'],
        }
        self.assertEqual(expected, actual)

        # Directoryd service not present
        get_service_config_value_mock.return_value = []
        actual = mconfigs.filter_configs_by_key(configs_by_key)
        expected = {'magmad': configs_by_key['magmad']}
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
        actual = mconfigs.unpack_mconfig_any(magmad_any)
        self.assertEqual(magmad_mconfig, actual)

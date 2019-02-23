"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
import json
import unittest
from unittest import mock

import functools
from google.protobuf.any_pb2 import Any
from google.protobuf.json_format import MessageToJson
from orc8r.protos import common_pb2
from orc8r.protos.mconfig.mconfigs_pb2 import MagmaD, MetricsD
from orc8r.protos.mconfig_pb2 import GatewayConfigs, OffsetGatewayConfigs
from orc8r.protos.streamer_pb2 import DataUpdate

from magma.configuration.mconfig_managers import StreamedMconfigManager
from magma.magmad.streaming_mconfig_callback import MCONFIG_VIEW_STREAM_NAME, \
    StreamingMconfigCallback


class StreamingMconfigCallbackTest(unittest.TestCase):
    @mock.patch('magma.configuration.service_configs.get_service_config_value')
    def test_process_update(self, config_mock):
        """
        Test an update with the following characteristics:
            - There is an unrecognized Any type
            - There is an unrecognized service
            - The magmad mconfig is updated and flips the new feature flag on
            - Another service has its config updated
            - Another service does not have its config updated
        """
        # Set up fixtures
        old_mconfig, new_mconfig = self._get_old_and_new_mconfig_fixtures()
        new_mconfig_ser = MessageToJson(new_mconfig)
        # Add a config with a bad type to the new mconfig (MessageToJson errors
        # out if we include it in the proto message)
        new_mconfig_json_deser = json.loads(new_mconfig_ser)
        new_mconfig_json_deser['configs']['configsByKey']['oops'] = {
            '@type': 'type.googleapis.com/not.a.real.type',
            'value': 'value',
        }
        new_mconfig_ser = json.dumps(new_mconfig_json_deser)

        updates = [
            DataUpdate(
                value=new_mconfig_ser.encode(),
                key='mock',
            ),
        ]

        # Set up mock dependencies
        config_mock.return_value = ['metricsd']

        @asyncio.coroutine
        def _mock_restart_services(): return 'mock'

        service_manager_mock = mock.Mock()
        service_manager_mock.restart_services = mock.MagicMock(
            wraps=_mock_restart_services,
        )

        mconfig_manager_mock = mock.Mock()
        mconfig_manager_mock.load_mconfig.return_value = old_mconfig
        mconfig_manager_mock.update_stored_mconfig = mock.MagicMock()
        # use original impl of deserialize
        mconfig_manager_mock.deserialize_mconfig = mock.MagicMock(
            wraps=functools.partial(
                StreamedMconfigManager.deserialize_mconfig,
                mconfig_manager_mock,
            ),
        )

        magma_service_mock = mock.Mock()
        magma_service_mock.reload_mconfig = mock.MagicMock()

        loop_mock = mock.MagicMock()

        # Run function, assert calls
        callback = StreamingMconfigCallback(
            ['metricsd'], service_manager_mock,
            magma_service_mock, mconfig_manager_mock,
            allow_unknown_fields=False, loop=loop_mock,
        )
        callback.process_update(MCONFIG_VIEW_STREAM_NAME, updates, True)

        # magmad and metricsd configs changed but only metricsd should be
        # restarted (magmad won't restart for config change)
        service_manager_mock.restart_services.assert_called_once_with(
            ['metricsd'],
        )
        # We should have written the whole serialized mconfig, bad keys and all
        mconfig_manager_mock.update_stored_mconfig.assert_called_once_with(
            new_mconfig_ser,
        )
        # Should have reloaded magmad mconfigs to pick up change
        magma_service_mock.reload_mconfig.assert_called_once_with()

    @staticmethod
    def _get_old_and_new_mconfig_fixtures():
        """
        Returns (old_mconfig, new_mconfig) where new_mconfig updates metricsd
        config and magmad config. new_mconfig also adds a config for an
        unrecognized service keyed by 'mock'
        """
        old_magmad_any = Any()
        old_magmad_config = MagmaD(checkin_interval=123)
        old_magmad_any.Pack(old_magmad_config)

        new_magmad_any = Any()
        new_magmad_config = MagmaD(
            checkin_interval=42,
            feature_flags={'kafka_config_streamer': True},
        )
        new_magmad_any.Pack(new_magmad_config)

        old_metricsd_any = Any()
        old_metricsd_config = MetricsD(log_level=common_pb2.ERROR)
        old_metricsd_any.Pack(old_metricsd_config)

        new_metricsd_any = Any()
        new_metricsd_config = MetricsD(log_level=common_pb2.INFO)
        new_metricsd_any.Pack(new_metricsd_config)

        # Unrecognized service, good type
        mock_service_any = Any()
        mock_service_config = MetricsD(log_level=common_pb2.FATAL)
        mock_service_any.Pack(mock_service_config)

        old_mconfig = OffsetGatewayConfigs(
            offset=42,
            configs=GatewayConfigs(
                configs_by_key={
                    'magmad': old_magmad_any,
                    'metricsd': old_metricsd_any,
                },
            ),
        )
        new_mconfig = OffsetGatewayConfigs(
            offset=43,
            configs=GatewayConfigs(
                configs_by_key={
                    'magmad': new_magmad_any,
                    'metricsd': new_metricsd_any,
                    'mock': mock_service_any,
                },
            ),
        )

        return old_mconfig, new_mconfig

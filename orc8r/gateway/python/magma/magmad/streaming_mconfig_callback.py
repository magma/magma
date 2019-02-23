"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
import logging
from typing import Any, List

from orc8r.protos.mconfig_pb2 import MconfigStreamRequest
from orc8r.protos.streamer_pb2 import DataUpdate

from magma.common.service import MagmaService
from magma.common.streamer import StreamerClient
from magma.configuration.mconfig_managers import MconfigManager
from magma.magmad.service_manager import ServiceManager

MCONFIG_VIEW_STREAM_NAME = 'mconfig_views'


class StreamingMconfigCallback(StreamerClient.Callback):
    """
    Stream callback for computed mconfig views (OffsetGatewayConfigs).
    """
    def __init__(self, services: List[str], service_manager: ServiceManager,
                 magmad_service: MagmaService, mconfig_manager: MconfigManager,
                 allow_unknown_fields: bool = True, loop=None) -> None:
        self._services = services
        self._service_manager = service_manager
        self._magmad_service = magmad_service
        self._mconfig_manager = mconfig_manager
        self._allow_unknown_fields = allow_unknown_fields
        self._loop = loop or asyncio.get_event_loop()

        # Load initial mconfig
        self._mconfig = self._mconfig_manager.load_mconfig()

    def get_request_args(self, stream_name: str) -> Any:
        return MconfigStreamRequest(offset=self._mconfig.offset)

    def process_update(self,
                       stream_name: str, updates: List[DataUpdate],
                       resync: bool) -> None:
        if len(updates) == 0:
            logging.info('No streamed mconfig to process')
            return

        for update in updates[:-1]:
            logging.info('Ignoring streamed config update %s', update.key)

        logging.info('Processing streamed config update %s', updates[-1].key)
        mconfig_str = updates[-1].value.decode()
        mconfig = self._mconfig_manager.deserialize_mconfig(
            mconfig_str,
            allow_unknown_fields=self._allow_unknown_fields,
        )

        if 'magmad' not in mconfig.configs.configs_by_key:
            logging.error('Invalid config! Magmad service config missing')
            return

        self._mconfig_manager.update_stored_mconfig(mconfig_str)

        def did_mconfig_change(serv_name):
            return mconfig.configs.configs_by_key.get(serv_name) != \
                   self._mconfig.configs.configs_by_key.get(serv_name)

        # Reload magmad configs locally
        if did_mconfig_change('magmad'):
            self._magmad_service.reload_mconfig()
            magmad_conf = self._mconfig_manager.load_service_mconfig('magmad')
            self._loop.create_task(
                self._service_manager.update_dynamic_services(
                    magmad_conf.dynamic_services,
                )
            )

        services_to_restart = [
            srv for srv in self._services if did_mconfig_change(srv)
        ]
        if services_to_restart:
            self._loop.create_task(
                self._service_manager.restart_services(services_to_restart),
            )

        self._mconfig = mconfig

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
import logging
import re
from typing import Any, List, NamedTuple, Optional

import magma.magmad.events as magmad_events
from magma.common.service import MagmaService
from magma.common.streamer import StreamerClient
from magma.configuration.mconfig_managers import (
    MconfigManager,
    load_service_mconfig,
)
from magma.magmad.service_manager import ServiceManager
from orc8r.protos.mconfig import mconfigs_pb2
from orc8r.protos.mconfig_pb2 import GatewayConfigsDigest

CONFIG_STREAM_NAME = 'configs'
SHARED_MCONFIG = 'shared_mconfig'
MAGMAD = 'magmad'


class VersionInfo(NamedTuple):
    agw_version: Optional[re.Match]
    orc8r_version: Optional[re.Match]


class ConfigManager(StreamerClient.Callback):
    """
    Manager for access gateway config. Updates are received as a stream and
    are guaranteed to be lossless and in-order. Config is written to file in
    JSON format.
    """

    def __init__(
        self, services: List[str], service_manager: ServiceManager,
        magmad_service: MagmaService, mconfig_manager: MconfigManager,
        allow_unknown_fields: bool = True, loop=None,
    ) -> None:
        """
        Args:
            services: List of services to manage
            service_manager: ServiceManager instance
            magmad_service: magmad service instance
            mconfig_manager: manager class for the mconfig
            allow_unknown_fields: set to True to suppress unknown field errors
            loop: asyncio event loop to run in
        """
        self._services = services
        self._service_manager = service_manager
        self._magmad_service = magmad_service
        self._mconfig_manager = mconfig_manager
        self._allow_unknown_fields = allow_unknown_fields
        self._loop = loop or asyncio.get_event_loop()

        # Load managed config
        self._mconfig = self._mconfig_manager.load_mconfig()

    def get_request_args(self, stream_name: str) -> Any:
        # Include an mconfig digest argument to allow cloud optimization of
        # not returning a non-updated mconfig.
        digest = getattr(self._mconfig.metadata, 'digest', None)
        if digest is None:
            return None
        mconfig_digest_proto = GatewayConfigsDigest(
            md5_hex_digest=digest.md5_hex_digest,
        )
        return mconfig_digest_proto

    def _parse_versions_and_log_warning(self, agw_version, unpacked_mconfig):
        # version should be in X.X.X format with only non-negative numbers allowed for X
        version_regex = re.compile(r"[0-9]+\.(?P<minor_version>[0-9]+)\.[0-9]+")

        # unpack the magmad structure to get orce_version field
        if unpacked_mconfig.Is(mconfigs_pb2.MagmaD.DESCRIPTOR):
            magmad_parsed = mconfigs_pb2.MagmaD()
            unpacked_mconfig.Unpack(magmad_parsed)
            orc8r_version = magmad_parsed.orc8r_version

            agw_version_parsed = version_regex.match(agw_version)
            orc8r_version_parsed = version_regex.match(orc8r_version)

            # agw_version is not in expected format
            if not agw_version_parsed:
                logging.warning("Gateway version: %s not valid", agw_version)

            # orc8r_version is not in expected format
            if not orc8r_version_parsed:
                logging.warning("Orchestrator version: %s not valid", orc8r_version)

            return VersionInfo(agw_version_parsed, orc8r_version_parsed)
        logging.error(
            "Expecting MagmaD mconfig structure, but received a different structure: %s.",
            unpacked_mconfig.type_url,
        )
        return VersionInfo(None, None)

    def process_update(self, stream_name, updates, resync):
        """
        Handle config updates. Resync is ignored since the entire config
        structure is passed in every update.
        Inputs:
         - updates - list of GatewayConfigs protobuf structures
         - resync - boolean indicating whether all database information will be
                    resent (hence cached data can be discarded). This is ignored
                    since config is contained in one DB element, hence all
                    data is sent in every update.
        """
        if len(updates) == 0:
            logging.info('No config update to process')
            return

        # We will only take the last update
        for update in updates[:-1]:
            logging.info('Ignoring config update %s', update.key)

        # Deserialize and store the last config update
        logging.info('Processing config update %s', updates[-1].key)
        mconfig_str = updates[-1].value.decode()
        mconfig = self._mconfig_manager.deserialize_mconfig(
            mconfig_str,
            self._allow_unknown_fields,
        )

        if MAGMAD not in mconfig.configs_by_key:
            logging.error('Invalid config! Magmad service config missing')
            return

        self._mconfig_manager.update_stored_mconfig(mconfig_str)
        self._magmad_service.reload_mconfig()

        def did_mconfig_change(serv_name):
            return mconfig.configs_by_key.get(serv_name) != \
                self._mconfig.configs_by_key.get(serv_name)

        # Reload magmad configs locally
        if did_mconfig_change(MAGMAD) or (
            SHARED_MCONFIG in mconfig.configs_by_key
            and did_mconfig_change(SHARED_MCONFIG)
        ):
            logging.info("Restarting dynamic services due to config change")
            self._loop.create_task(
                self._service_manager.update_dynamic_services(
                    load_service_mconfig(MAGMAD, mconfigs_pb2.MagmaD())
                    .dynamic_services,
                ),
            )

        if SHARED_MCONFIG in mconfig.configs_by_key and did_mconfig_change(SHARED_MCONFIG):
            logging.info("Shared config changed. Restarting all services.")
            services_to_restart = self._services
        else:
            services_to_restart = [
                srv for srv in self._services if did_mconfig_change(srv)
            ]

        if services_to_restart:
            self._loop.create_task(
                self._service_manager.restart_services(services_to_restart),
            )

        self._mconfig = mconfig

        configs_by_key = {}
        for srv in self._services:
            if srv in mconfig.configs_by_key:
                configs_by_key[srv] = mconfig.configs_by_key.get(srv)

        agw_version = self._magmad_service.version
        unpacked_mconfig = mconfig.configs_by_key.get(MAGMAD)
        version_info = self._parse_versions_and_log_warning(agw_version, unpacked_mconfig)
        agw_version_parsed = version_info.agw_version
        orc8r_version_parsed = version_info.orc8r_version

        if agw_version_parsed and orc8r_version_parsed:
            agw_minor = int(agw_version_parsed.group('minor_version'))
            orc8r_minor = int(orc8r_version_parsed.group('minor_version'))
            if agw_minor - orc8r_minor <= -1:
                logging.warning("Gateway is more than one minor version behind orc8r. Please consider updating it.")

        magmad_events.processed_updates(configs_by_key)

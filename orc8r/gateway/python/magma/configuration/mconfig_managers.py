"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import contextlib
import json
import logging
from typing import Any, Generic, TypeVar

import abc
import os
from google.protobuf import json_format
from magma.common import serialization_utils
from magma.configuration.exceptions import LoadConfigError
from magma.configuration.mconfigs import filter_configs_by_key, \
    unpack_mconfig_any
from orc8r.protos.mconfig.mconfigs_pb2 import MagmaD
from orc8r.protos.mconfig_pb2 import GatewayConfigs, OffsetGatewayConfigs

T = TypeVar('T')

MCONFIG_DIR = '/etc/magma'
MCONFIG_OVERRIDE_DIR = '/var/opt/magma/configs'
DEFAULT_MCONFIG_DIR = os.environ.get('MAGMA_CONFIG_LOCATION', MCONFIG_DIR)


def get_mconfig_manager():
    """
    Get the mconfig manager implementation that the system is configured to
    use.

    Returns: MconfigManager implementation
    """
    magmad_mconfig = _load_magmad_mconfig_value()
    if magmad_mconfig is None or magmad_mconfig.feature_flags is None:
        return MconfigManagerImpl()

    if magmad_mconfig.feature_flags.get('kafka_config_streamer', False):
        return StreamedMconfigManager()
    else:
        return MconfigManagerImpl()


def _load_magmad_mconfig_value() -> MagmaD:
    # Try to load magmad mconfig value from new mconfigs, then fall back on
    # old mconfigs
    try:
        new_mconfig_manager = StreamedMconfigManager()
        return new_mconfig_manager.load_service_mconfig('magmad')
    except LoadConfigError:
        logging.debug('Could not load offset mconfig, falling back to '
                      'default mconfig manager')

    try:
        old_mconfig_manager = MconfigManagerImpl()
        return old_mconfig_manager.load_service_mconfig('magmad')
    except LoadConfigError as e:
        logging.error(
            'Could not load magmad mconfig value '
            'to check feature flags: %s', e
        )
        return MagmaD()


def load_service_mconfig(service: str) -> Any:
    """
    Utility function to load the mconfig for a specific service using the
    configured mconfig manager.
    """
    return get_mconfig_manager().load_service_mconfig(service)


class MconfigManager(Generic[T]):
    """
    Interface for a class which handles loading and updating some cloud-
    managed configuration (mconfig).
    """

    @abc.abstractmethod
    def load_mconfig(self) -> T:
        """
        Load the managed configuration from its stored location.

        Returns: Loaded mconfig
        """
        pass

    @abc.abstractmethod
    def load_service_mconfig(self, service_name: str) -> Any:
        """
        Load a specific service's managed configuration.

        Args:
            service_name (str): name of the service to load a config for

        Returns: Loaded config value for the service
        """
        pass

    @abc.abstractmethod
    def update_stored_mconfig(self, updated_value: str):
        """
        Update the stored mconfig to the provided serialized value

        Args:
            updated_value: Serialized value of new mconfig value to store
        """
        pass

    @abc.abstractmethod
    def deserialize_mconfig(self, serialized_value: str,
                            allow_unknown_fields: bool = True) -> T:
        """
        Deserialize the given string to the managed mconfig.

        Args:
            serialized_value:
                Serialized value of a managed mconfig
            allow_unknown_fields:
                Set to true to suppress errors from parsing unknown fields

        Returns: deserialized mconfig value
        """
        pass

    @abc.abstractmethod
    def delete_stored_mconfig(self):
        """
        Delete the stored mconfig file.
        """
        pass


class MconfigManagerImpl(MconfigManager[GatewayConfigs]):
    """
    Legacy mconfig manager for non-offset mconfigs
    """

    MCONFIG_FILE_NAME = 'gateway.mconfig'
    MCONFIG_PATH = os.path.join(MCONFIG_OVERRIDE_DIR, MCONFIG_FILE_NAME)

    def load_mconfig(self) -> GatewayConfigs:
        cfg_file_name = self._get_mconfig_file_path()
        try:
            with open(cfg_file_name, 'r') as cfg_file:
                mconfig_str = cfg_file.read()
            return self.deserialize_mconfig(mconfig_str)
        except (OSError, json.JSONDecodeError, json_format.ParseError) as e:
            raise LoadConfigError('Error loading mconfig') from e

    def load_service_mconfig(self, service_name: str) -> Any:
        mconfig = self.load_mconfig()
        if service_name not in mconfig.configs_by_key:
            raise LoadConfigError(
                "Service ({}) missing in mconfig".format(service_name),
            )

        service_mconfig = mconfig.configs_by_key[service_name]
        return unpack_mconfig_any(service_mconfig)

    def deserialize_mconfig(self, serialized_value: str,
                            allow_unknown_fields: bool = True,
                            ) -> GatewayConfigs:
        # First parse as JSON in case there are types unrecognized by
        # protobuf symbol database
        json_mconfig = json.loads(serialized_value)
        cfgs_by_key_json = json_mconfig.get('configs_by_key', {})
        cfgs_by_key_json.update(json_mconfig.get('configsByKey', {}))
        filtered_cfgs_by_key = filter_configs_by_key(cfgs_by_key_json)

        # Set configs to filtered map, re-dump and parse
        if 'configs_by_key' in json_mconfig:
            json_mconfig.pop('configs_by_key')
        json_mconfig['configsByKey'] = filtered_cfgs_by_key
        json_mconfig_dumped = json.dumps(json_mconfig)

        # Workaround for outdated protobuf library on sandcastle
        if allow_unknown_fields:
            return json_format.Parse(
                json_mconfig_dumped,
                GatewayConfigs(),
                ignore_unknown_fields=True,
            )
        else:
            return json_format.Parse(json_mconfig_dumped, GatewayConfigs())

    def delete_stored_mconfig(self):
        with contextlib.suppress(FileNotFoundError):
            os.remove(self.MCONFIG_PATH)

    def update_stored_mconfig(self, updated_value: str) -> GatewayConfigs:
        serialization_utils.write_to_file_atomically(
            self.MCONFIG_PATH, updated_value,
        )

    def _get_mconfig_file_path(self):
        if os.path.isfile(self.MCONFIG_PATH):
            return self.MCONFIG_PATH
        else:
            return os.path.join(DEFAULT_MCONFIG_DIR, self.MCONFIG_FILE_NAME)


class StreamedMconfigManager(MconfigManager[OffsetGatewayConfigs]):
    """
    Manager for OffsetGatewayConfigs mconfigs streamed from the new computed
    mconfig views
    """
    MCONFIG_FILENAME = 'gateway.streamed.mconfig'
    MCONFIG_PATH = os.path.join(MCONFIG_OVERRIDE_DIR, MCONFIG_FILENAME)

    def load_mconfig(self) -> OffsetGatewayConfigs:
        filename = self._get_mconfig_file_path()

        try:
            # load as json object
            # pass in configs_by_key to filter
            with open(filename, 'r') as cfg_file:
                mconfig_str = cfg_file.read()
            return self.deserialize_mconfig(mconfig_str)
        except (OSError, json.JSONDecodeError, json_format.ParseError) as e:
            raise LoadConfigError from e

    def load_service_mconfig(self, service_name: str) -> Any:
        offset_mconfig = self.load_mconfig()
        if service_name not in offset_mconfig.configs.configs_by_key:
            raise LoadConfigError(
                'Service {} missing in mconfig'.format(service_name),
            )

        service_mconfig = offset_mconfig.configs.configs_by_key[service_name]
        return unpack_mconfig_any(service_mconfig)

    def deserialize_mconfig(self, serialized_value: str,
                            allow_unknown_fields: bool = True,
                            ) -> OffsetGatewayConfigs:
        # Parse as json first in case there are types unrecognized by
        # protobuf symbol database
        json_offset_config = json.loads(serialized_value)
        json_mconfig = json_offset_config.get('configs', {})
        cfgs_by_key_json = json_mconfig.get('configs_by_key', {})
        cfgs_by_key_json.update(json_mconfig.get('configsByKey', {}))
        filtered_cfgs_by_key_json = filter_configs_by_key(cfgs_by_key_json)

        # Set the configs by key to the filtered one
        json_offset_config['configs'] = {}
        json_offset_config['configs']['configsByKey'] = \
            filtered_cfgs_by_key_json
        json_offset_config_dumped = json.dumps(json_offset_config)

        # Now re-parse the filtered config
        # Switch on allow_unknown_fields because sandcastle protobuf library
        # is outdated and I have no clue how to update that thing.
        if allow_unknown_fields:
            return json_format.Parse(
                json_offset_config_dumped,
                OffsetGatewayConfigs(),
                ignore_unknown_fields=True,
            )
        else:
            return json_format.Parse(
                json_offset_config_dumped,
                OffsetGatewayConfigs(),
            )

    def update_stored_mconfig(self, updated_value: str):
        serialization_utils.write_to_file_atomically(
            self.MCONFIG_PATH, updated_value,
        )

    def delete_stored_mconfig(self):
        with contextlib.suppress(FileNotFoundError):
            os.remove(self.MCONFIG_PATH)

    def _get_mconfig_file_path(self):
        if os.path.isfile(self.MCONFIG_PATH):
            return self.MCONFIG_PATH
        else:
            return os.path.join(DEFAULT_MCONFIG_DIR, self.MCONFIG_FILENAME)

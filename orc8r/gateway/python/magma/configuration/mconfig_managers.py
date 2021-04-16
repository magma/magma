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
import abc
import contextlib
import json
import os
from typing import Any, Generic, TypeVar

import magma.configuration.events as magma_configuration_events
from google.protobuf import json_format
from magma.common import serialization_utils
from magma.configuration.exceptions import LoadConfigError
from magma.configuration.mconfigs import (
    filter_configs_by_key,
    unpack_mconfig_any,
)
from orc8r.protos.mconfig_pb2 import GatewayConfigs, GatewayConfigsMetadata

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
    # This is stubbed out after deleting the streamed mconfig manager
    return MconfigManagerImpl()


def load_service_mconfig(service: str, mconfig_struct: Any) -> Any:
    """
    Utility function to load the mconfig for a specific service using the
    configured mconfig manager.
    """
    return get_mconfig_manager().load_service_mconfig(service, mconfig_struct)


def load_service_mconfig_as_json(service_name: str) -> Any:
    """
    Loads the managed configuration from its json file stored on disk.

    Args:
        service_name (str): name of the service to load the config for

    Returns: Loaded config value for the service as parsed json struct, not
    protobuf message struct
    """
    return get_mconfig_manager().load_service_mconfig_as_json(service_name)


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
    def load_service_mconfig(
        self, service_name: str,
        mconfig_struct: Any,
    ) -> Any:
        """
        Load a specific service's managed configuration.

        Args:
            service_name (str): name of the service to load a config for
            mconfig_struct (Any): protobuf message struct of the managed config
            for the service

        Returns: Loaded config value for the service
        """
        pass

    @abc.abstractmethod
    def load_mconfig_metadata(self) -> GatewayConfigsMetadata:
        """
        Load the metadata of the managed configuration.

        Returns: Loaded mconfig metadata
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
    def deserialize_mconfig(
        self, serialized_value: str,
        allow_unknown_fields: bool = True,
    ) -> T:
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

    def load_service_mconfig(
        self, service_name: str,
        mconfig_struct: Any,
    ) -> Any:
        mconfig = self.load_mconfig()
        if service_name not in mconfig.configs_by_key:
            raise LoadConfigError(
                "Service ({}) missing in mconfig".format(service_name),
            )

        service_mconfig = mconfig.configs_by_key[service_name]
        return unpack_mconfig_any(service_mconfig, mconfig_struct)

    def load_service_mconfig_as_json(self, service_name) -> Any:
        cfg_file_name = self._get_mconfig_file_path()
        with open(cfg_file_name, 'r') as f:
            json_mconfig = json.load(f)
            service_configs = json_mconfig.get('configsByKey', {})
            service_configs.update(json_mconfig.get('configs_by_key', {}))
        if service_name not in service_configs:
            raise LoadConfigError(
                "Service ({}) missing in mconfig".format(service_name),
            )

        return service_configs[service_name]

    def load_mconfig_metadata(self) -> GatewayConfigsMetadata:
        mconfig = self.load_mconfig()
        return mconfig.metadata

    def deserialize_mconfig(
        self, serialized_value: str,
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
        magma_configuration_events.deleted_stored_mconfig()

    def update_stored_mconfig(self, updated_value: str) -> GatewayConfigs:
        serialization_utils.write_to_file_atomically(
            self.MCONFIG_PATH, updated_value,
        )
        magma_configuration_events.updated_stored_mconfig()

    def _get_mconfig_file_path(self):
        if os.path.isfile(self.MCONFIG_PATH):
            return self.MCONFIG_PATH
        else:
            return os.path.join(DEFAULT_MCONFIG_DIR, self.MCONFIG_FILE_NAME)

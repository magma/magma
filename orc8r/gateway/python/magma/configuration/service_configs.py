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

import logging
import os
from typing import Optional  # noqa: lint doesn't handle inline typehints
from typing import Any, Dict

import yaml
from magma.configuration.exceptions import LoadConfigError

# Location of configs (both service config and mconfig)
CONFIG_DIR = '/etc/magma'
CONFIG_OVERRIDE_DIR = '/var/opt/magma/configs'


def load_override_config(service_name: str) -> Optional[Any]:
    """
    Load override service configuration from the file in the override
    directory.

    Args:
        service_name: service to pull configs for; name of config file

    Returns: json-decoded value of the service config, None if it's not found

    Raises:
        LoadConfigError:
            Unable to load config due to missing file or missing key
    """
    override_file_name = _override_file_name(service_name)
    if os.path.isfile(override_file_name):
        return _load_yaml_file(override_file_name)
    return None


def save_override_config(service_name: str, cfg: Any):
    """
    Write the configuration object to its corresponding file in the override
    directory.

    Args:
        service_name: service to write config object to; name of config file
        cfg: json-decoded value of the service config
    """
    override_file_name = _override_file_name(service_name)
    os.makedirs(CONFIG_OVERRIDE_DIR, exist_ok=True)
    with open(override_file_name, 'w') as override_file:
        yaml.dump(cfg, override_file, default_flow_style=False)


def load_service_config(service_name: str) -> Any:
    """
    Load service configuration from file. Also check override directory,
    and, if service file present there, override the values.

    Args:
        service_name: service to pull configs for; name of config file

    Returns: json-decoded value of the service config

    Raises:
        LoadConfigError:
            Unable to load config due to missing file or missing key
    """
    cfg_file_name = os.path.join(CONFIG_DIR, '%s.yml' % service_name)
    cfg = _load_yaml_file(cfg_file_name)

    overrides = load_override_config(service_name)
    if overrides is not None:
        # Update the keys in the config if they are present in the override
        cfg.update(overrides)
    return cfg


cached_service_configs = {}     # type: Dict[str, Any]


def get_service_config_value(service: str, param: str, default: Any) -> Any:
    """
    Get a config value for :service:, falling back to a :default: value.

    Log error if the default config is returned.

    Args:
        service: name of service to get config for
        param: config key to fetch the value for
        default: default value to return on failure

    Returns:
        value of :param: in the config files for :service:
    """
    service_configs = cached_service_configs.get(service)
    try:
        service_configs = service_configs or load_service_config(service)
    except LoadConfigError as e:
        logging.error('Error retrieving config: %s', e)
        return default

    # Handle empty file
    if not service_configs:
        logging.error('Error retrieving config, file empty for: %s', service)
        return default

    cached_service_configs[service] = service_configs

    config_value = service_configs.get(param)
    if config_value is not None:
        return config_value
    else:
        logging.error(
            'Error retrieving config for %s, key not found: %s',
            service, param,
        )
        return default


def _override_file_name(service_name: str) -> str:
    return os.path.join(CONFIG_OVERRIDE_DIR, '%s.yml' % service_name)


def _load_yaml_file(file_name: str) -> Any:
    """
    Load the yaml file and returns the python object.

    Args:
        file_name: name of the .yml file

    Returns:
        Contents of the yml file deserialized into a Python object

    Raises:
        LoadConfigError: on error
    """

    try:
        with open(file_name, 'r') as stream:
            data = yaml.safe_load(stream)
            return data
    except (OSError, yaml.YAMLError) as e:
        raise LoadConfigError('Error loading yml config') from e

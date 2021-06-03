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

import json
from typing import Any, List

from magma.enodebd.data_models.data_model import DataModel
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.exceptions import ConfigurationError
from magma.enodebd.logger import EnodebdLogger as logger


class EnodebConfiguration():
    """
    This represents the data model configuration for a single
    eNodeB device. This can correspond to either the current configuration
    of the device, or what configuration we desire to have for the device.
    """

    def __init__(self, data_model: DataModel) -> None:
        """
        The fields initialized in the constructor here should be enough to
        track state across any data model configuration.

        Most objects for eNodeB data models cannot be added or deleted.
        For those objects, we just track state with a simple mapping from
        parameter name to value.

        For objects which can be added/deleted, we track them separately.
        """

        # DataModel
        self._data_model = data_model

        # Dict[ParameterName, Any]
        self._param_to_value = {}

        # Dict[ParameterName, Dict[ParameterName, Any]]
        self._numbered_objects = {}
        # If adding a PLMN object, then you would set something like
        # self._numbered_objects['PLMN_1'] = {'PLMN_1_ENABLED': True}

    @property
    def data_model(self) -> DataModel:
        """
        The data model configuration is tied to a single data model
        """
        return self._data_model

    def get_parameter_names(self) -> List[ParameterName]:
        """
        Returns: list of ParameterName
        """
        return list(self._param_to_value.keys())

    def has_parameter(self, param_name: ParameterName) -> bool:
        return param_name in self._param_to_value

    def get_parameter(self, param_name: ParameterName) -> Any:
        """
        Args:
            param_name: ParameterName
        Returns:
            Any, value of the parameter, formatted to be understood by enodebd
        """
        self._assert_param_in_model(param_name)
        return self._param_to_value[param_name]

    def set_parameter(
        self,
        param_name: ParameterName,
        value: Any,
    ) -> None:
        """
        Args:
            param_name: the parameter name to configure
            value: the value to set, formatted to be understood by enodebd
        """
        self._assert_param_in_model(param_name)
        self._param_to_value[param_name] = value

    def delete_parameter(self, param_name: ParameterName) -> None:
        del self._param_to_value[param_name]

    def get_object_names(self) -> List[ParameterName]:
        return list(self._numbered_objects.keys())

    def has_object(self, param_name: ParameterName) -> bool:
        """
        Args:
            param_name: The ParameterName of the object
        Returns: True if set in configuration
        """
        self._assert_param_in_model(param_name)
        return param_name in self._numbered_objects

    def add_object(self, param_name: ParameterName) -> None:
        if param_name in self._numbered_objects:
            raise ConfigurationError("Configuration already has object")
        self._numbered_objects[param_name] = {}

    def delete_object(self, param_name: ParameterName) -> None:
        if param_name not in self._numbered_objects:
            raise ConfigurationError("Configuration does not have object")
        del self._numbered_objects[param_name]

    def get_parameter_for_object(
        self,
        param_name: ParameterName,
        object_name: ParameterName,
    ) -> Any:
        return self._numbered_objects[object_name].get(param_name)

    def set_parameter_for_object(
        self,
        param_name: ParameterName,
        value: Any,
        object_name: ParameterName,
    ) -> None:
        """
        Args:
            param_name: the parameter name to configure
            value: the value to set, formatted to be understood by enodebd
            object_name: ParameterName of object
        """
        self._assert_param_in_model(object_name)
        self._assert_param_in_model(param_name)
        self._numbered_objects[object_name][param_name] = value

    def get_parameter_names_for_object(
        self,
        object_name: ParameterName,
    ) -> List[ParameterName]:
        return list(self._numbered_objects[object_name].keys())

    def get_debug_info(self) -> str:
        debug_info = 'Param values: {}, \n Object values: {}'
        return debug_info.format(
            json.dumps(self._param_to_value, indent=2),
            json.dumps(
                self._numbered_objects,
                indent=2,
            ),
        )

    def _assert_param_in_model(self, param_name: ParameterName) -> None:
        trparam_model = self.data_model
        tr_param = trparam_model.get_parameter(param_name)
        if tr_param is None:
            logger.warning('Parameter <%s> not defined in model', param_name)
            raise ConfigurationError("Parameter not defined in model.")
